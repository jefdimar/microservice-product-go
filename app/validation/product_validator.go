package validation

import (
	"fmt"
	"go-microservice-product-porto/app/models"
	"strings"
)

// ProductValidator handles validation rules for Product model
type ProductValidator struct {
	Product *models.Product
}

// Validate runs all validation rules for the Product
func (v *ProductValidator) Validate() error {
	if err := v.validateName(); err != nil {
		return err
	}
	if err := v.validatePrice(); err != nil {
		return err
	}
	if err := v.validateDescription(); err != nil {
		return err
	}
	if err := v.validateStock(); err != nil {
		return err
	}
	return nil
}

// NewProductValidator creates a new instance of ProductValidator
func NewProductValidator(product *models.Product) *ProductValidator {
	return &ProductValidator{
		Product: product,
	}
}

// validateName check if the product name meets the requirements
func (v *ProductValidator) validateName() error {
	if v.Product.Name == "" {
		return fmt.Errorf("name is required")
	}

	if len(v.Product.Name) < 3 {
		return fmt.Errorf("name must be at least 3 characters long")
	}

	if len(v.Product.Name) > 100 {
		return fmt.Errorf("name must be at most 100 characters long")
	}

	return nil
}

// validatePrice checks if the price is valid
func (v *ProductValidator) validatePrice() error {
	if v.Product.Price < 0 {
		return fmt.Errorf("price cannot be negative")
	}

	if v.Product.Price > 1000000000 {
		return fmt.Errorf("price cannot be greater than 1 billion")
	}

	return nil
}

// validateStock checks if the stock quantity is valid
func (v *ProductValidator) validateStock() error {
	if v.Product.Stock < 0 {
		return fmt.Errorf("stock cannot be negative")
	}

	if v.Product.Stock > 999999 {
		return fmt.Errorf("stock cannot be greater than 999,999 units")
	}

	return nil
}

// validateDescription checks if the product description is valid
func (v *ProductValidator) validateDescription() error {
	if v.Product.Description == "" {
		return fmt.Errorf("description is required")
	}

	description := strings.TrimSpace(v.Product.Description)
	if len(description) < 10 {
		return fmt.Errorf("description must be at least 10 characters long")
	}
	if len(description) > 3000 {
		return fmt.Errorf("description must be at most 3000 characters long")
	}

	return nil
}
