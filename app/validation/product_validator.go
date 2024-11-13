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

type ValidationErrors struct {
	Errors map[string]string
}

func (ve *ValidationErrors) AddError(field, message string) {
	if ve.Errors == nil {
		ve.Errors = make(map[string]string)
	}
	ve.Errors[field] = message
}

func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// NewProductValidator creates a new instance of ProductValidator
func NewProductValidator(product *models.Product) *ProductValidator {
	return &ProductValidator{
		Product: product,
	}
}

// Validate runs all validation rules for the Product
func (v *ProductValidator) Validate() error {
	validationErrors := &ValidationErrors{}

	v.validateName(validationErrors)
	v.validatePrice(validationErrors)
	v.validateDescription(validationErrors)
	v.validateStock(validationErrors)

	if validationErrors.HasErrors() {
		return validationErrors
	}
	return nil
}

// validateName check if the product name meets the requirements
func (v *ProductValidator) validateName(ve *ValidationErrors) {
	if v.Product.Name == "" {
		ve.AddError("name", "name is required")
	} else if len(v.Product.Name) < 3 {
		ve.AddError("name", "name must be at least 3 characters long")
	} else if len(v.Product.Name) > 100 {
		ve.AddError("name", "name must be at most 100 characters long")
	}
}

// validatePrice checks if the price is valid
func (v *ProductValidator) validatePrice(ve *ValidationErrors) {
	if v.Product.Price < 0 {
		ve.AddError("price", "price cannot be negative")
	} else if v.Product.Price > 1000000000 {
		ve.AddError("price", "price cannot be greater than 1 billion")
	}
}

// validateStock checks if the stock quantity is valid
func (v *ProductValidator) validateStock(ve *ValidationErrors) {
	if v.Product.Stock < 0 {
		ve.AddError("stock", "stock cannot be negative")
	} else if v.Product.Stock > 999999 {
		ve.AddError("stock", "stock cannot be greater than 999,999 units")
	}
}

// validateDescription checks if the product description is valid
func (v *ProductValidator) validateDescription(ve *ValidationErrors) {
	if v.Product.Description == "" {
		ve.AddError("description", "description is required")
	} else {
		description := strings.TrimSpace(v.Product.Description)
		if len(description) < 10 {
			ve.AddError("description", "description must be at least 10 characters long")
		} else if len(description) > 3000 {
			ve.AddError("description", "description must be at most 3000 characters long")
		}
	}
}

func (ve *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %v", ve.Errors)
}
