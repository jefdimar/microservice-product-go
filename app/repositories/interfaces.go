package repositories

import (
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/services"
)

type ProductRepository interface {
	CreateInPostgres(product *models.Product) error
	FindAllInPostgres() ([]models.Product, error)
	FindByIDInPostgres(id uint) (*models.Product, error)
	CreateInMongo(product *models.Product) error
	UpdateInMongo(idString string, updates map[string]interface{}) error
	DeleteInMongo(idString string) error
	FindAllInMongo(page, pageSize int, sortBy, sortDir string, filters map[string]interface{}) ([]models.Product, error)
	FindByIDInMongo(idString string) (*models.Product, error)
	CountDocuments(filters map[string]interface{}) (int64, error)
	GetCacheService() services.CacheService
	UpdateStock(id string, newStock int, reason string) error
	CreateStockMovement(movement *models.StockMovement) error
	GetStockMovement(productID string) ([]models.StockMovement, error)
}
