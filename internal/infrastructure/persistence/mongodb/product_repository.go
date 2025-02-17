package mongodb

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/pkg/errors"
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
			return errors.StandardError(errors.ECONFLICT, product.ErrProductAlreadyExists)
		}
		return errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to create product: %v", err))
	}
	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*product.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.StandardError(errors.EINVALID, fmt.Errorf("invalid product ID: %v", err))
	}

	var prod product.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&prod)
	if err == mongo.ErrNoDocuments {
		return nil, errors.StandardError(errors.ENOTFOUND, product.ErrProductNotFound)
	}
	if err != nil {
		return nil, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to find product: %v", err))
	}
	return &prod, nil
}

func (r *ProductRepository) FindAll(ctx context.Context, page, pageSize int, sortBy, sortDir string) ([]*product.Product, int64, error) {
	skip := (page - 1) * pageSize

	sortFieldMap := map[string]string{
		"name":       "name",
		"price":      "price",
		"stock":      "stock",
		"created_at": "created_at",
	}

	sortOpts := bson.D{}
	if sortBy != "" {
		if mongoField, exists := sortFieldMap[sortBy]; exists {
			sortValue := 1
			if strings.ToLower(sortDir) == "desc" {
				sortValue = -1
			}
			if sortBy == "price" {
				sortOpts = bson.D{{Key: "price.amount", Value: sortValue}}
			} else {
				sortOpts = bson.D{{Key: mongoField, Value: sortValue}}
			}
		}
	} else {
		sortOpts = bson.D{{Key: "_id", Value: 1}}
	}

	findOptions := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(sortOpts)

	cursor, err := r.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, 0, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to find products: %v", err))
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, 0, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to decode products: %v", err))
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, 0, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to count products: %v", err))
	}

	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, prod *product.Product) error {
	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": prod.ID}, prod)
	if err != nil {
		return errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to update product: %v", err))
	}
	if result.MatchedCount == 0 {
		return errors.StandardError(errors.ENOTFOUND, product.ErrProductNotFound)
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.StandardError(errors.EINVALID, fmt.Errorf("invalid product ID: %v", err))
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to delete product: %v", err))
	}
	if result.DeletedCount == 0 {
		return errors.StandardError(errors.ENOTFOUND, product.ErrProductNotFound)
	}
	return nil
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

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to search products: %v", err))
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to decode search results: %v", err))
	}

	return products, nil
}
