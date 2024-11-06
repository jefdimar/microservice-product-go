package usecase

import (
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/repositories"
)

type ProductUsecase struct {
	repo *repositories.ProductRepository
}

func NewProductUsecase(repo *repositories.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo}
}

func (b *ProductUsecase) CreateProduct(product *models.Product) error {
	return b.repo.CreateInMongo(product)
}

func (b *ProductUsecase) GetAllProducts() ([]models.Product, error) {
	return b.repo.FindAllInMongo()
}

func (b *ProductUsecase) GetProductByID(id string) (*models.Product, error) {
	return b.repo.FindByIDInMongo(id)
}