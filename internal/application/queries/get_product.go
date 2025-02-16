package queries

import (
	"context"

	"go-microservice-product-porto/internal/domain/product"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetProductQuery struct {
	ID string `json:"id"`
}

func (h *ProductQueryHandler) HandleGetProduct(ctx context.Context, query GetProductQuery) (*product.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, err
	}

	// Try to get from cache first
	cachedProduct, _ := h.cache.Get(objectID.Hex())
	if productMap, ok := cachedProduct.(map[string]interface{}); ok {
		// Convert map to Product struct
		product := &product.Product{
			ID:          objectID,
			Name:        productMap["name"].(string),
			Description: productMap["description"].(string),
			Price:       productMap["price"].(float64),
			Stock:       int(productMap["stock"].(float64)),
		}
		return product, nil
	}

	// Get from repository if not in cache
	product, err := h.repo.FindByID(ctx, objectID.Hex())
	if err != nil {
		return nil, err
	}

	// Update cache
	h.cache.Set(objectID.Hex(), product)

	return product, nil
}
