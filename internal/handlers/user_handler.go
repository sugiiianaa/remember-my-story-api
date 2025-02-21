package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sugiiianaa/remember-my-story/internal/apperrors"
	"github.com/sugiiianaa/remember-my-story/internal/models"
	"github.com/sugiiianaa/remember-my-story/internal/services"
	"github.com/sugiiianaa/remember-my-story/pkg/helpers"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.InvalidRequestData,
			err.Error(),
		))
		return
	}

	userID, err := h.authService.Register(req.Email, req.FullName, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.UserAlreadyExist,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusCreated, helpers.SuccessResponse(map[string]interface{}{
		"user_id": userID,
	}))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.InvalidRequestData,
			err.Error(),
		))
		return
	}

	token, err := h.authService.Login(req.Email, req.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrorResponse(
			apperrors.InvalidCredentials,
			err.Error(),
		))
		return
	}

	c.JSON(http.StatusOK, helpers.SuccessResponse(map[string]interface{}{
		"token": token,
	}))
}
