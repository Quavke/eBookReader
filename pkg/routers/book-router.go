package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterBookRoutes(router *gin.Engine, ctrl *controllers.BookController){
	group := router.Group("/api/v1")
	{
		group.GET("/books", ctrl.GetAll)
		group.GET("/book/:id", ctrl.GetByID)
		group.POST("/books", ctrl.Create)
		group.PUT("/books/:id", ctrl.Update)
		group.DELETE("/books/:id", ctrl.Delete)
	}
}