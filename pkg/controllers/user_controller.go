package controllers

import (
	"ebookr/pkg/models"
	"ebookr/pkg/services"
	"ebookr/pkg/utils"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(service services.UserService) *UserController{
	return &UserController{UserService: service}
}

func (ctrl *UserController) GetAll(c *gin.Context){
	users, err := ctrl.UserService.GetAllUsers()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot get all users"})
		log.Printf("User controller GetAll error, service method GetAllUsers. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success", "error": nil, "users": users})
}

func (ctrl *UserController) GetByID(c *gin.Context){
	user, err := ctrl.UserService.GetUserByID(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot get all users"})
		log.Printf("User controller GetAll error, service method GetUserByID. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successful", "error": nil, "user": user})
}

func (ctrl *UserController) Create(c *gin.Context){
	var user models.RegisterReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message":"error", "error": "something wrong with your request. You need to sent Username and Password"})
		log.Printf("User controller Create error, bind. Error: %s", err.Error())
		return
	}

	if err := ctrl.UserService.CreateUser(user.Username, []byte(user.Password)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot create user"})
		log.Printf("User controller Create error, repo method CreateUser. Error: %s", err.Error())
		return
	}
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully", "error": nil})
}

func (ctrl *UserController) Login(c *gin.Context){
	var user models.RegisterReq
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "something wrong with your request. You need to sent Username and Password"})
		log.Printf("User controller Login error, bind. Error: %s", err.Error())
		return
	}

	claims, err := ctrl.UserService.LoginUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot login user"})
		log.Printf("User controller Login error, repo method LoginUser. Error: %s", err.Error())
		return
	}

	token, err := utils.GenerateToken(claims, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot generate token"})
		log.Printf("User controller Login error, gen token. Error: %s", err.Error())
		return
	}
	c.SetCookie("Authorization", token, 86400, "/", "", false, true) 
	c.JSON(http.StatusOK, gin.H{"message": "Successful login", "error": nil})
}

func (ctrl *UserController) Update(c *gin.Context){
	claims := c.MustGet("claims").(*models.Claims)
	var user models.UpdateReq
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot create integer id"})
    log.Printf("User controller Update error, cast id to int. Error: %s", err.Error())
		return
	}

	if id != int(claims.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "error", "error": "It is not you!"})
    log.Println("User controller Update error, access denied, invalid ID")
		return
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "something wrong with your request. You need to sent Username and Password"})
    log.Printf("User controller Update error, bind. Error: %s", err.Error())
		return
	}
	if err := ctrl.UserService.UpdateUser(&user, claims.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot update user"})
    log.Printf("User controller Update error, service method UpdateUser. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful update", "error": nil})
}

func (ctrl *UserController) Delete(c *gin.Context){
	claims := c.MustGet("claims").(*models.Claims)
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot create integer id"})
    log.Printf("User controller Delete error, cast id to int. Error: %s", err.Error())
		return
	}

	if id != int(claims.UserID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "error", "error": "It is not you!"})
    log.Println("User controller Delete error, access denied, invalid ID")
		return
	}

	if err := ctrl.UserService.DeleteUser(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot create integer id"})
    log.Printf("User controller Delete error, service method DeleteUser. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "successful delete", "error": nil})
}