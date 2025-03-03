package main

import (
	"context"
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
	logger := initLogger(env)

	// Log the environment and log level
	logger.WithFields(logrus.Fields{
		"environment": env,
		"logLevel":    logger.GetLevel().String(),
	}).Info("Starting application with the following settings")

	db := initDatabase(logger)
	router := setupRouter(logger, env, db)
	startServer(router, logger)
}

// --------------------------
// Configuration functions
// --------------------------

func configureEnvironment() string {
	if err := godotenv.Load(); err != nil {
		_, err := os.Stat(".env")
		if os.IsNotExist(err) {
			log.Println(".env file does not exist!")
		} else if err != nil {
			log.Println("Error accessing .env:", err)
		} else {
			log.Println(".env file found")
		}
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

func initLogger(env string) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set log level based on environment
	if env == "debug" {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

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
	// jwt setup
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Fatal("JWT_SECRET environment variable not set")
	}
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	// journal setup
	journalRepo := repositories.NewJournalRepository(db)
	journalService := services.NewJournalService(journalRepo)
	journalHandler := handlers.NewJournalHandler(journalService)

	// auth setup
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(*userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(*authService)

	router := gin.New()

	router.Use(
		middleware.LoggingMiddleware(logger, env),
	)

	registerRoutes(router, journalHandler, authHandler, authMiddleware)
	return router
}

func registerRoutes(
	router *gin.Engine,
	handler *handlers.JournalHandler,
	authHandler *handlers.AuthHandler,
	authMiddleware gin.HandlerFunc) {
	api := router.Group("api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		journals := api.Group("/journals")
		journals.Use(authMiddleware)
		{
			journals.POST("", handler.CreateEntry)
			// journals.GET("/:id", handler.GetEntry)
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
