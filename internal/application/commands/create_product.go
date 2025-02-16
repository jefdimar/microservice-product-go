package commands

import (
	"context"

	"go-microservice-product-porto/internal/domain/product"
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
		return product.ErrInvalidProduct
	}

	if err := h.repo.Create(ctx, newProduct); err != nil {
		if err == product.ErrProductAlreadyExists {
			return err
		}
		return err
	}

	// Update single product cache
	h.cache.Set(newProduct.ID.Hex(), newProduct)

	// Invalidate list cache to ensure fresh data
	h.cache.Delete("products:all")

	return nil
}
