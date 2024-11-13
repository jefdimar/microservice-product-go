package usecase

import (
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/repositories"
	"go-microservice-product-porto/app/validation"
)

type ProductUsecase struct {
	repo *repositories.ProductRepository
}

func NewProductUsecase(repo *repositories.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo}
}

func (b *ProductUsecase) CreateProduct(product *models.Product) error {
	validator := validation.NewProductValidator(product)
	if err := validator.Validate(); err != nil {
		return err
	}
	return b.repo.CreateInMongo(product)
}

func (b *ProductUsecase) GetAllProducts(page, pageSize int, sortBy, sortDir string) ([]models.Product, error) {
	return b.repo.FindAllInMongo(page, pageSize, sortBy, sortDir)
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
