package queries

import (
	"context"

	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/pkg/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetProductQuery struct {
	ID string `json:"id"`
}

func (h *ProductQueryHandler) HandleGetProduct(ctx context.Context, query GetProductQuery) (*product.Product, error) {
	objectID, err := primitive.ObjectIDFromHex(query.ID)
	if err != nil {
		return nil, errors.StandardError(errors.EINVALID, err)
	}

	// Try to get from cache first
	cachedProduct, err := h.cache.Get(objectID.Hex())
	if err != nil {
		// Check list cache
		listCacheKey := "products_list_p1_s10_name_asc" // Default list cache key
		cachedList, err := h.cache.Get(listCacheKey)
		if err == nil && cachedList != nil {
			if listResponse, ok := cachedList.(*ListProductsResponse); ok {
				// Search for product in cached list
				for _, p := range listResponse.Products {
					if p.ID == objectID {
						// Found in list cache, store in individual cache
						h.cache.Set(objectID.Hex(), p)
						return p, nil
					}
				}
			}
		}
	}

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
		if err == err {
			return nil, errors.StandardError(errors.ENOTFOUND, err)
		}
		return nil, errors.StandardError(errors.EREPOSITORY, err)
	}

	// Update cache
	if err := h.cache.Set(objectID.Hex(), product); err != nil {
		return nil, errors.StandardError(errors.ECACHE, err)
	}

	return product, nil
}
