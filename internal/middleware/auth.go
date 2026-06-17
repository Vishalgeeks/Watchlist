package middleware

import (
	"net/http"
	"strings"
	"watchlist-backend/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Header se token lo
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Message: "token required",
			})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Token verify karo
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.Response{
				Success: false,
				Message: "invalid token",
			})
			c.Abort()
			return
		}

		// user_id context mein daalo
		claims := token.Claims.(jwt.MapClaims)
		userID := int(claims["user_id"].(float64))
		c.Set("user_id", userID)
		c.Next()
	}
}