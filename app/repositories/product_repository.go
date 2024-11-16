package repositories

import (
	"context"
	"fmt"
	"go-microservice-product-porto/app/helpers"
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/services"
	"go-microservice-product-porto/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

type ProductRepositoryImpl struct {
	postgresDB      *gorm.DB
	mongoCollection *mongo.Collection
	cacheService    services.CacheService
}

func NewProductRepository(cacheService services.CacheService) ProductRepository {
	return &ProductRepositoryImpl{
		mongoCollection: config.DBConn.MongoDB.Collection("products"),
		postgresDB:      config.DBConn.PostgreDB,
		cacheService:    cacheService,
	}
}

func (r *ProductRepositoryImpl) GetCacheService() services.CacheService {
	return r.cacheService
}

// PostgreSQL operations
func (r *ProductRepositoryImpl) CreateInPostgres(product *models.Product) error {
	return r.postgresDB.Create(product).Error
}

func (r *ProductRepositoryImpl) FindAllInPostgres() ([]models.Product, error) {
	var products []models.Product
	err := r.postgresDB.Find(&products).Error
	return products, err
}

func (r *ProductRepositoryImpl) FindByIDInPostgres(id uint) (*models.Product, error) {
	var product models.Product
	err := r.postgresDB.First(&product, id).Error
	return &product, err
}

// MongoDB operations
func (r *ProductRepositoryImpl) CreateInMongo(product *models.Product) error {
	product.ID = primitive.NewObjectID()
	product.SKU = helpers.GenerateSKU()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.FormattedPrice = helpers.FormatPrice(product.Price)
	product.FormattedCreatedAt = helpers.FormatDateTime(product.CreatedAt)
	product.FormattedUpdatedAt = helpers.FormatDateTime(product.UpdatedAt)

	if !product.IsActive {
		product.IsActive = true
	}

	if product.Stock < 0 {
		product.Stock = 0
	}

	_, err := r.mongoCollection.InsertOne(context.Background(), product)
	if err == nil {
		r.cacheService.DeletePattern("product:list:*")
	}
	return err
}

func (r *ProductRepositoryImpl) FindAllInMongo(page, pageSize int, sortBy, sortDir string, filters map[string]interface{}) ([]models.Product, error) {
	cacheKey := fmt.Sprintf("product:list:p%d:s%d:%s:%s:%v", page, pageSize, sortBy, sortDir, filters)

	if cachedProducts, err := r.cacheService.GetList(cacheKey); err == nil {
		return cachedProducts.Data, nil
	}

	var products []models.Product
	skip := (page - 1) * pageSize

	filterQuery := bson.M{}

	for key, value := range filters {
		switch key {
		case "search":
			filterQuery["$or"] = []bson.M{
				{"name": bson.M{"$regex": value.(string), "$options": "i"}},
				{"description": bson.M{"$regex": value.(string), "$options": "i"}},
				{"sku": bson.M{"$regex": value.(string), "$options": "i"}},
			}
		case "name":
			filterQuery["name"] = bson.M{"$regex": value.(string), "$options": "i"}
		case "sku":
			filterQuery["sku"] = bson.M{"$regex": value.(string), "$options": "i"}
		case "price_min":
			filterQuery["price"] = bson.M{"$gte": value.(float64)}
		case "price_max":
			if _, exists := filterQuery["price"]; exists {
				filterQuery["price"].(bson.M)["$lte"] = value.(float64)
			} else {
				filterQuery["price"] = bson.M{"$lte": value.(float64)}
			}
		case "start_date":
			filterQuery["created_at"] = bson.M{"$gte": value.(time.Time)}
		case "end_date":
			if _, exists := filterQuery["created_at"]; exists {
				filterQuery["created_at"].(bson.M)["$lte"] = value.(time.Time)
			} else {
				filterQuery["created_at"] = bson.M{"$lte": value.(time.Time)}
			}
		case "stock_min":
			filterQuery["stock"] = bson.M{"$gte": value.(int)}
		case "stock_max":
			if _, exists := filterQuery["stock"]; exists {
				filterQuery["stock"].(bson.M)["$lte"] = value.(int)
			} else {
				filterQuery["stock"] = bson.M{"$lte": value.(int)}
			}
		case "is_active":
			filterQuery["is_active"] = value.(bool)
		}
	}

	// Get total items count before pagination
	totalItems, err := r.CountDocuments(filters)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize != 0 {
		totalPages++
	}

	// Map the sortBy field to the correct MongoDB field name
	mongoSortField := "created_at"
	if sortBy != "" {
		switch sortBy {
		case "created_at":
			mongoSortField = "created_at"
		case "updated_at":
			mongoSortField = "updated_at"
		case "price":
			mongoSortField = "price"
		case "name":
			mongoSortField = "name"
		}
	}

	sortValue := 1
	if sortDir == "desc" {
		sortValue = -1
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: mongoSortField, Value: sortValue}})

	cursor, err := r.mongoCollection.Find(context.Background(), filterQuery, opts)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &products)
	if err != nil {
		return nil, err
	}

	for i := range products {
		products[i].FormattedPrice = helpers.FormatPrice(products[i].Price)
		products[i].FormattedCreatedAt = helpers.FormatDateTime(products[i].CreatedAt)
		products[i].FormattedUpdatedAt = helpers.FormatDateTime(products[i].UpdatedAt)
	}

	r.cacheService.SetList(cacheKey, &models.PaginatedResponse{
		Data: products,
		Pagination: models.PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	})

	return products, err
}

