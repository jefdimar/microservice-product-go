package mongodb

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"go-microservice-product-porto/internal/domain/product"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(client *mongo.Client) *ProductRepository {
	collection := client.Database("products_db").Collection("products")
	return &ProductRepository{
		collection: collection,
	}
}

func (r *ProductRepository) Create(ctx context.Context, prod *product.Product) error {
	_, err := r.collection.InsertOne(ctx, prod)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return product.ErrProductAlreadyExists
		}
		return err
	}
	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*product.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var prod product.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&prod)
	if err == mongo.ErrNoDocuments {
		return nil, product.ErrProductNotFound
	}
	return &prod, err
}

func (r *ProductRepository) FindAll(ctx context.Context) ([]*product.Product, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Update(ctx context.Context, prod *product.Product) error {
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": prod.ID}, prod)
	return err
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}

func (r *ProductRepository) Search(ctx context.Context, name string, minPrice, maxPrice float64) ([]*product.Product, error) {
	matchStage := bson.M{}

	if name != "" {
		matchStage["name"] = bson.M{
			"$regex":   name,
			"$options": "i",
		}
	}

	if minPrice > 0 || maxPrice > 0 {
		priceMatch := bson.M{}
		if minPrice > 0 {
			priceMatch["$gte"] = minPrice
		}
		if maxPrice > 0 {
			priceMatch["$lte"] = maxPrice
		}
		matchStage["price"] = priceMatch
	}

	pipeline := []bson.M{
		{"$match": matchStage},
	}

	log.Printf("Pipeline: %+v", pipeline)

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}

	return products, nil
}
