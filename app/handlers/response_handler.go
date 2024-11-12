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
}

func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string, err interface{}) {
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
