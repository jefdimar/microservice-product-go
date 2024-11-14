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

type PaginatedResponse struct {
	Data       []models.Product `json:"data"`
	Pagination PaginationMeta   `json:"pagination"`
}

type PaginationMeta struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalItems  int64 `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
}

func (b *ProductUsecase) GetAllProducts(page, pageSize int, sortBy, sortDir string, filters map[string]interface{}) (*PaginatedResponse, error) {
	products, err := b.repo.FindAllInMongo(page, pageSize, sortBy, sortDir, filters)
	if err != nil {
		return nil, err
	}

	totalItems, err := b.repo.CountDocuments(filters)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalItems) / float64(pageSize)))

	return &PaginatedResponse{
		Data: products,
		Pagination: PaginationMeta{
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
	return b.repo.CreateInMongo(product)
}

func (b *ProductUsecase) GetProductByID(id string) (*models.Product, error) {
	return b.repo.FindByIDInMongo(id)
}

func (b *ProductUsecase) UpdateProduct(id string, product *models.Product) error {
	validator := validation.NewProductValidator(product)
	if err := validator.Validate(); err != nil {
		return err
	}
	return b.repo.UpdateInMongo(id, product)
}

func (b *ProductUsecase) DeleteProduct(id string) error {
	return b.repo.DeleteInMongo(id)
}
