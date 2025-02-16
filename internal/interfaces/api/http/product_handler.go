package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"go-microservice-product-porto/internal/application/commands"
	"go-microservice-product-porto/internal/application/queries"
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
	var request struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	query := queries.GetProductQuery{ID: productID}
	product, err := h.queryHandler.HandleGetProduct(c.Request.Context(), query)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	products, err := h.queryHandler.HandleListProducts(c.Request.Context(), queries.ListProductsQuery{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	var cmd commands.UpdateStockCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	cmd.ProductID = productID

	if err := h.commandHandler.HandleUpdateStock(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stock updated successfully"})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product id is required"})
		return
	}

	cmd := commands.DeleteProductCommand{ProductID: productID}
	if err := h.commandHandler.HandleDeleteProduct(c.Request.Context(), cmd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

func (h *ProductHandler) SearchProducts(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": products,
		"filters": gin.H{
			"name":      query.Name,
			"min_price": query.MinPrice,
			"max_price": query.MaxPrice,
		},
	})
}
