package handlers

import "github.com/gin-gonic/gin"

type ResponseHandler interface {
	SuccessResponse(c *gin.Context, statusCode int, message string, data interface{})
	ErrorResponse(c *gin.Context, statusCode int, message string, err interface{})
	ValidationErrorResponse(c *gin.Context, err interface{})
	BadRequestResponse(c *gin.Context, err interface{})
	NotFoundResponse(c *gin.Context, message string)
	InternalServerErrorResponse(c *gin.Context, err interface{})
	ConflictResponse(c *gin.Context, err interface{})
	UnauthorizedResponse(c *gin.Context)
	ForbiddenResponse(c *gin.Context)
	NoContentResponse(c *gin.Context)
	PaginatedResponse(c *gin.Context, data interface{}, meta interface{})
}
