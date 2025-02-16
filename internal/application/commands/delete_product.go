package commands

import (
	"context"
	"go-microservice-product-porto/pkg/errors"
)

type DeleteProductCommand struct {
	ProductID string `json:"product_id"`
}

func (h *ProductCommandHandler) HandleDeleteProduct(ctx context.Context, cmd DeleteProductCommand) error {
	// Check if product exists before deletion
	if _, err := h.repo.FindByID(ctx, cmd.ProductID); err != nil {
		return errors.StandardError(errors.ENOTFOUND, err)
	}

	if err := h.repo.Delete(ctx, cmd.ProductID); err != nil {
		return errors.StandardError(errors.EREPOSITORY, err)
	}

	// Invalidate product cache
	if err := h.cache.Delete(cmd.ProductID); err != nil {
		return errors.StandardError(errors.ECACHE, err)
	}

	// Invalidate list cache
	if err := h.cache.Delete("products:all"); err != nil {
		return errors.StandardError(errors.ECACHE, err)
	}

	return nil
}
