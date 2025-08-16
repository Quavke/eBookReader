package controllers

import (
	"ebookr/pkg/models"
	"ebookr/pkg/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthorController struct {
	AuthorService services.AuthorService
}

func NewAuthorController(service services.AuthorService) *AuthorController {
	return &AuthorController{AuthorService: service}
}

func (ctrl *AuthorController) GetAll(c *gin.Context){
	authors, err := ctrl.AuthorService.GetAllAuthors()
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot get all authors"})
    log.Printf("Author controller GetAll error, service method GetAllAuthors. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful", Data: authors})
}

func (ctrl *AuthorController) GetByID(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create integer id"})
    log.Printf("Author controller GetByID error, cast id to int. Error: %s", err.Error())
		return
	}
	author, err := ctrl.AuthorService.GetAuthorByID(uint(id))
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot get author by this id"})
    log.Printf("Author controller GetByID error, service method GetAuthorByID. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful create", Data: author})
}

func (ctrl *AuthorController) Create(c *gin.Context){
	var author models.Author

	user := c.MustGet("claims").(*models.Claims)

	if err := c.ShouldBindBodyWithJSON(&author); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "something wrong with your request. You need to sent Firstname, Lastname, Birthday(yyyy-mm-dd)"})
    log.Printf("Author controller Create error, bind. Error: %s", err.Error())
		return
	}
	author.UserID = user.UserID
	if err := ctrl.AuthorService.CreateAuthor(&author); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "cannot create author"})
    log.Printf("Author controller Create error, service method CreateAuthor. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful create"})
}
// TODO Delete и Update проверка на то, что пользователь владеет данными
func (ctrl *AuthorController) Update(c *gin.Context){
	var author models.Author
	
	claims := c.MustGet("claims").(*models.Claims)


	if err := c.ShouldBindBodyWithJSON(&author); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "something wrong with your request. You need to sent Firstname, Lastname, Birthday(yyyy-mm-dd)"})
    log.Printf("Author controller Update error, bind. Error: %s", err.Error())
		return
	}
	if err := ctrl.AuthorService.UpdateAuthor(&author, claims.UserID); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "cannot update author"})
    log.Printf("Author controller Update error, service method UpdateAuthor. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful update"})
}

func (ctrl *AuthorController) Delete(c *gin.Context){
	claims := c.MustGet("claims").(*models.Claims)

	if err := ctrl.AuthorService.DeleteAuthor(claims.UserID); err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "can't delete author by this id"})
    log.Printf("Author controller Delete error, service method DeleteAuthor. Error: %s", err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}