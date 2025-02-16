package helpers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return 0, fmt.Errorf("user not authenticated")
	}
	return userID.(uint), nil
}
