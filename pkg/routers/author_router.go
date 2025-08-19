package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthorRoutes(group *gin.RouterGroup, ctrl *controllers.AuthorController, AuthMiddleware gin.HandlerFunc) {
	group.GET("/authors", ctrl.GetAll)
	group.GET("/authors/:id", ctrl.GetByID)
	group.GET("/authors/create", ctrl.GetCreateMock)
	auth := group.Group("/")
	auth.Use(AuthMiddleware)
	{
		auth.POST("/authors", ctrl.Create)
		auth.PUT("/authors/me", ctrl.Update)
		auth.DELETE("/authors/me", ctrl.Delete)
	}
}