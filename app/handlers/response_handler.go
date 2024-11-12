package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Links   interface{} `json:"links,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func PaginatedResponse(c *gin.Context, data interface{}, meta interface{}) {
	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Data retrieved successfully",
		Data:    data,
		Errors:  meta,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
	c.JSON(statusCode, Response{
		Status:  "error",
		Message: message,
		Errors:  err,
	})
}

func ValidationErrorResponse(c *gin.Context, err interface{}) {
	c.JSON(http.StatusUnprocessableEntity, Response{
		Status:  "error",
		Message: "Validation failed",
		Errors:  err,
	})
}

func BadRequestResponse(c *gin.Context, err interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Status:  "error",
		Message: "Invalid Request",
		Errors:  err,
	})
}

func InternalServerErrorResponse(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, Response{
		Status:  "error",
		Message: "Internal Server Error",
		Errors:  err,
	})
}

func NotFoundResponse(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Status:  "error",
		Message: message,
	})
}

func ConflictResponse(c *gin.Context, err interface{}) {
	c.JSON(http.StatusConflict, Response{
		Status:  "error",
		Message: "Resource conflict",
		Errors:  err,
	})
}

func UnauthorizedResponse(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		Status:  "error",
		Message: "Unauthorized",
	})
}

func ForbiddenResponse(c *gin.Context) {
	c.JSON(http.StatusForbidden, Response{
		Status:  "error",
		Message: "Forbidden",
	})
}

func NoContentResponse(c *gin.Context) {
	c.JSON(http.StatusNoContent, nil)
}
