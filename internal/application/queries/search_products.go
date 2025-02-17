package queries

import (
	"context"
	"fmt"
	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/pkg/errors"
)

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type SearchProductsQuery struct {
	Name       string     `json:"name"`
	MinPrice   float64    `json:"min_price"`
	MaxPrice   float64    `json:"max_price"`
	Pagination Pagination `json:"pagination"`
}

func (h *ProductQueryHandler) HandleSearchProducts(ctx context.Context, query SearchProductsQuery) ([]*product.Product, error) {
	// Generate cache key based on search parameters
	cacheKey := fmt.Sprintf("search_products_%s_%.2f_%.2f", query.Name, query.MinPrice, query.MaxPrice)

	// Try to get from cache first
	cachedResults, err := h.cache.Get(cacheKey)
	if err == nil && cachedResults != nil {
		if products, ok := cachedResults.([]*product.Product); ok {
			return products, nil
		}
	}

	// Perform search in repository
	products, err := h.repo.Search(ctx, query.Name, query.MinPrice, query.MaxPrice)
	if err != nil {
		return nil, errors.StandardError(errors.EREPOSITORY, err)
	}

	// Store results in cache
	if err := h.cache.Set(cacheKey, products); err != nil {
		return nil, errors.StandardError(errors.ECACHE, err)
	}

	return products, nil
}
