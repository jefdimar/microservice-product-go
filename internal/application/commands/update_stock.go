package commands

import (
	"context"
	"time"

	"go-microservice-product-porto/internal/domain/product"
)

type UpdateStockCommand struct {
	ProductID string `json:"product_id"`
	Stock     int    `json:"stock"`
}

func (h *ProductCommandHandler) HandleUpdateStock(ctx context.Context, cmd UpdateStockCommand) error {
	prod, err := h.repo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	if cmd.Stock < 0 {
		return product.ErrInvalidStock
	}

	oldStock := prod.Stock
	prod.Stock = cmd.Stock
	prod.UpdatedAt = time.Now()

	if err := h.repo.Update(ctx, prod); err != nil {
		return err
	}

	h.eventHandler.HandleStockUpdated(&product.ProductStockUpdatedEvent{
		Product:  prod,
		OldStock: oldStock,
		NewStock: cmd.Stock,
	})

	return nil
}
