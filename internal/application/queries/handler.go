package queries

import (
	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/internal/infrastructure/cache"
)

type ProductQueryHandler struct {
	repo  product.Repository
	cache cache.CacheService
}

func NewProductQueryHandler(repo product.Repository, cache cache.CacheService) *ProductQueryHandler {
	return &ProductQueryHandler{
		repo:  repo,
		cache: cache,
	}
}
