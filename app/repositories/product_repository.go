package repositories

import (
	"context"
	"go-microservice-product-porto/app/models"
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
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		mongoCollection: config.DBConn.MongoDB.Collection("products"),
		postgresDB:      config.DBConn.PostgreDB,
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
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	_, err := r.mongoCollection.InsertOne(context.Background(), product)
	return err
}

func (r *ProductRepository) FindAllInMongo() ([]models.Product, error) {
	var products []models.Product

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.mongoCollection.Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return nil, err
	}

	err = cursor.All(context.Background(), &products)
	return products, err
}

func (r *ProductRepository) FindByIDInMongo(idString string) (*models.Product, error) {
	var product models.Product
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"id": objectId}
	err = r.mongoCollection.FindOne(context.Background(), filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) UpdateInMongo(idString string, product *models.Product) error {
	objectId, err := primitive.ObjectIDFromHex(idString)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"price":       product.Price,
			"description": product.Description,
			"updated_at":  time.Now(),
		},
	}

	filter := bson.M{"id": objectId}
	_, err = r.mongoCollection.UpdateOne(context.Background(), filter, update)
	return err
}
