package routers

import (
	"ebookr/pkg/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(group *gin.RouterGroup, ctrl *controllers.UserController, AuthMiddleware gin.HandlerFunc){
	group.POST("/users/login", ctrl.Login)
	group.POST("/users", ctrl.Create)
	group.POST("/users", ctrl.Logout)
	group.GET("/users", ctrl.GetAll)
	group.GET("/users/:id", ctrl.GetByID)
	auth := group.Group("/")
	auth.Use(AuthMiddleware)
	{
		auth.PUT("/users/me", ctrl.Update)
		auth.DELETE("/users/me", ctrl.Delete)
	}
}