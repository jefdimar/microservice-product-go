package queries

import (
	"context"

	"go-microservice-product-porto/internal/domain/product"
)

type ListProductsQuery struct {
	// Add any filtering/pagination parameters here
}

func (h *ProductQueryHandler) HandleListProducts(ctx context.Context, query ListProductsQuery) ([]*product.Product, error) {
	cacheKey := "products:all"

	// Try to get from cache first
	if cachedProducts, err := h.cache.Get(cacheKey); err == nil {
		if products, ok := cachedProducts.([]*product.Product); ok {
			return products, nil
		}
	}

	// Get from repository if not in cache
	products, err := h.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Update cache
	h.cache.Set(cacheKey, products)

	return products, nil
}
