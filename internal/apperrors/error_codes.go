package apperrors

import "net/http"

type ErrorCode struct {
	Code    string
	Message string
	Status  int
}

var (
	// ======================
	// Generic Errors
	// ======================
	InvalidRequestData = ErrorCode{
		Code:    "invalid_request_data",
		Message: "There is something wrong with your request payload",
		Status:  http.StatusBadRequest,
	}

	Unauthorized = ErrorCode{
		Code:    "unauthorized_request",
		Message: "Your request is unauthorized, please login to continue",
		Status:  http.StatusUnauthorized,
	}

	Forbidden = ErrorCode{
		Code:    "forbidden_request",
		Message: "You don't have permission to access this resource",
		Status:  http.StatusForbidden,
	}

	NotFound = ErrorCode{
		Code:    "resource_not_found",
		Message: "The requested resource was not found",
		Status:  http.StatusNotFound,
	}

	InternalServerError = ErrorCode{
		Code:    "internal_server_error",
		Message: "Something went wrong on our end. Please try again later.",
		Status:  http.StatusInternalServerError,
	}

	ServiceUnavailable = ErrorCode{
		Code:    "service_unavailable",
		Message: "The service is temporarily unavailable. Please try again later.",
		Status:  http.StatusServiceUnavailable,
	}

	TimeoutError = ErrorCode{
		Code:    "request_timeout",
		Message: "The request timed out. Please try again.",
		Status:  http.StatusRequestTimeout,
	}

	Conflict = ErrorCode{
		Code:    "conflict_error",
		Message: "A conflict occurred while processing the request.",
		Status:  http.StatusConflict,
	}

	TooManyRequests = ErrorCode{
		Code:    "too_many_requests",
		Message: "Too many requests. Please try again later.",
		Status:  http.StatusTooManyRequests,
	}

	// ======================
	// Authentication Errors
	// ======================
	InvalidCredentials = ErrorCode{
		Code:    "invalid_credentials",
		Message: "The provided email or password is invalid",
		Status:  http.StatusUnauthorized,
	}

	TokenExpired = ErrorCode{
		Code:    "token_expired",
		Message: "Your session has expired. Please log in again.",
		Status:  http.StatusUnauthorized,
	}

	InvalidToken = ErrorCode{
		Code:    "invalid_token",
		Message: "The provided token is invalid",
		Status:  http.StatusUnauthorized,
	}

	// ======================
	// User Related Errors
	// ======================
	UserAlreadyExist = ErrorCode{
		Code:    "user_already_exist",
		Message: "User already exists",
		Status:  http.StatusBadRequest,
	}

	UserNotFound = ErrorCode{
		Code:    "user_not_found",
		Message: "User not found",
		Status:  http.StatusNotFound,
	}

	// ======================
	// Database Errors
	// ======================
	DatabaseConnectionError = ErrorCode{
		Code:    "database_connection_error",
		Message: "Unable to connect to the database",
		Status:  http.StatusInternalServerError,
	}

	RecordNotFound = ErrorCode{
		Code:    "record_not_found",
		Message: "The requested record was not found in the database",
		Status:  http.StatusNotFound,
	}
)
