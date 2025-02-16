package commands

import (
	eventhandlers "go-microservice-product-porto/internal/application/event_handlers"
	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/internal/infrastructure/cache"
)

type ProductCommandHandler struct {
	repo         product.Repository
	eventHandler *eventhandlers.ProductEventHandler
	cache        cache.CacheService
}

func NewProductCommandHandler(repo product.Repository, eventHandler *eventhandlers.ProductEventHandler, cache cache.CacheService) *ProductCommandHandler {
	return &ProductCommandHandler{
		repo:         repo,
		eventHandler: eventHandler,
		cache:        cache,
	}
}
