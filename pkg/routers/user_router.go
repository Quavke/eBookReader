package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.Engine, ctrl *controllers.UserController){ // , jwt gin.HandlerFunc
	group := router.Group("/api/v1")
	{
		group.GET("/users", ctrl.GetAll)
		group.GET("/users/:id", ctrl.GetByID)
		group.POST("/users", ctrl.Create)
		group.PUT("/users/:id", ctrl.Update)
		group.DELETE("/users/:id", ctrl.Delete)
	}
}