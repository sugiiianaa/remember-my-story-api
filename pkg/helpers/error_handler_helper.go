package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sugiiianaa/remember-my-story/internal/apperrors"
)

func HandleError(c *gin.Context, err error, errCode apperrors.ErrorCode) {
	c.JSON(http.StatusBadRequest, ErrorResponse(errCode, err.Error()))
}

func HandleInternalError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, ErrorResponse(apperrors.InternalServerError, err.Error()))
}

func HandleUnauthorized(c *gin.Context, err error) {
	c.JSON(http.StatusUnauthorized, ErrorResponse(apperrors.Unauthorized, err.Error()))
}
