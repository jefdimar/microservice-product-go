package usecase

import "go-microservice-product-porto/app/models"

type ProductUsecase interface {
	GetAllProducts(page, pageSize int, sortBy, sortDir string, filters map[string]interface{}) (*models.PaginatedResponse, error)
	CreateProduct(product *models.Product) error
	GetProductByID(id string) (*models.Product, error)
	UpdateProduct(id string, product *models.Product) error
	DeleteProduct(id string) error
	InvalidateRelatedCaches(productID string) error
	InvalidateListCaches() error
}
