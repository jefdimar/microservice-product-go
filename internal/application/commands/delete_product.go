package commands

import (
	"context"
)

type DeleteProductCommand struct {
	ProductID string `json:"product_id"`
}

func (h *ProductCommandHandler) HandleDeleteProduct(ctx context.Context, cmd DeleteProductCommand) error {
	if err := h.repo.Delete(ctx, cmd.ProductID); err != nil {
		return err
	}

	// Invalidate product cache
	h.cache.Delete(cmd.ProductID)

	// Invalidate list cache
	h.cache.Delete("products:all")

	return nil
}
