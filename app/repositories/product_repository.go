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

func (r *ProductRepository) FindAllInMongo(page, pageSize int, sortBy, sortDir string) ([]models.Product, error) {
	var products []models.Product

	skip := (page - 1) * pageSize

	sortValue := 1
	if sortDir == "desc" {
		sortValue = -1
	}

	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: sortBy, Value: sortValue}})

	// opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.mongoCollection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &products)
	if err != nil {
		return nil, err
	}

	for i := range products {
		products[i].FormattedPrice = helpers.FormatPrice(products[i].Price)
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
