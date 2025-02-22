package helpers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sugiiianaa/remember-my-story/internal/apperrors"
)

type ApiResponse struct {
	Success bool        `json:"success"`         // Indicates if the request was successful
	Data    interface{} `json:"data,omitempty"`  // The actual response data
	Error   *ApiError   `json:"error,omitempty"` // Error details (if any)
}

type ApiError struct {
	ErrorCode string `json:"error_code,omitempty"` // Error code to determine which error are happen
	Message   string `json:"message"`              // Error message for clients
	Details   string `json:"details,omitempty"`    // Additional error details (for debugging)
	RequestId string `json:"request_id,omitempty"` // Request id for production
}

// SuccessResponse returns a standardized success response
func SuccessResponse(data interface{}) ApiResponse {
	return ApiResponse{
		Success: true,
		Data:    data,
	}
}

// ErrorResponse returns a standardized error response
func ErrorResponse(errCode apperrors.ErrorCode, details string) ApiResponse {

	serverEnvironment := os.Getenv("SERVER_ENV")

	apiError := ApiError{
		ErrorCode: errCode.Code,
		Message:   errCode.Message,
	}

	if serverEnvironment == "debug" {
		apiError.Details = details
	}

	return ApiResponse{
		Success: false,
		Error:   &apiError,
	}
}

func RespondWithSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, SuccessResponse(data))
}

func RespondWithError(c *gin.Context, statusCode int, errCode apperrors.ErrorCode, details string) {
	c.JSON(statusCode, ErrorResponse(errCode, details))
}
