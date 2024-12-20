package controllers

import (
	"go-microservice-product-porto/app/handlers"
	"go-microservice-product-porto/app/helpers"
	"go-microservice-product-porto/app/models"
	"go-microservice-product-porto/app/usecase"
	"go-microservice-product-porto/app/validation"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type ProductController struct {
	business    usecase.ProductUsecase
	respHandler handlers.ResponseHandler
}

func NewProductController(usecase usecase.ProductUsecase, handler handlers.ResponseHandler) *ProductController {
	return &ProductController{
		business:    usecase,
		respHandler: handler,
	}
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
		c.respHandler.BadRequestResponse(ctx, err.Error())
		return
	}

	if err := c.business.CreateProduct(&product); err != nil {
		if strings.Contains(err.Error(), "validation") {
			c.respHandler.ValidationErrorResponse(ctx, err.Error())
			return
		}
		c.respHandler.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	c.business.InvalidateListCaches()

	c.respHandler.SuccessResponse(ctx, http.StatusCreated, "Product created successfully", product)
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
		SortDir:  strings.ToLower(ctx.Query("sortDir")), // Normalize the direction
	}

	// Validate sort direction
	if params.SortDir != "" && params.SortDir != "asc" && params.SortDir != "desc" {
		params.SortDir = "desc" // Set default if invalid
	}

	filters := make(map[string]interface{})

	// Search filter
	if search := ctx.Query("search"); search != "" {
		filters["search"] = search
	}

	// Name filter
	if name := ctx.Query("name"); name != "" {
		filters["name"] = name
	}

	// Price range filter
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

	// Date range filter with format "02-Jan-2006"
	if startDate := ctx.Query("start_date"); startDate != "" {
		if date, err := time.Parse("02-Jan-2006", startDate); err == nil {
			// Set to start of day
			date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
			filters["start_date"] = date
		}
	}

	if endDate := ctx.Query("end_date"); endDate != "" {
		if date, err := time.Parse("02-Jan-2006", endDate); err == nil {
			// Set to end of day
			date = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 999999999, time.UTC)
			filters["end_date"] = date
		}
	}

	// Stock filter
	if minStock := ctx.Query("stock_min"); minStock != "" {
		if stock, err := strconv.Atoi(minStock); err == nil {
			filters["stock_min"] = stock
		}
	}
	if maxStock := ctx.Query("stock_max"); maxStock != "" {
		if stock, err := strconv.Atoi(maxStock); err == nil {
			filters["stock_max"] = stock
		}
	}

	// Active status filter
	if isActive := ctx.Query("is_active"); isActive != "" {
		if active, err := strconv.ParseBool(isActive); err == nil {
			filters["is_active"] = active
		}
	}

	// SKU filter
	if sku := ctx.Query("sku"); sku != "" {
		filters["sku"] = sku
	}

	if err := validation.ValidateQueryParams(params); err != nil {
		c.respHandler.ValidationErrorResponse(ctx, err.Error())
		return
	}

	paginatedResponse, err := c.business.GetAllProducts(params.Page, params.PageSize, params.SortBy, params.SortDir, filters)
	if err != nil {
		if err == redis.Nil {
			// Cache miss - fetch from database
			paginatedResponse, err = c.business.GetAllProducts(params.Page, params.PageSize, params.SortBy, params.SortDir, filters)
			if err != nil {
				c.respHandler.InternalServerErrorResponse(ctx, err.Error())
				return
			}
		} else {
			c.respHandler.InternalServerErrorResponse(ctx, err.Error())
			return
		}
	}

	c.respHandler.SuccessResponse(ctx, http.StatusOK, "Products retrieved successfully", paginatedResponse)
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
		c.respHandler.BadRequestResponse(ctx, err.Error())
		return
	}

	product, err := c.business.GetProductByID(id)
	if err != nil {
		if err.Error() == "product not found" {
			c.respHandler.NotFoundResponse(ctx, "Product not found")
			return
		}
		// Handle cache miss specifically
		if err == redis.Nil {
			// Cache miss - continue with database lookup
			product, err = c.business.GetProductByID(id)
			if err != nil {
				if err.Error() == "mongo: no documents in result" {
					c.respHandler.NotFoundResponse(ctx, "Product not found")
					return
				}
				c.respHandler.InternalServerErrorResponse(ctx, err.Error())
				return
			}
		} else if err.Error() == "mongo: no documents in result" {
			c.respHandler.NotFoundResponse(ctx, "Product not found")
			return
		} else {
			c.respHandler.InternalServerErrorResponse(ctx, err.Error())
			return
		}

	}

	c.respHandler.SuccessResponse(ctx, http.StatusOK, "Product retrieved successfully", product)
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
		c.respHandler.BadRequestResponse(ctx, err.Error())
		return
	}

	var updates models.ProductUpdate
	if err := ctx.ShouldBindJSON(&updates); err != nil {
		c.respHandler.BadRequestResponse(ctx, err.Error())
		return
	}

	if err := c.business.UpdateProduct(id, &updates); err != nil {

		if err.Error() == "no fields to update" {
			c.respHandler.BadRequestResponse(ctx, err.Error())
			return
		}
		c.respHandler.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	// Fetch the updated product
	updatedProduct, _ := c.business.GetProductByID(id)

	c.respHandler.SuccessResponse(ctx, http.StatusOK, "Product updated successfully", updatedProduct)
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
		c.respHandler.BadRequestResponse(ctx, err.Error())
		return
	}

	err := c.business.DeleteProduct(id)
	if err != nil {
		if err.Error() == "product not found" {
			c.respHandler.NotFoundResponse(ctx, "Product not found")
			return
		}

		if err.Error() == "mongo: no documents in result" {
			c.respHandler.NotFoundResponse(ctx, "Product not found")
			return
		}
		c.respHandler.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	// After successful deletion, invalidate related caches
	if err := c.business.InvalidateRelatedCaches(id); err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	c.respHandler.SuccessResponse(ctx, http.StatusOK, "Product deleted successfully", nil)
}

// @Summary Update product stock
// @Description Update product stock with movement tracking
// @Tags stock
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body StockUpdateRequest true "Stock update details"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id}/stock [put]
func (c *ProductController) UpdateStock(ctx *gin.Context) {
	id := ctx.Param("id")

	var request struct {
		Stock  int    `json:"stock" binding:"required"`
		Reason string `json:"reason" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		c.respHandler.BadRequestResponse(ctx, err.Error())
		return
	}

	err := c.business.UpdateProductStock(id, request.Stock, request.Reason)
	if err != nil {
		c.respHandler.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	c.respHandler.SuccessResponse(ctx, http.StatusOK, "Stock updated successfully", nil)
}

// @Summary Get stock movement history
// @Description Get stock movement history for a product
// @Tags stock
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {array} models.StockMovement
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /products/{id}/stock-movements [get]
func (c *ProductController) GetStockMovements(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := validation.ValidateObjectID(id); err != nil {
		c.respHandler.BadRequestResponse(ctx, err.Error())
		return
	}

	movements, err := c.business.GetStockMovement(id)
	if err != nil {
		c.respHandler.InternalServerErrorResponse(ctx, err.Error())
		return
	}

	c.respHandler.SuccessResponse(ctx, http.StatusOK, "Stock movements retrieved successfully", movements)
}
