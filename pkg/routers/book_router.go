package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(group *gin.RouterGroup, ctrl *controllers.BookController){
	{
		group.GET("/books", ctrl.GetAll)
		group.GET("/books/:id", ctrl.GetByID)
		group.POST("/books", ctrl.Create)
		group.PUT("/books/:id", ctrl.Update)
		group.DELETE("/books/:id", ctrl.Delete)
	}
}