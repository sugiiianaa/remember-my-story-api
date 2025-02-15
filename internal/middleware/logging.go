package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Collect information
		duration := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userAgent := c.Request.UserAgent()

		// Create log entry
		entry := logger.WithFields(logrus.Fields{
			"client_ip":  clientIP,
			"method":     method,
			"path":       path,
			"status":     status,
			"duration":   duration,
			"user_agent": userAgent,
		})

		// Log message based on status code
		if status >= 500 {
			entry.Error("server error")
		} else if status >= 400 {
			entry.Warn("client error")
		} else {
			entry.Info("request processed")
		}

		// Log request body for debugging (be cautious with sensitive data)
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			logger.Debugf("Request body: %s", getRequestBody(c))
		}
	}
}

func getRequestBody(c *gin.Context) string {
	body, err := c.GetRawData()
	if err != nil {
		return fmt.Sprintf("error reading body: %v", err)
	}

	// Restore the body back to the context
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return string(body)
}
