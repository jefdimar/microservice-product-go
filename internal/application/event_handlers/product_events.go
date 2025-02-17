package eventhandlers

import (
	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/internal/infrastructure/cache"
	"go-microservice-product-porto/pkg/errors"
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
	if err := h.cache.Delete("products_list"); err != nil {
		log.Printf("Error deleting products_list from cache: %v", errors.StandardError(errors.ECACHE, err))
	}
}

func (h *ProductEventHandler) HandleStockUpdated(event *product.ProductStockUpdatedEvent) {
	if err := h.cache.Set(event.Product.ID.Hex(), event.Product); err != nil {
		log.Printf("Error updating cache: %v", errors.StandardError(errors.ECACHE, err))
		return
	}

	log.Printf("Stock updated for product %s from %d to %d",
		event.Product.ID.Hex(), event.OldStock, event.NewStock)
}
func (h *ProductEventHandler) HandleProductDeleted(event *product.ProductDeletedEvent) {
	if err := h.cache.Delete(event.ProductID); err != nil {
		log.Printf("Error deleting product from cache: %v", errors.StandardError(errors.ECACHE, err))
	}

	if err := h.cache.Delete("products_list"); err != nil {
		log.Printf("Error deleting products_list from cache: %v", errors.StandardError(errors.ECACHE, err))
	}
}
