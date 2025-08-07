package controllers

import (
	"ebookr/pkg/models"
	"ebookr/pkg/services"
	"net/http"

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
	var user models.RegisterReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var userDB models.UserDB
	userDB.Username = user.Username
	userDB.Password = []byte(user.Password)
	user.Password = ""
	if err := ctrl.UserService.CreateUser(&userDB); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "error": nil})
}

func (ctrl *UserController) Update(c *gin.Context){
	ctrl.UserService.UpdateUser(&models.UserDB{}, 1)
}

func (ctrl *UserController) Delete(c *gin.Context){
	ctrl.UserService.DeleteUser(1)
}