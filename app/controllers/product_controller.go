package controllers

import (
	"go-microservice-product-porto/app/handlers"
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/usecase"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	business *usecase.ProductUsecase
}

func NewProductController(usecase *usecase.ProductUsecase) *ProductController {
	return &ProductController{usecase}
}

// @Summary Create a new product
// @Description Create a new product with the provided input data
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Create product"
// @Success 201 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [post]
func (c *ProductController) Create(ctx *gin.Context) {
	var product models.Product
	if err := ctx.ShouldBindJSON(&product); err != nil {
		handlers.BadRequestResponse(ctx, err.Error())
		return
	}

	if err := c.business.CreateProduct(&product); err != nil {
		if strings.Contains(err.Error(), "validation") {
			handlers.ValidationErrorResponse(ctx, err.Error())
			return
		}
		handlers.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	handlers.SuccessResponse(ctx, http.StatusCreated, "Product created successfully", product)
}

// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]interface{}
// @Router /products [get]
func (c *ProductController) GetAll(ctx *gin.Context) {
	products, err := c.business.GetAllProducts()
	if err != nil {
		handlers.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	handlers.SuccessResponse(ctx, http.StatusOK, "Products retrieved successfully", products)
}

// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [get]
func (c *ProductController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	product, err := c.business.GetProductByID(id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			handlers.NotFoundResponse(ctx, "Product not found")
			return
		}
		handlers.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	handlers.SuccessResponse(ctx, http.StatusOK, "Product retrieved successfully", product)
}

// @Summary Update a product
// @Description Update a product with the provided input data
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.Product true "Update product"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [put]
func (c *ProductController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var product models.Product

	if err := ctx.ShouldBindJSON(&product); err != nil {
		handlers.BadRequestResponse(ctx, err.Error())
		return
	}

	_, err := c.business.GetProductByID(id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			handlers.NotFoundResponse(ctx, "Product not found")
			return
		}
		handlers.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	if err := c.business.UpdateProduct(id, &product); err != nil {
		if strings.Contains(err.Error(), "validation") {
			handlers.ValidationErrorResponse(ctx, err.Error())
			return
		}
		handlers.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	handlers.SuccessResponse(ctx, http.StatusOK, "Product updated successfully", product)
}

// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id} [delete]
func (c *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	err := c.business.DeleteProduct(id)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			handlers.NotFoundResponse(ctx, "Product not found")
			return
		}
		handlers.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	handlers.SuccessResponse(ctx, http.StatusOK, "Product deleted successfully", nil)
}
