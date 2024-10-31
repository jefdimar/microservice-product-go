package repositories

import (
	"context"
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type ProductRepository struct {
	postgresDB *gorm.DB
	mongoCollection *mongo.Collection
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		mongoCollection: config.DBConn.MongoDB.Collection("products"),
		postgresDB: config.DBConn.PostgreDB,
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
	cursor, err := r.mongoCollection.Find(context.Background(), bson.M{})
	if err != nil {
			return nil, err
	}
	err = cursor.All(context.Background(), &products)
	return products, err
}

func (r *ProductRepository) FindByIDInMongo(id primitive.ObjectID) (*models.Product, error) {
	var product models.Product
	err := r.mongoCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&product)
	return &product, err
}