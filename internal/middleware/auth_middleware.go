package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sugiiianaa/remember-my-story/internal/apperrors"
	"github.com/sugiiianaa/remember-my-story/pkg/helpers"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(
				apperrors.Unauthorized,
				"Valid authorization header required",
			))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(
				apperrors.Unauthorized,
				"Valid authorization header required",
			))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(
				apperrors.Unauthorized,
				"Invalid token",
			))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(
				apperrors.Unauthorized,
				"Invalid token",
			))
			return
		}

		sub, err := claims.GetSubject()

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(
				apperrors.Unauthorized,
				"Invalid token",
			))
			return
		}

		userID, err := strconv.ParseUint(sub, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrorResponse(
				apperrors.Unauthorized,
				"Invalid token",
			))
			return
		}

		c.Set("userID", uint(userID))
		c.Next()
	}

}
