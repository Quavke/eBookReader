package middlewares

import (
	"ebookr/pkg/models"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func AuthMiddleware(jwtSecretKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == ""{
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error", "error": "Authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error", "error": "Invalid authorization format. User 'Bearer token'"})
			return
		}

		tokenStr := parts[1]

		token, err := jwt.ParseWithClaims(tokenStr, &models.Claims{}, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return jwtSecretKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error", "error": err.Error()})
			return
		}

		if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
			user := &models.User{
				Model: gorm.Model{
					ID: uint(claims.UserID),
				},
				Username: claims.Username,
			}

			c.Set("user", user)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "error", "error": "Invalid token"})
		}
	}
}