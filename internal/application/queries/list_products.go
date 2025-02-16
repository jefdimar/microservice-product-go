package queries

import (
	"context"

	"go-microservice-product-porto/internal/domain/product"
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

	products, total, err := h.repo.FindAll(ctx, query.Page, query.PageSize, query.SortBy, query.SortDir)
	if err != nil {
		return nil, err
	}

	response := &ListProductsResponse{
		Products: products,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}

	return response, nil
}