func (r *ProductRepositoryImpl) FindByIDInMongo(idString string) (*models.Product, error) {
	if !r.Exists(idString) {
		return nil, fmt.Errorf("product not found")
	}
	product, err := r.cacheService.Get("product:" + idString)
	if err == nil {
		return product, nil
	}

	var mongoProduct models.Product
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"id": objectId}
	err = r.mongoCollection.FindOne(context.Background(), filter).Decode(&mongoProduct)
	if err != nil {
		return nil, err
	}

	mongoProduct.FormattedPrice = helpers.FormatPrice(mongoProduct.Price)
	mongoProduct.FormattedCreatedAt = helpers.FormatDateTime(mongoProduct.CreatedAt)
	mongoProduct.FormattedUpdatedAt = helpers.FormatDateTime(mongoProduct.UpdatedAt)

	r.cacheService.Set("product:"+idString, &mongoProduct)

	return &mongoProduct, nil

}

func (r *ProductRepositoryImpl) UpdateInMongo(idString string, updates map[string]interface{}) error {

	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return err
	}

	updates["updated_at"] = time.Now()

	if price, ok := updates["price"]; ok {
		updates["formatted_price"] = helpers.FormatPrice(price.(float64))
	}

	update := bson.M{
		"$set": updates,
	}

	filter := bson.M{"id": objectId}
	result, err := r.mongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount > 0 {
		// Clear both specific product cache and list cache
		r.cacheService.Delete("product:" + idString)
		r.cacheService.DeletePattern("product:list:*")
	}

	return nil
}

func (r *ProductRepositoryImpl) DeleteInMongo(idString string) error {
	if !r.Exists(idString) {
		return fmt.Errorf("product not found")
	}

	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return err
	}

	filter := bson.M{"id": objectId}
	result, err := r.mongoCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	r.cacheService.Delete("product:" + idString)

	return nil
}

func (r *ProductRepositoryImpl) UpdateStock(id string, newStock int, reason string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	var product models.Product
	filter := bson.M{"id": objectId}
	err = r.mongoCollection.FindOne(context.Background(), filter).Decode(&product)
	if err != nil {
		return err
	}

	movementType := "decrease"
	if newStock > product.Stock {
		movementType = "increase"
	}

	movement := models.StockMovement{
		ID:            primitive.NewObjectID(),
		ProductID:     product.ID.(primitive.ObjectID),
		Type:          movementType,
		Quantity:      newStock - product.Stock,
		PreviousStock: product.Stock,
		NewStock:      newStock,
		Reason:        reason,
		CreatedAt:     time.Now(),
	}

	update := bson.M{
		"$set": bson.M{
			"stock":      newStock,
			"updated_at": time.Now(),
		},
	}

	_, err = r.mongoCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return r.CreateStockMovement(&movement)
}
func (r *ProductRepositoryImpl) CreateStockMovement(movement *models.StockMovement) error {
	_, err := r.mongoCollection.Database().Collection("stock_movements").InsertOne(context.Background(), movement)
	return err
}

func (r *ProductRepositoryImpl) GetStockMovement(productID string) ([]models.StockMovement, error) {
	objectId, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return nil, err
	}

	var movements []models.StockMovement
	filter := bson.M{"product_id": objectId}
	cursor, err := r.mongoCollection.Database().Collection("stock_movements").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &movements)
	return movements, err
}

func (r *ProductRepositoryImpl) CountDocuments(filters map[string]interface{}) (int64, error) {
	// Use the same filter logic as FindAllInMongo
	filterQuery := bson.M{}
	for key, value := range filters {
		switch key {
		case "search":
			filterQuery["$or"] = []bson.M{
				{"name": bson.M{"$regex": value.(string), "$options": "i"}},
				{"description": bson.M{"$regex": value.(string), "$options": "i"}},
				{"sku": bson.M{"$regex": value.(string), "$options": "i"}},
			}
		case "name":
			filterQuery["name"] = bson.M{"$regex": value.(string), "$options": "i"}
		case "sku":
			filterQuery["sku"] = bson.M{"$regex": value.(string), "$options": "i"}
		case "price_min":
			filterQuery["price"] = bson.M{"$gte": value.(float64)}
		case "price_max":
			if _, exists := filterQuery["price"]; exists {
				filterQuery["price"].(bson.M)["$lte"] = value.(float64)
			} else {
				filterQuery["price"] = bson.M{"$lte": value.(float64)}
			}
		case "start_date":
			filterQuery["created_at"] = bson.M{"$gte": value.(time.Time)}
		case "end_date":
			if _, exists := filterQuery["created_at"]; exists {
				filterQuery["created_at"].(bson.M)["$lte"] = value.(time.Time)
			} else {
				filterQuery["created_at"] = bson.M{"$lte": value.(time.Time)}
			}
		case "stock_min":
			filterQuery["stock"] = bson.M{"$gte": value.(int)}
		case "stock_max":
			if _, exists := filterQuery["stock"]; exists {
				filterQuery["stock"].(bson.M)["$lte"] = value.(int)
			} else {
				filterQuery["stock"] = bson.M{"$lte": value.(int)}
			}
		case "is_active":
			filterQuery["is_active"] = value.(bool)
		}
	}

	return r.mongoCollection.CountDocuments(context.Background(), filterQuery)
}

func (r *ProductRepositoryImpl) Exists(idString string) bool {
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return false
	}

	filter := bson.M{"id": objectId}
	count, err := r.mongoCollection.CountDocuments(context.Background(), filter)
	return err == nil && count > 0
}
