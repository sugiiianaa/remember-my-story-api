package helpers

import "os"

type ApiResponse struct {
	Success bool        `json:"success"`           // Indicates if the request was successful
	Message string      `json:"message,omitempty"` // Optional message for additional context
	Data    interface{} `json:"data,omitempty"`    // The actual response data
	Error   *ApiError   `json:"error,omitempty"`   // Error details (if any)
}

type ApiError struct {
	Code      int    `json:"code"`                 // Error code (e.g., HTTP status code)
	Message   string `json:"message"`              // Error message for clients
	Details   string `json:"details,omitempty"`    // Additional error details (for debugging)
	RequestId string `json:"request_id,omitempty"` // Request id for production
}

// SuccessResponse returns a standardized success response
func SuccessResponse(message string, data interface{}) ApiResponse {
	return ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse returns a standardized error response
func ErrorResponse(code int, message string, details string) ApiResponse {

	serverEnvironment := os.Getenv("SERVER_ENV")

	apiError := ApiError{
		Code:    code,
		Message: message,
	}

	if serverEnvironment == "debug" {
		apiError.Details = details
	}

	return ApiResponse{
		Success: false,
		Error:   &apiError,
	}
}
