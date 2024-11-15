package usecase

import (
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/repositories"
	"go-microservice-product-porto/app/validation"
	"math"
)

type ProductUsecase struct {
	repo *repositories.ProductRepository
}

func NewProductUsecase(repo *repositories.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo}
}

func (b *ProductUsecase) GetAllProducts(page, pageSize int, sortBy, sortDir string, filters map[string]interface{}) (*models.PaginatedResponse, error) {
	products, err := b.repo.FindAllInMongo(page, pageSize, sortBy, sortDir, filters)
	if err != nil {
		return nil, err
	}

	totalItems, err := b.repo.CountDocuments(filters)
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

func (b *ProductUsecase) CreateProduct(product *models.Product) error {
	validator := validation.NewProductValidator(product)
	if err := validator.Validate(); err != nil {
		return err
	}

	err := b.repo.CreateInMongo(product)
	if err == nil {
		b.repo.GetCacheService().DeletePattern("products:list:*")
	}
	return err
}

func (b *ProductUsecase) GetProductByID(id string) (*models.Product, error) {
	return b.repo.FindByIDInMongo(id)
}

func (b *ProductUsecase) UpdateProduct(id string, product *models.Product) error {
	validator := validation.NewProductValidator(product)
	if err := validator.Validate(); err != nil {
		return err
	}

	err := b.repo.UpdateInMongo(id, product)
	if err == nil {
		b.repo.GetCacheService().Delete("product:" + id)
		b.repo.GetCacheService().DeletePattern("products:list:*")
	}
	return err
}

func (b *ProductUsecase) DeleteProduct(id string) error {
	err := b.repo.DeleteInMongo(id)
	if err == nil {
		b.repo.GetCacheService().Delete("product:" + id)
		b.repo.GetCacheService().DeletePattern("products:list:*")
	}
	return err
}
