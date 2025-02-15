package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.AbortWithStatusJSON(statusCode, gin.H{
		"error":  message,
		"status": statusCode,
	})
}

func ValidationErrorResponse(c *gin.Context, errors map[string]string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"error":   "validation failed",
		"details": errors,
		"status":  http.StatusBadRequest,
	})
}
