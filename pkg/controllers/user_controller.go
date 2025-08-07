package controllers

import (
	"ebookr/pkg/models"
	"ebookr/pkg/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

// TODO Написать реализацию контроллеров

func NewUserController(service services.UserService) *UserController{
	return &UserController{UserService: service}
}

func (ctrl *UserController) GetAll(c *gin.Context){
	ctrl.UserService.GetAllUsers()
}

func (ctrl *UserController) GetByID(c *gin.Context){
	ctrl.UserService.GetUserByID(1)
}

func (ctrl *UserController) Create(c *gin.Context){
	ctrl.UserService.CreateUser(&models.User{})
}

func (ctrl *UserController) Update(c *gin.Context){
	ctrl.UserService.UpdateUser(&models.User{}, 1)
}

func (ctrl *UserController) Delete(c *gin.Context){
	ctrl.UserService.DeleteUser(1)
}