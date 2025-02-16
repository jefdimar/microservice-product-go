package eventhandlers

import (
	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/internal/infrastructure/cache"
	"log"
)

type ProductEventHandler struct {
	cache cache.CacheService
	repo  product.Repository
}

func NewProductEventHandler(cache cache.CacheService, repo product.Repository) *ProductEventHandler {
	return &ProductEventHandler{
		cache: cache,
		repo:  repo,
	}
}

func (h *ProductEventHandler) HandleProductCreated(event *product.ProductCreatedEvent) {
	// Clear products list cache
	h.cache.Delete("products_list")
}

func (h *ProductEventHandler) HandleStockUpdated(event *product.ProductStockUpdatedEvent) {
	// Update product cache
	h.cache.Set(event.Product.ID.Hex(), event.Product)

	// Additional logging or monitoring could be added here
	log.Printf("Stock updated for product %s from %d to %d",
		event.Product.ID.Hex(), event.OldStock, event.NewStock)
}
func (h *ProductEventHandler) HandleProductDeleted(event *product.ProductDeletedEvent) {
	// Remove from cache
	h.cache.Delete(event.ProductID)
	h.cache.Delete("products_list")
}
