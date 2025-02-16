package commands

import (
	"context"
	"go-microservice-product-porto/internal/domain/product"
	"go-microservice-product-porto/pkg/errors"
)

type CreateProductCommand struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

func (h *ProductCommandHandler) HandleCreateProduct(ctx context.Context, cmd CreateProductCommand) error {
	newProduct := product.NewProduct(cmd.Name, cmd.Description, cmd.Price, cmd.Stock)

	if !newProduct.IsValid() {
		return errors.StandardError(errors.EVALIDATION, product.ErrInvalidProduct)
	}

	if err := h.repo.Create(ctx, newProduct); err != nil {
		if err == product.ErrProductAlreadyExists {
			return errors.StandardError(errors.ECONFLICT, err)
		}
		return errors.StandardError(errors.EREPOSITORY, err)
	}

	// Update single product cache
	if err := h.cache.Set(newProduct.ID.Hex(), newProduct); err != nil {
		return errors.StandardError(errors.ECACHE, err)
	}

	// Invalidate list cache to ensure fresh data
	if err := h.cache.Delete("products:all"); err != nil {
		return errors.StandardError(errors.ECACHE, err)
	}

	return nil
}
