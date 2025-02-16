package commands

import (
	"context"
	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/pkg/errors"
	"time"
)

type UpdateStockCommand struct {
	ProductID string `json:"product_id"`
	Stock     int    `json:"stock"`
}

func (h *ProductCommandHandler) HandleUpdateStock(ctx context.Context, cmd UpdateStockCommand) error {
	prod, err := h.repo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return errors.StandardError(errors.ENOTFOUND, err)
	}

	if cmd.Stock < 0 {
		return errors.StandardError(errors.EVALIDATION, product.ErrInvalidStock)
	}

	oldStock := prod.Stock
	prod.Stock = cmd.Stock
	prod.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, prod); err != nil {
		return errors.StandardError(errors.EREPOSITORY, err)
	}

	// Handle cache update
	if err := h.cache.Set(prod.ID.Hex(), prod); err != nil {
		return errors.StandardError(errors.ECACHE, err)
	}

	h.eventHandler.HandleStockUpdated(&product.ProductStockUpdatedEvent{
		Product:  prod,
		OldStock: oldStock,
		NewStock: cmd.Stock,
	})

	return nil
}
