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
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Configure environment
	env := strings.ToLower(strings.Trim(os.Getenv("SERVER_ENV"), "\" "))
	if env != "release" {
		env = "debug"
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize logger
	logger := middleware.Logger(env)

	// Completely disable Gin's default logging
	gin.DefaultWriter = io.Discard
	gin.DisableConsoleColor()

	// Log startup configuration
	logger.WithFields(logrus.Fields{
		"env":  env,
		"port": os.Getenv("APP_PORT"),
	}).Info("Starting server in ", strings.ToUpper(env), " mode")

	// Initialize database
	db, err := database.NewPostgresConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}

	// Initialize application layers
	journalRepo := repositories.NewJournalRepository(db)
	journalService := services.NewJournalService(journalRepo)
	journalHandler := handlers.NewJournalHandler(journalService)

	// Configure router
	router := gin.New()

	// Middlewares
	router.Use(
		middleware.RequestIDMiddleware(),
		middleware.LoggingMiddleware(logger, env),
	)

	// Routes
	api := router.Group("/api/v1")
	{
		journals := api.Group("/journals")
		{
			journals.POST("", journalHandler.CreateEntry)
			journals.GET("/:id", journalHandler.GetEntry)
			journals.GET("", journalHandler.GetEntriesByDate)
		}
	}

	// Configure server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown setup
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exited properly")
}
