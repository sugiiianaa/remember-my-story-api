package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	repositories "github.com/sugiiianaa/remember-my-story/internal/Repositories"
	"github.com/sugiiianaa/remember-my-story/internal/models"
	"github.com/sugiiianaa/remember-my-story/internal/models/enums"
	"github.com/sugiiianaa/remember-my-story/internal/services"
)

type JournalHandler struct {
	service services.JournalService
}

func NewJournalHandler(service services.JournalService) *JournalHandler {
	return &JournalHandler{service: service}
}

func (h *JournalHandler) CreateEntry(c *gin.Context) {
	var entry models.JournalEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate mood
	if entry.Mood == enums.Mood.Unknown {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mood value"})
		return
	}

	// Get userID from context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if uid, ok := userID.(uint); ok {
		entry.UserID = uid
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	if err := h.service.CreateEntry(c.Request.Context(), &entry); err != nil {
		if _, ok := err.(services.ValidationError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create entry"})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *JournalHandler) GetEntry(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.service.GetEntry(c.Request.Context(), uint(id))
	if err != nil {
		if err == repositories.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "entry not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get entry"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *JournalHandler) GetEntriesByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
		return
	}

	entries, err := h.service.GetEntriesByDate(c.Request.Context(), date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get entries"})
		return
	}

	c.JSON(http.StatusOK, entries)
}
