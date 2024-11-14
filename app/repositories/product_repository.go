package repositories

import (
	"context"
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

type ProductRepository struct {
	postgresDB      *gorm.DB
	mongoCollection *mongo.Collection
	cacheService    *services.CacheService
}

func NewProductRepository() *ProductRepository {
	cacheService := services.NewCacheService(config.DBConn.Redis)
	return &ProductRepository{
		mongoCollection: config.DBConn.MongoDB.Collection("products"),
		postgresDB:      config.DBConn.PostgreDB,
		cacheService:    cacheService,
	}
}

// PostgreSQL operations
func (r *ProductRepository) CreateInPostgres(product *models.Product) error {
	return r.postgresDB.Create(product).Error
}

func (r *ProductRepository) FindAllInPostgres() ([]models.Product, error) {
	var products []models.Product
	err := r.postgresDB.Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByIDInPostgres(id uint) (*models.Product, error) {
	var product models.Product
	err := r.postgresDB.First(&product, id).Error
	return &product, err
}

// MongoDB operations
func (r *ProductRepository) CreateInMongo(product *models.Product) error {
	product.ID = primitive.NewObjectID()
	product.SKU = helpers.GenerateSKU()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.FormattedPrice = helpers.FormatPrice(product.Price)

	if !product.IsActive {
		product.IsActive = true
	}

	if product.Stock < 0 {
		product.Stock = 0
	}

	_, err := r.mongoCollection.InsertOne(context.Background(), product)
	return err
}

func (r *ProductRepository) FindAllInMongo(page, pageSize int, sortBy, sortDir string, filters map[string]interface{}) ([]models.Product, error) {
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

	sortValue := 1
	if sortDir == "desc" {
		sortValue = -1
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: sortBy, Value: sortValue}})
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

	return products, err
}

func (r *ProductRepository) FindByIDInMongo(idString string) (*models.Product, error) {
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

func (r *ProductRepository) UpdateInMongo(idString string, product *models.Product) error {
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":            product.Name,
			"price":           product.Price,
			"description":     product.Description,
			"formatted_price": helpers.FormatPrice(product.Price),
			"updated_at":      time.Now(),
		},
	}

	filter := bson.M{"id": objectId}
	_, err = r.mongoCollection.UpdateOne(context.Background(), filter, update)
	return err
}

func (r *ProductRepository) DeleteInMongo(idString string) error {
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

func (r *ProductRepository) CountDocuments(filters map[string]interface{}) (int64, error) {
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
