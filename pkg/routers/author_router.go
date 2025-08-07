package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthorRoutes(router *gin.Engine, ctrl *controllers.AuthorController) {
	group := router.Group("/api/v1")
	{
		group.GET("/authors", ctrl.GetAll)
		group.GET("/authors/:id", ctrl.GetByID)
		group.POST("/authors", ctrl.Create)
		group.PUT("/authors/:id", ctrl.Update)
		group.DELETE("/authors/:id", ctrl.Delete)
	}
}