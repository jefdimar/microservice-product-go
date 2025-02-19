package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"go-microservice-product-porto/internal/application/commands"
	"go-microservice-product-porto/internal/application/queries"
	"go-microservice-product-porto/pkg/common"
	"go-microservice-product-porto/pkg/logger"
)

type ProductHandler struct {
	commandHandler *commands.ProductCommandHandler
	queryHandler   *queries.ProductQueryHandler
}

func NewProductHandler(commandHandler *commands.ProductCommandHandler, queryHandler *queries.ProductQueryHandler) *ProductHandler {
	return &ProductHandler{
		commandHandler: commandHandler,
		queryHandler:   queryHandler,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	logger.Info().
		Str("handler", "CreateProduct").
		Msg("Creating a new product")

	var request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error().
			Str("handler", "CreateProduct").
			Err(err).
			Msg("Error binding JSON")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := commands.CreateProductCommand{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
	}

	if err := h.commandHandler.HandleCreateProduct(c.Request.Context(), cmd); err != nil {
		logger.Error().
			Str("handler", "CreateProduct").
			Err(err).
			Msg("Error handling create product command")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info().
		Str("handler", "CreateProduct").
		Msg("Product created successfully")

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	logger.Info().
		Str("handler", "GetProduct").
		Msg("Fetching product details")

	productID := c.Param("id")
	if productID == "" {
		logger.Error().
			Str("handler", "GetProduct").
			Msg("Product ID is required")

		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	query := queries.GetProductQuery{ID: productID}
	product, err := h.queryHandler.HandleGetProduct(c.Request.Context(), query)
	if err != nil {
		logger.Error().
			Str("handler", "GetProduct").
			Err(err).
			Msg("Error fetching product details")

		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	logger.Info().
		Str("handler", "GetProduct").
		Msg("Product details fetched successfully")

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	logger.Info().
		Str("handler", "ListProducts").
		Msg("Fetching list of products")

	query := queries.ListProductsQuery{
		Page:     common.ParseInt(c.DefaultQuery("page", "1")),
		PageSize: common.ParseInt(c.DefaultQuery("page_size", "10")),
		SortBy:   c.DefaultQuery("sort_by", ""),
		SortDir:  c.DefaultQuery("sort_dir", "asc"),
	}

	result, err := h.queryHandler.HandleListProducts(c.Request.Context(), query)
	if err != nil {
		logger.Error().
			Str("handler", "ListProducts").
			Err(err).
			Msg("Error fetching list of products")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info().
		Str("handler", "ListProducts").
		Msg("List of products fetched successfully")

	c.JSON(http.StatusOK, result)
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	logger.Info().
		Str("handler", "UpdateStock").
		Msg("Updating product stock")

	productID := c.Param("id")
	if productID == "" {
		logger.Error().
			Str("handler", "UpdateStock").
			Msg("Product ID is required")

		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	var cmd commands.UpdateStockCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		logger.Error().
			Str("handler", "UpdateStock").
			Err(err).
			Msg("Error binding JSON")

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProductID = productID

	if err := h.commandHandler.HandleUpdateStock(c.Request.Context(), cmd); err != nil {
		logger.Error().
			Str("handler", "UpdateStock").
			Err(err).
			Msg("Error updating stock")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info().
		Str("handler", "UpdateStock").
		Msg("Stock updated successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	logger.Info().
		Str("handler", "DeleteProduct").
		Msg("Deleting product")

	productID := c.Param("id")
	if productID == "" {
		logger.Error().
			Str("handler", "DeleteProduct").
			Msg("Product ID is required")

		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	cmd := commands.DeleteProductCommand{ProductID: productID}
	if err := h.commandHandler.HandleDeleteProduct(c.Request.Context(), cmd); err != nil {
		logger.Error().
			Str("handler", "DeleteProduct").
			Err(err).
			Msg("Error deleting product")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info().
		Str("handler", "DeleteProduct").
		Msg("Product deleted successfully")

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
	logger.Info().
		Str("handler", "SearchProducts").
		Msg("Searching for products")

	query := queries.SearchProductsQuery{
		Name: strings.TrimSpace(c.Query("name")),
	}

	if minPriceStr := c.Query("min_price"); minPriceStr != "" {
		minPrice, _ := strconv.ParseFloat(minPriceStr, 64)
		query.MinPrice = minPrice
	}

	if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
		maxPrice, _ := strconv.ParseFloat(maxPriceStr, 64)
		query.MaxPrice = maxPrice
	}

	products, err := h.queryHandler.HandleSearchProducts(c.Request.Context(), query)
	if err != nil {
		logger.Error().
			Str("handler", "SearchProducts").
			Err(err).
			Msg("Error searching for products")

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Info().
		Str("handler", "SearchProducts").
		Msg("Products searched successfully")

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"filters": gin.H{
			"name":      query.Name,
			"min_price": query.MinPrice,
			"max_price": query.MaxPrice,
		},
	})
}
