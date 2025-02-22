package handlers

import (
	"errors"
	"net/http"
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
		helpers.RespondWithError(c, http.StatusBadRequest, apperrors.InvalidRequestData, err.Error())
		return
	}

	if err := helpers.ValidateMood(entry.Mood.String()); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, apperrors.InvalidRequestData, err.Error())
		return
	}

	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, apperrors.Unauthorized, err.Error())
		return
	}

	entry.UserID = userID

	journalID, err := h.service.CreateEntry(&entry)
	if err != nil {
		helpers.HandleInternalError(c, err)
		return
	}

	helpers.RespondWithSuccess(c, http.StatusCreated, map[string]interface{}{
		"journal_id": journalID,
	})
}

func (h *JournalHandler) UpdateEntry(c *gin.Context) {
	journalID, err := helpers.GetUintParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, apperrors.InvalidRequestData, err.Error())
		return
	}

	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, apperrors.Unauthorized, err.Error())
		return
	}

	var updateData struct {
		Mood               *string    `json:"mood,omitempty"`
		Date               *time.Time `json:"date,omitempty"`
		ThisDayDescription *string    `json:"this_day_description,omitempty"`
		DailyReflection    *string    `json:"daily_reflection,omitempty"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, apperrors.InvalidRequestData, err.Error())
		return
	}

	updateFields := make(map[string]interface{})

	if updateData.Mood != nil {
		if err := helpers.ValidateMood(*updateData.Mood); err != nil {
			helpers.RespondWithError(c, http.StatusBadRequest, apperrors.InvalidRequestData, err.Error())
			return
		}
		updateFields["mood"] = enums.MoodFromString(*updateData.Mood)
	}

	if updateData.Date != nil {
		updateFields["date"] = *updateData.Date
	}
	if updateData.ThisDayDescription != nil {
		updateFields["this_day_description"] = *updateData.ThisDayDescription
	}
	if updateData.DailyReflection != nil {
		updateFields["daily_reflection"] = *updateData.DailyReflection
	}

	err = h.service.UpdateEntry(journalID, userID, updateFields)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			helpers.RespondWithError(c, http.StatusNotFound, apperrors.NotFound, "Journal entry not found")
			return
		}
		helpers.HandleInternalError(c, err)
		return
	}

	helpers.RespondWithSuccess(c, http.StatusOK, map[string]interface{}{
		"message": "Journal entry updated successfully",
	})
}

func (h *JournalHandler) GetAllEntry(c *gin.Context) {
	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, apperrors.Unauthorized, err.Error())
		return
	}

	entries, err := h.service.GetAllEntry(userID)
	if err != nil {
		helpers.HandleInternalError(c, err)
		return
	}

	helpers.RespondWithSuccess(c, http.StatusOK, entries)
}

func (h *JournalHandler) DeleteEntry(c *gin.Context) {
	journalID, err := helpers.GetUintParam(c, "id")
	if err != nil {
		helpers.RespondWithError(c, http.StatusBadRequest, apperrors.InvalidRequestData, err.Error())
		return
	}

	userID, err := helpers.GetUserIDFromContext(c)
	if err != nil {
		helpers.RespondWithError(c, http.StatusUnauthorized, apperrors.Unauthorized, err.Error())
		return
	}

	err = h.service.DeleteEntry(journalID, userID)
	if err != nil {
		helpers.HandleError(c, err, apperrors.InvalidRequestData)
		return
	}

	helpers.RespondWithSuccess(c, http.StatusOK, map[string]interface{}{
		"journal_id": journalID,
	})
}
