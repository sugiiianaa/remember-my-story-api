package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	repositories "github.com/sugiiianaa/remember-my-story/internal/Repositories"
	"github.com/sugiiianaa/remember-my-story/internal/database"
	"github.com/sugiiianaa/remember-my-story/internal/handlers"
	"github.com/sugiiianaa/remember-my-story/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db, err := database.NewPostgresConnection(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	// Initialize layers
	journalRepo := repositories.NewJournalRepository(db)
	journalService := services.NewJournalService(journalRepo)
	journalHandler := handlers.NewJournalHandler(journalService)

	// Create Gin router
	router := gin.Default()

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

	// Start server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("failed to start server: " + err.Error())
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic("server forced to shutdown: " + err.Error())
	}
}
