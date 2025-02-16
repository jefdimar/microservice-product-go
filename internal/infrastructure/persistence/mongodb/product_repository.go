package mongodb

import (
	"context"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

func (r *ProductRepository) FindAll(ctx context.Context, page, pageSize int, sortBy, sortDir string) ([]*product.Product, int64, error) {
	skip := (page - 1) * pageSize

	// Map API sort fields to actual MongoDB document field names
	sortFieldMap := map[string]string{
		"name":       "name",       // Verify this matches your MongoDB field
		"price":      "price",      // For nested fields
		"stock":      "stock",      // Verify this matches your MongoDB field
		"created_at": "created_at", // Verify this matches your MongoDB field
	}

	// Build sort options
	sortOpts := bson.D{}
	if sortBy != "" {
		if mongoField, exists := sortFieldMap[sortBy]; exists {
			sortValue := 1
			if strings.ToLower(sortDir) == "desc" {
				sortValue = -1
			}

			// Special handling for price field
			if sortBy == "price" {
				sortOpts = bson.D{{Key: "price.amount", Value: sortValue}}
			} else {
				sortOpts = bson.D{{Key: mongoField, Value: sortValue}}
			}
		}
	} else {
		// Default sort
		sortOpts = bson.D{{Key: "_id", Value: 1}}
	}

	// Add debug logging
	log.Printf("Applying sort options: %+v", sortOpts)

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(sortOpts)

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Printf("MongoDB Find error: %v", err)
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err = cursor.All(ctx, &products); err != nil {
		log.Printf("Cursor decode error: %v", err)
		return nil, 0, err
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
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
