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
	"go-microservice-product-porto/pkg/logger"
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
	logger.Debug().
		Str("product_name", prod.Name).
		Float64("price", prod.Price).
		Msg("attempting to create product")

	_, err := r.collection.InsertOne(ctx, prod)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			logger.Error().
				Str("product_name", prod.Name).
				Err(err).
				Msg("product already exists")
			return errors.StandardError(errors.ECONFLICT, product.ErrProductAlreadyExists)
		}
		logger.Error().
			Str("product_name", prod.Name).
			Err(err).
			Msg("failed to create product")
		return errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to create product: %v", err))
	}
	logger.Info().
		Str("product_name", prod.Name).
		Msg("product created successfully")

	return nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*product.Product, error) {
	logger.Debug().
		Str("product_id", id).
		Msg("attempting to find product by ID")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error().
			Str("product_id", id).
			Err(err).
			Msg("invalid product ID")
		return nil, errors.StandardError(errors.EINVALID, fmt.Errorf("invalid product ID: %v", err))
	}

	var prod product.Product
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&prod)
	if err == mongo.ErrNoDocuments {
		logger.Error().
			Str("product_id", id).
			Err(err).
			Msg("product not found")
		return nil, errors.StandardError(errors.ENOTFOUND, product.ErrProductNotFound)
	}
	if err != nil {
		logger.Error().
			Str("product_id", id).
			Err(err).
			Msg("failed to find product")
		return nil, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to find product: %v", err))
	}
	logger.Info().
		Str("product_id", id).
		Msg("product found successfully")
	return &prod, nil
}

func (r *ProductRepository) FindAll(ctx context.Context, page, pageSize int, sortBy, sortDir string) ([]*product.Product, int64, error) {
	logger.Debug().
		Int("page", page).
		Int("page_size", pageSize).
		Str("sort_by", sortBy).
		Str("sort_dir", sortDir).
		Msg("attempting to find all products with all filters")

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
		logger.Error().
			Err(err).
			Msg("failed to find products")
		return nil, 0, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to find products: %v", err))
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err = cursor.All(ctx, &products); err != nil {
		logger.Error().
			Err(err).
			Msg("failed to decode products")
		return nil, 0, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to decode products: %v", err))
	}

	total, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to count products")
		return nil, 0, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to count products: %v", err))
	}

	logger.Info().
		Int64("total", total).
		Msg("products found successfully with total products is :")
	return products, total, nil
}

func (r *ProductRepository) Update(ctx context.Context, prod *product.Product) error {
	logger.Debug().
		Str("product_id", prod.ID.Hex()).
		Msg("attempting to update product")

	result, err := r.collection.ReplaceOne(ctx, bson.M{"_id": prod.ID}, prod)
	if err != nil {
		logger.Error().
			Str("product_id", prod.ID.Hex()).
			Err(err).
			Msg("failed to update product")
		return errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to update product: %v", err))
	}
	if result.MatchedCount == 0 {
		logger.Error().
			Str("product_id", prod.ID.Hex()).
			Msg("product not found")
		return errors.StandardError(errors.ENOTFOUND, product.ErrProductNotFound)
	}
	logger.Info().
		Str("product_id", prod.ID.Hex()).
		Msg("product updated successfully")
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	logger.Debug().
		Str("product_id", id).
		Msg("attempting to delete product")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error().
			Str("product_id", id).
			Err(err).
			Msg("invalid product ID")
		return errors.StandardError(errors.EINVALID, fmt.Errorf("invalid product ID: %v", err))
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		logger.Error().
			Str("product_id", id).
			Err(err).
			Msg("failed to delete product")
		return errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to delete product: %v", err))
	}
	if result.DeletedCount == 0 {
		logger.Error().
			Str("product_id", id).
			Msg("product not found")
		return errors.StandardError(errors.ENOTFOUND, product.ErrProductNotFound)
	}
	logger.Info().
		Str("product_id", id).
		Msg("product deleted successfully")

	return nil
}

func (r *ProductRepository) Search(ctx context.Context, name string, minPrice, maxPrice float64) ([]*product.Product, error) {
	logger.Debug().
		Str("name", name).
		Float64("min_price", minPrice).
		Float64("max_price", maxPrice).
		Msg("attempting to search products with parameters")

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
		logger.Error().
			Err(err).
			Msg("failed to search products")

		return nil, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to search products: %v", err))
	}
	defer cursor.Close(ctx)

	var products []*product.Product
	if err := cursor.All(ctx, &products); err != nil {
		logger.Error().
			Err(err).
			Msg("failed to decode search results")
		return nil, errors.StandardError(errors.EREPOSITORY, fmt.Errorf("failed to decode search results: %v", err))
	}

	logger.Info().
		Int("total", len(products)).
		Msg("products found successfully")
	return products, nil
}
