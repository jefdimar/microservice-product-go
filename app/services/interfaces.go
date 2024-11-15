package services

import "go-microservice-product-porto/app/models"

type CacheService interface {
	Get(key string) (*models.Product, error)
	Set(key string, product *models.Product) error
	Delete(key string) error
	GetList(key string) (*models.PaginatedResponse, error)
	SetList(key string, value *models.PaginatedResponse) error
	DeletePattern(pattern string) error
	InvalidateRelatedCaches(productID string) error
}
