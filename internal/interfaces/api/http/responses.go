package http

import "go-microservice-product-porto/pkg/logger"

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewErrorResponse(err error) ErrorResponse {
	logger.Error().
		Err(err).
		Msg("Error occurred")

	return ErrorResponse{
		Error: err.Error(),
	}
}

func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	logger.Info().
		Str("message", message).
		Msg(message)

	return SuccessResponse{
		Message: message,
		Data:    data,
	}
}
