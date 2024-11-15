package usecase

import (
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/repositories"
	"go-microservice-product-porto/app/services"
	"go-microservice-product-porto/app/validation"
	"math"
)

type ProductUsecaseImpl struct {
	repository   repositories.ProductRepository
	cacheService services.CacheService
}

func NewProductUsecase(repo repositories.ProductRepository, cache services.CacheService) ProductUsecase {
	return &ProductUsecaseImpl{
		repository:   repo,
		cacheService: cache,
	}
}

func (b *ProductUsecaseImpl) GetAllProducts(page, pageSize int, sortBy, sortDir string, filters map[string]interface{}) (*models.PaginatedResponse, error) {
	products, err := b.repository.FindAllInMongo(page, pageSize, sortBy, sortDir, filters)
	if err != nil {
		return nil, err
	}

	totalItems, err := b.repository.CountDocuments(filters)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return &models.PaginatedResponse{
		Data: products,
		Pagination: models.PaginationMeta{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalItems:  totalItems,
			TotalPages:  totalPages,
		},
	}, nil
}

func (b *ProductUsecaseImpl) CreateProduct(product *models.Product) error {
	validator := validation.NewProductValidator(product)
	if err := validator.Validate(); err != nil {
		return err
	}

	err := b.repository.CreateInMongo(product)
	if err == nil {
		b.repository.GetCacheService().DeletePattern("products:list:*")
	}
	return err
}

func (b *ProductUsecaseImpl) GetProductByID(id string) (*models.Product, error) {
	return b.repository.FindByIDInMongo(id)
}
func (b *ProductUsecaseImpl) UpdateProduct(id string, product *models.Product) error {
	validator := validation.NewProductValidator(product)
	if err := validator.Validate(); err != nil {
		return err
	}

	err := b.repository.UpdateInMongo(id, product)
	if err == nil {
		b.repository.GetCacheService().Delete("product:" + id)
		b.repository.GetCacheService().DeletePattern("products:list:*")
	}
	return err
}
func (b *ProductUsecaseImpl) DeleteProduct(id string) error {
	err := b.repository.DeleteInMongo(id)
	if err == nil {
		b.repository.GetCacheService().Delete("product:" + id)
		b.repository.GetCacheService().DeletePattern("products:list:*")
	}
	return err
}
func (u *ProductUsecaseImpl) InvalidateRelatedCaches(productID string) error {
	return u.cacheService.InvalidateRelatedCaches(productID)
}
func (u *ProductUsecaseImpl) InvalidateListCaches() error {
	return u.cacheService.DeletePattern("products:list:*")
}
