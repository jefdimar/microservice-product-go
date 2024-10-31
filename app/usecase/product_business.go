package usecase

import (
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (b *ProductUsecase) GetProductByID(id primitive.ObjectID) (*models.Product, error) {
	mongoProduct, err := b.repo.FindByIDInMongo(id)
	if err != nil {
		return nil, err
	}
	product := &models.Product{
		ID:          mongoProduct.ID,
		Name:        mongoProduct.Name,
		Description: mongoProduct.Description,
		Price:       mongoProduct.Price,
		CreatedAt: 	 mongoProduct.CreatedAt,
		UpdatedAt:   mongoProduct.UpdatedAt,
	}
	return product, nil
}