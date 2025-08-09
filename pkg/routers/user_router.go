package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(group *gin.RouterGroup, ctrl *controllers.UserController){ // , jwt gin.HandlerFunc
	{
		group.GET("/users", ctrl.GetAll)
		group.GET("/users/:id", ctrl.GetByID)
		group.POST("/users", ctrl.Create)
		group.POST("/users/login", ctrl.Login)
		group.PUT("/users/:id", ctrl.Update)
		group.DELETE("/users/:id", ctrl.Delete)
	}
}