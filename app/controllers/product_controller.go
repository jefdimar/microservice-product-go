package controllers

import (
	"go-microservice-product-porto/app/handlers"
	"go-microservice-product-porto/app/helpers"
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/usecase"
	"go-microservice-product-porto/app/validation"
	"net/http"
	"strconv"
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
// @Failure 422 {object} map[string]interface{}
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
// @Description Get a list of all products with pagination and sorting
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param pageSize query int false "Items per page (default: 10)"
// @Param sortBy query string false "Sort field (name, price, created_at, updated_at)"
// @Param sortDir query string false "Sort direction (asc, desc)"
// @Success 200 {array} models.Product
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products [get]
func (c *ProductController) GetAll(ctx *gin.Context) {
	params := &validation.QueryParams{
		Page:     helpers.ParseInt(ctx.Query("page"), 1),
		PageSize: helpers.ParseInt(ctx.Query("pageSize"), 10),
		SortBy:   ctx.Query("sortBy"),
		SortDir:  ctx.Query("sortDir"),
	}

	filters := make(map[string]interface{})

	if search := ctx.Query("search"); search != "" {
		filters["search"] = search
	}

	if name := ctx.Query("name"); name != "" {
		filters["name"] = name
	}

	if minPrice := ctx.Query("price_min"); minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			filters["price_min"] = price
		}
	}

	if maxPrice := ctx.Query("price_max"); maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			filters["price_max"] = price
		}
	}

	if isActive := ctx.Query("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			filters["is_active"] = active
		}
	}

	if err := validation.ValidateQueryParams(params); err != nil {
		handlers.ValidationErrorResponse(ctx, err.Error())
		return
	}
	paginatedResponse, err := c.business.GetAllProducts(params.Page, params.PageSize, params.SortBy, params.SortDir, filters)
	if err != nil {
		handlers.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	handlers.SuccessResponse(ctx, http.StatusOK, "Products retrieved successfully", paginatedResponse)
}

// @Summary Get product by ID
// @Description Get a product by its ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID (MongoDB ObjectID)"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [get]
func (c *ProductController) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := validation.ValidateObjectID(id); err != nil {
		handlers.BadRequestResponse(ctx, err.Error())
		return
	}

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
// @Param id path string true "Product ID (MongoDB ObjectID)"
// @Param product body models.Product true "Update product"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [put]
func (c *ProductController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := validation.ValidateObjectID(id); err != nil {
		handlers.BadRequestResponse(ctx, err.Error())
		return
	}
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
// @Param id path string true "Product ID (MongoDB ObjectID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /products/{id} [delete]
func (c *ProductController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := validation.ValidateObjectID(id); err != nil {
		handlers.BadRequestResponse(ctx, err.Error())
		return
	}

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
