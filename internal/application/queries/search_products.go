package queries

import (
	"context"
	"go-microservice-product-porto/internal/domain/product"
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
	// Skip cache for search queries to ensure fresh results
	products, err := h.repo.Search(ctx, query.Name, query.MinPrice, query.MaxPrice)
	if err != nil {
		return nil, err
	}

	return products, nil
}
