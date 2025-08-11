package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(group *gin.RouterGroup, ctrl *controllers.BookController, AuthMiddleware gin.HandlerFunc){
	group.GET("/books", ctrl.GetAll)
	group.GET("/books/:id", ctrl.GetByID)
	auth := group.Group("/")
	auth.Use(AuthMiddleware)
	{
		auth.POST("/books", ctrl.Create)
		auth.PUT("/books/:id", ctrl.Update)
		auth.DELETE("/books/:id", ctrl.Delete)
	}
}