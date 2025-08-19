package controllers

import (
	"github.com/Quavke/eBookReader/pkg/models"
	"github.com/Quavke/eBookReader/pkg/services"
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

	
	authors, err := ctrl.AuthorService.GetAllAuthors(limit, page, "user_id desc")
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

func (ctrl *AuthorController) Update(c *gin.Context){
	var author models.UpdateAuthorReq
	
	claims := c.MustGet("claims").(*models.Claims)


	if err := c.ShouldBindBodyWithJSON(&author); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "something wrong with your request. You need to sent Firstname or Lastname or Birthday(yyyy-mm-dd)"})
    log.Printf("Author controller Update error, bind. Error: %s", err.Error())
		return
	}

	if author.Firstname == "" && author.Lastname == "" && author.Birthday.IsZero() {
		c.JSON(http.StatusOK, models.APIResponse[any]{
				Message: "successful update", 
				Data: "No fields to update"})
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

func (ctrl *AuthorController) GetCreateMock(c *gin.Context) {
	authors := make([]models.Author, 200)
	for i := range authors {
    authors[i] = models.Author{
        UserID: uint(i+1),
				Firstname: "firstname" + strconv.Itoa(i+1),
				Lastname: "lastname" + strconv.Itoa(i+1),
				Birthday: models.DateOnly{Time: models.DateOnly{}.Time.AddDate(0, 0, i+1),
			},
		}
	}
	for i := range authors {
		if err := ctrl.AuthorService.CreateAuthor(&authors[i]); err != nil {
			log.Printf("User controller GetCreateMock error, service method CreateUser. Error: %s", err.Error())
			c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create mock users"})
			return
		}
	}

	c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful create"})
}
