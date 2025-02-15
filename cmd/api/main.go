package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	repositories "github.com/sugiiianaa/remember-my-story/internal/Repositories"
	"github.com/sugiiianaa/remember-my-story/internal/database"
	"github.com/sugiiianaa/remember-my-story/internal/handlers"
	"github.com/sugiiianaa/remember-my-story/internal/middleware"
	"github.com/sugiiianaa/remember-my-story/internal/services"
	"gorm.io/gorm"
)

func main() {
	env := configureEnvironment()
	logger := setupLogger(env)
	db := initDatabase(logger)
	router := setupRouter(logger, env, db)
	startServer(router, logger)
}

// --------------------------
// Configuration functions
// --------------------------

func configureEnvironment() string {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	env := strings.ToLower(strings.Trim(os.Getenv("SERVER_ENV"), "\" "))

	if env == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		env = "debug"
		gin.SetMode(gin.DebugMode)
	}

	return env
}

func setupLogger(env string) *logrus.Logger {
	logger := middleware.Logger(env)

	// Group all gin logging configurations
	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()

	logger.WithFields(logrus.Fields{
		"env":  env,
		"port": os.Getenv("APP_PORT"),
	}).Info("Starting server in ", strings.ToUpper(env), " mode")

	return logger
}

func initDatabase(logger *logrus.Logger) *gorm.DB {
	db, err := database.NewPostgresConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	return db
}

// --------------------------
// Router/Server functions
// --------------------------
func setupRouter(logger *logrus.Logger, env string, db *gorm.DB) *gin.Engine {
	// Initialize layers
	journalRepo := repositories.NewJournalRepository(db)
	journalService := services.NewJournalService(journalRepo)
	journalHandler := handlers.NewJournalHandler(journalService)

	router := gin.New()

	router.Use(
		middleware.RequestIDMiddleware(),
		middleware.LoggingMiddleware(logger, env),
	)

	registerRoutes(router, journalHandler)
	return router
}

func registerRoutes(router *gin.Engine, handler *handlers.JournalHandler) {
	api := router.Group("api/v1")
	{
		journals := api.Group("/journals")
		{
			journals.POST("", handler.CreateEntry)
			journals.GET("/:id", handler.GetEntry)
			journals.GET("", handler.GetEntriesByDate)
		}
	}
}

func startServer(router *gin.Engine, logger *logrus.Logger) {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown setup
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exited properly")
}
