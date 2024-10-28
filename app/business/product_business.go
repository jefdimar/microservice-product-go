package business

import (
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/repositories"
)

type ProductBusiness struct {
	repo *repositories.ProductRepository
}

func NewProductBusiness(repo *repositories.ProductRepository) *ProductBusiness {
	return &ProductBusiness{repo}
}

func (b *ProductBusiness) CreateProduct(product *models.Product) error {
	return b.repo.Create(product)
}

func (b *ProductBusiness) GetAllProducts() ([]models.Product, error) {
	return b.repo.FindAll()
}

func (b *ProductBusiness) GetProductByID(id uint) (*models.Product, error) {
	return b.repo.FindByID(id)
}
