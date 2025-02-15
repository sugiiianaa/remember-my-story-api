package middleware

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct {
	logrus.TextFormatter
	env string
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Handle non-request logs differently
	if _, ok := entry.Data["method"]; !ok {
		return f.formatApplicationLog(entry)
	}

	if f.env == "debug" {
		return f.formatDevelopment(entry)
	}
	return f.formatProduction(entry)
}

func (f *CustomFormatter) formatApplicationLog(entry *logrus.Entry) ([]byte, error) {
	if f.env == "debug" {
		return f.formatDevelopment(entry)
	}

	return (&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}).Format(entry)
}

func (f *CustomFormatter) formatDevelopment(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = 36 // Cyan
	case logrus.InfoLevel:
		levelColor = 32 // Green
	case logrus.WarnLevel:
		levelColor = 33 // Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // Red
	default:
		levelColor = 37 // White
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("\x1b[%dm%-7s\x1b[0m [%s] %s",
		levelColor,
		entry.Level.String(),
		timestamp,
		entry.Message,
	)

	for k, v := range entry.Data {
		msg += fmt.Sprintf(" \x1b[90m%s\x1b[0m=%v", k, v)
	}
	msg += "\n"

	return []byte(msg), nil
}

func (f *CustomFormatter) formatProduction(entry *logrus.Entry) ([]byte, error) {
	// Filter and sanitize production logs
	safeEntry := logrus.Entry{
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
		Data: logrus.Fields{
			"client_ip":  maskIP(entry.Data["client_ip"]),
			"request_id": entry.Data["request_id"],
			"method":     entry.Data["method"],
			"path":       sanitizePath(entry.Data["path"]),
			"status":     entry.Data["status"],
		},
	}

	return (&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			"msg": "message",
		},
	}).Format(&safeEntry)
}

func maskIP(ip interface{}) string {
	if ipStr, ok := ip.(string); ok {
		if ipStr == "::1" {
			return "localhost"
		}
		// Add IP masking logic for production
		// Example: "192.168.1.1" -> "192.168.x.x"
		parts := strings.Split(ipStr, ".")
		if len(parts) == 4 {
			return fmt.Sprintf("%s.%s.x.x", parts[0], parts[1])
		}
		return ipStr
	}
	return "unknown"
}

func sanitizePath(path interface{}) string {
	if pathStr, ok := path.(string); ok {
		// Sanitize UUIDs in paths
		return uuidRegEx.ReplaceAllString(pathStr, ":id")
	}
	return "unknown"
}

var uuidRegEx = regexp.MustCompile(`[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`)

func Logger(env string) *logrus.Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)

	log.SetFormatter(&CustomFormatter{
		TextFormatter: logrus.TextFormatter{
			DisableTimestamp: true,
			ForceColors:      env == "debug",
		},
		env: env,
	})

	if env == "debug" {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	return log
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

func LoggingMiddleware(log *logrus.Logger, env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetString("request_id")

		// Log request body only in debug mode
		var body string
		if env == "debug" {
			body = getRequestBody(c)
		}

		// Log request start
		if env == "debug" {
			log.WithFields(logrus.Fields{
				"method":     c.Request.Method,
				"path":       c.Request.URL.Path,
				"request_id": requestID,
			}).Debug("Request started")
		}

		c.Next()

		// Collect logging data
		latency := time.Since(start)
		status := c.Writer.Status()

		entry := log.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     status,
			"latency":    fmt.Sprintf("%.3fms", latency.Seconds()*1000),
			"client_ip":  c.ClientIP(),
			"request_id": requestID,
		})

		if env == "debug" && body != "" {
			entry = entry.WithField("body", body)
		}

		// Security: Never log 4xx bodies in production
		switch {
		case status >= 500:
			entry.Error("Server error")
		case status >= 400:
			if env == "debug" {
				entry.Warn("Client error")
			} else {
				entry.Warn("Client error")
				// Log additional security info in production
				log.WithFields(logrus.Fields{
					"request_id": requestID,
					"client_ip":  maskIP(c.ClientIP()),
					"user_agent": c.Request.UserAgent(),
				}).Warn("Potential client issue")
			}
		default:
			entry.Info("Request processed")
		}
	}
}

func getRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Sprintf("error reading body: %v", err)
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return string(bodyBytes)
}
