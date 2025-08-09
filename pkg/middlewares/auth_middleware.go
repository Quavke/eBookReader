package middlewares

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(repo repositories.UserRepo) gin.HandlerFunc {
	return func (c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods([]string{"HS256"}), jwt.WithIssuer("eBookReader"),)

		if err != nil || !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		
		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > float64(claims.ExpiresAt.Unix()){
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			if float64(time.Now().Unix()) < float64(claims.NotBefore.Unix()){
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			if err := repo.IsExists(uint(claims.UserID)); err != nil{
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		}

		c.Set("claims", claims)
		c.Next()
	}
}