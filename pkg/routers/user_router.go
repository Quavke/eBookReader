package routers

import (
	"ebookr/pkg/controllers"
	"ebookr/pkg/middlewares"
	"ebookr/pkg/repositories"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(group *gin.RouterGroup, ctrl *controllers.UserController, repo repositories.UserRepo){
	group.POST("/users/login", ctrl.Login)
	group.POST("/users", ctrl.Create)
	auth := group.Group("/")
	auth.Use(middlewares.AuthMiddleware(repo))
	{
		auth.GET("/users", ctrl.GetAll)
		auth.GET("/users/:id", ctrl.GetByID)
		auth.PUT("/users/:id", ctrl.Update)
		auth.DELETE("/users/:id", ctrl.Delete)
	}
}