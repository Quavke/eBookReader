package middlewares

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BooksMiddleware(repo repositories.UserRepo) gin.HandlerFunc {
	return func (c *gin.Context) {
		claims, exists := c.Get("claims")
        if !exists {
            c.JSON(http.StatusUnauthorized, models.APIResponse[any]{Message: "something went wrong. You may not be logged in."})
						log.Println("Books middleware error, cannot find claims")
            c.Abort()
            return
        }
        
    userClaims := claims.(*models.Claims)

		isAuthor, err := repo.IsAuthor(userClaims.UserID)
		
		if err != nil && !isAuthor {
			c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "something went wrong. You may not be the author"})
			log.Printf("Books middleware error, user repo method isAuthor. Error: %s", err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}