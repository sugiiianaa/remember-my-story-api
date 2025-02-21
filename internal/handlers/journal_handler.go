package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sugiiianaa/remember-my-story/internal/apperrors"
	"github.com/sugiiianaa/remember-my-story/internal/models"
	"github.com/sugiiianaa/remember-my-story/internal/models/enums"
	"github.com/sugiiianaa/remember-my-story/internal/services"
	"github.com/sugiiianaa/remember-my-story/pkg/helpers"
	"gorm.io/gorm"
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

func (h *JournalHandler) UpdateEntry(c *gin.Context) {
	// Get journal entry ID from URL
	journalIDParam := c.Param("id")
	fmt.Printf("journalIDParam: %s\n", journalIDParam)
	if journalIDParam == "" {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.InvalidRequestData,
			"Missing journal ID",
		))
		return
	}

	// Convert ID to uint
	journalID, err := strconv.ParseUint(journalIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.InvalidRequestData,
			"Invalid journal ID",
		))
		return
	}

	// Get user ID from context
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, helpers.ErrorResponse(
			apperrors.Unauthorized,
			err.Error(),
		))
		return
	}

	// Define a struct for expected update fields with pointers to check presence
	var updateData struct {
		Mood               *string    `json:"mood,omitempty"`
		Date               *time.Time `json:"date,omitempty"`
		ThisDayDescription *string    `json:"this_day_description,omitempty"`
		DailyReflection    *string    `json:"daily_reflection,omitempty"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.InvalidRequestData,
			err.Error(),
		))
		return

	}
	// Convert the struct to a map for the service layer
	updateFields := make(map[string]interface{})

	// Handle Mood field separately to convert string to MoodType
	if updateData.Mood != nil {
		mood := enums.MoodFromString(*updateData.Mood)
		if mood == enums.Mood.Unknown {
			c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
				apperrors.InvalidRequestData,
				fmt.Sprintf("%s is an invalid mood.", *updateData.Mood),
			))
			return
		}
		updateFields["mood"] = mood
	}

	// Add other fields to the map if they are provided
	if updateData.Date != nil {
		updateFields["date"] = *updateData.Date
	}
	if updateData.ThisDayDescription != nil {
		updateFields["this_day_description"] = *updateData.ThisDayDescription
	}
	if updateData.DailyReflection != nil {
		updateFields["daily_reflection"] = *updateData.DailyReflection
	}

	// Call service layer to handle update logic
	err = h.service.UpdateEntry(uint(journalID), userID, updateFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, helpers.ErrorResponse(
				apperrors.NotFound,
				"Journal entry not found",
			))
			return
		}

		c.JSON(http.StatusInternalServerError, helpers.ErrorResponse(
			apperrors.InternalServerError,
			err.Error(),
		))
		return
	}

	// Success response
	c.JSON(http.StatusOK, helpers.SuccessResponse(map[string]interface{}{
		"message": "Journal entry updated successfully",
	}))
}
