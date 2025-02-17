package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = "unknown"
		}

		// Capture request body for logging (only in debug mode)
		var requestBody string
		if env == "debug" && c.Request.Body != nil {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			requestBody = string(bodyBytes)
			// Restore body for later use
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// Log request start
		logger.WithFields(logrus.Fields{
			"requestId": requestID,
			"method":    c.Request.Method,
			"path":      path,
			"ip":        c.ClientIP(),
			"body":      requestBody, // Only in debug mode
		}).Info("Request started")

		// Recover from panic
		defer func() {
			if err := recover(); err != nil {
				fields := logrus.Fields{
					"requestId": requestID,
					"method":    c.Request.Method,
					"path":      path,
					"ip":        c.ClientIP(),
				}

				if env == "debug" {
					fields["error"] = err
					logger.WithFields(fields).Error("Panic occurred")
				} else {
					logger.WithFields(fields).Error("Internal Server Error")
				}
				c.AbortWithStatus(500) // Return 500 status
			}
		}()

		// Process the request
		c.Next()

		// Log response
		latency := time.Since(start)
		status := c.Writer.Status()

		fields := logrus.Fields{
			"requestId": requestID,
			"status":    status,
			"method":    c.Request.Method,
			"path":      path,
			"ip":        c.ClientIP(),
			"latency":   latency,
		}

		// Log errors
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				if env == "debug" {
					fields["error"] = e.Error()
					logger.WithFields(fields).Error("Request failed")
				} else {
					logger.WithFields(fields).Error("An error occurred")
				}
			}
		} else {
			logger.WithFields(fields).Info("Request completed")
		}
	}
}
