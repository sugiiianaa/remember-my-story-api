package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sugiiianaa/remember-my-story/internal/apperrors"
	"github.com/sugiiianaa/remember-my-story/internal/models"
	"github.com/sugiiianaa/remember-my-story/internal/models/enums"
	"github.com/sugiiianaa/remember-my-story/internal/services"
	"github.com/sugiiianaa/remember-my-story/pkg/helpers"
)

type JournalHandler struct {
	service *services.JournalService
}

func NewJournalHandler(service *services.JournalService) *JournalHandler {
	return &JournalHandler{service: service}
}

func (h *JournalHandler) CreateEntry(c *gin.Context) {
	var entry models.JournalEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.InvalidRequestData,
			err.Error(),
		))
		return
	}

	// Validate mood
	if entry.Mood == enums.Mood.Unknown {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.InvalidRequestData,
			fmt.Sprintf("%s is invalid mood.", entry.Mood),
		))
	}

	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.ErrorResponse(
			apperrors.Unauthorized,
			err.Error(),
		))
	}

	entry.UserID = userID

	journalID, err := h.service.CreateEntry(&entry)

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(
			apperrors.InternalServerError,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, helpers.SuccessResponse(map[string]interface{}{
		"journal_id": journalID,
	}))
}

func (h *JournalHandler) GetEntry(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	entry, err := h.service.GetEntry(c.Request.Context(), uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(
			apperrors.InternalServerError,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, helpers.SuccessResponse(entry))
}
