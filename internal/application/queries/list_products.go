package queries

import (
	"context"

	"fmt"
	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/pkg/errors"
)

type ListProductsQuery struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	SortBy   string `json:"sort_by"`
	SortDir  string `json:"sort_dir"` // "asc" or "desc"
}

type ListProductsResponse struct {
	Products []*product.Product `json:"products"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

func (h *ProductQueryHandler) HandleListProducts(ctx context.Context, query ListProductsQuery) (*ListProductsResponse, error) {
	// Set default values if not provided
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}

	// Validate sort direction
	if query.SortDir != "" && query.SortDir != "asc" && query.SortDir != "desc" {
		query.SortDir = "asc"
	}

	// Generate cache key based on query parameters
	cacheKey := fmt.Sprintf("products_list_p%d_s%d_%s_%s", query.Page, query.PageSize, query.SortBy, query.SortDir)

	// Try to get from cache first
	cachedData, err := h.cache.Get(cacheKey)
	if err == nil && cachedData != nil {
		if response, ok := cachedData.(*ListProductsResponse); ok {
			return response, nil
		}
	}

	// Get from repository if not in cache
	products, total, err := h.repo.FindAll(ctx, query.Page, query.PageSize, query.SortBy, query.SortDir)
	if err != nil {
		return nil, errors.StandardError(errors.EREPOSITORY, err)
	}

	response := &ListProductsResponse{
		Products: products,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	// Store in cache
	if err := h.cache.Set(cacheKey, response); err != nil {
		return nil, errors.StandardError(errors.ECACHE, err)
	}

	return response, nil
}
