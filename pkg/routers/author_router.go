package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthorRoutes(group *gin.RouterGroup, ctrl *controllers.AuthorController) {
	{
		group.GET("/authors", ctrl.GetAll)
		group.GET("/authors/:id", ctrl.GetByID)
		group.POST("/authors", ctrl.Create)
		group.PUT("/authors/:id", ctrl.Update)
		group.DELETE("/authors/:id", ctrl.Delete)
	}
}