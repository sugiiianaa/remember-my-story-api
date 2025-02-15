package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorHandlerMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// Log the error
			logger.WithFields(logrus.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"client": c.ClientIP(),
			}).Errorf("Request error: %v", err.Err)

			// Send standardized error response
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":  "internal server error",
				"status": http.StatusInternalServerError,
			})
		}
	}
}

func RecoveryMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				logger.WithFields(logrus.Fields{
					"path":   c.Request.URL.Path,
					"method": c.Request.Method,
					"client": c.ClientIP(),
					"panic":  err,
				}).Error("Recovered from panic")

				// Send standardized error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":  "internal server error",
					"status": http.StatusInternalServerError,
				})
			}
		}()

		c.Next()
	}
}
