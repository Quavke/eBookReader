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
	limitStr := c.DefaultQuery("l", "50")
	pageStr := c.DefaultQuery("p", "1")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create integer limit"})
		log.Printf("Author controller GetAll error, cast limit to int. Error: %s", err.Error())
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create integer page"})
		log.Printf("Author controller GetAll error, cast page to int. Error: %s", err.Error())
		return
	}
	users, err := ctrl.UserService.GetAllUsers(limit, page, "id desc")
	if err != nil{
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot get all users"})
		log.Printf("User controller GetAll error, service method GetAllUsers. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "Success", Data: users})
}

func (ctrl *UserController) GetByID(c *gin.Context){
	user, err := ctrl.UserService.GetUserByID(1)
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot user by this ID"})
		log.Printf("User controller GetAll error, service method GetUserByID. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "Successful", Data: user})
}

func (ctrl *UserController) Create(c *gin.Context){
	var user models.RegisterReq
	if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "something wrong with your request. You need to sent Username and Password"})
		log.Printf("User controller Create error, bind. Error: %s", err.Error())
		return
	}

	if err := ctrl.UserService.CreateUser(user.Username, []byte(user.Password)); err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create user"})
		log.Printf("User controller Create error, repo method CreateUser. Error: %s", err.Error())
		return
	}
	user.Password = ""

  c.JSON(http.StatusCreated, models.APIResponse[any]{Message: "User created successfully"})
}

func (ctrl *UserController) Login(c *gin.Context){
	var user models.RegisterReq
	if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "something wrong with your request. You need to sent Username and Password"})
		log.Printf("User controller Login error, bind. Error: %s", err.Error())
		return
	}

	claims, err := ctrl.UserService.LoginUser(&user)
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot login user"})
		log.Printf("User controller Login error, repo method LoginUser. Error: %s", err.Error())
		return
	}

	token, err := utils.GenerateToken(claims, []byte(os.Getenv("JWT_SECRET")))
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot generate token"})
		log.Printf("User controller Login error, gen token. Error: %s", err.Error())
		return
	}
	c.SetCookie("Authorization", token, 86400, "/", "", false, true) 
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "Successful login"})
}

func (ctrl *UserController) Logout(c *gin.Context){
	isProd := c.MustGet("isProd").(bool)
	c.SetCookie("Authorization", "", -1, "/", "", isProd, true)

	c.Status(http.StatusNoContent)
}

func (ctrl *UserController) Update(c *gin.Context){
	claims := c.MustGet("claims").(*models.Claims)
	var user models.UpdateReq

	if err := c.ShouldBindJSON(&user); err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "something wrong with your request. You need to sent Username and Password"})
    log.Printf("User controller Update error, bind. Error: %s", err.Error())
		return
	}
	if err := ctrl.UserService.UpdateUser(&user, claims.UserID); err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot update user"})
    log.Printf("User controller Update error, service method UpdateUser. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful update"})
}

func (ctrl *UserController) Delete(c *gin.Context){
	claims := c.MustGet("claims").(*models.Claims)
	isProd := c.MustGet("isProd").(bool)
	if err := ctrl.UserService.DeleteUser(claims.UserID); err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create integer id"})
    log.Printf("User controller Delete error, service method DeleteUser. Error: %s", err.Error())
		return
	}
	c.SetCookie("Authorization", "", -1, "/", "", isProd, true)
	c.Status(http.StatusNoContent)
}

func (ctrl *UserController) GetCreateMock(c *gin.Context) {
	users := make([]models.RegisterReq, 200)
	for i := range users {
    users[i] = models.RegisterReq{
        Username: "user" + strconv.Itoa(i+1),
				Password: "password" + strconv.Itoa(i+1),
    }
	}
	for i := range users {
		if err := ctrl.UserService.CreateUser(users[i].Username, []byte(users[i].Password)); err != nil {
			log.Printf("User controller GetCreateMock error, service method CreateUser. Error: %s", err.Error())
			c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create mock users"})
			return
		}
	}

	c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful create"})
}