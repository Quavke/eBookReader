package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Quavke/eBookReader/pkg/models"
	"github.com/Quavke/eBookReader/pkg/services"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	BookService services.BookService
}

func NewBookController(service services.BookService) *BookController {
	return &BookController{BookService: service}
}

func (ctrl *BookController) GetAll(c *gin.Context){
	limitStr := c.DefaultQuery("l", "50")
	pageStr := c.DefaultQuery("p", "1")

	limit, err := strconv.ParseUint(limitStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create integer limit"})
		log.Printf("Author controller GetAll error, cast limit to int. Error: %s", err.Error())
		return
	}

	page, err := strconv.ParseUint(pageStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create integer page"})
		log.Printf("Author controller GetAll error, cast page to int. Error: %s", err.Error())
		return
	}

	books, err := ctrl.BookService.GetAllBooks(uint(limit), uint(page), "id desc")
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot get all books"})
		log.Printf("Book controller GetAll error. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful", Data: books})
}

func (ctrl *BookController) GetByID(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create integer id"})
    log.Printf("Book controller GetByID error, cast id to int. Error: %s", err.Error())
		return
	}
	book, err := ctrl.BookService.GetBookByID(uint(id))
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot get book by this id"})
    log.Printf("Book controller GetByID error. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful create", Data: book})
}

func (ctrl *BookController) Create(c *gin.Context){
	var book models.Book

	claims := c.MustGet("claims").(*models.Claims)

	if err := c.ShouldBindJSON(&book); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "something wrong with your request. You need to sent Title(min 3 chars), Content(min 50 chars)"})
    log.Printf("Book controller Create error, bind. Error: %s", err.Error())
		return
	}

	book.AuthorID = claims.UserID

	if err := ctrl.BookService.CreateBook(&book); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "cannot create book"})
        log.Printf("Book controller Create error, service method CreateBook. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful create"})
}

func (ctrl *BookController) Update(c *gin.Context){
	var book models.Book

	claims := c.MustGet("claims").(*models.Claims)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot cast id to integer"})
    log.Printf("Book controller Update error, cast id to int. Error: %s", err.Error())
		return
	}


	if err := c.ShouldBindJSON(&book); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "something wrong with your request. You need to sent Title(min 3 chars), Content(min 50 chars)"})
    log.Printf("Book controller Update error, bind. Error: %s", err.Error())
		return
	}

	if len(book.Title) < 3 || len(book.Title) > 400 && len(book.Content) < 10 {
		c.JSON(http.StatusOK, models.APIResponse[any]{
				Message: "successful update", 
				Data: "No fields to update"})
		return
	}

	if err := ctrl.BookService.UpdateBook(&book, uint(id), claims.UserID); err != nil {
    c.JSON(http.StatusBadRequest, models.APIResponse[any]{Message: "error", Error: "cannot update book"})
    log.Printf("Book controller Update error, service method UpdateBook. Error: %s", err.Error())
		return
	}
  c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful update"})
}

func (ctrl *BookController) Delete(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot cast id to integer"})
    log.Printf("Book controller Delete error, cast id to int. Error: %s", err.Error())
		return
	}

	claims := c.MustGet("claims").(*models.Claims)

	if err := ctrl.BookService.DeleteBook(uint(id), claims.UserID); err != nil {
    c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot delete book by this id"})
    log.Printf("Book controller Delete error, service method DeleteBook. Error: %s", err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}

func (ctrl *BookController) GetCreateMock(c *gin.Context) {
	books := make([]models.Book, 200)
	for i := range books {
    books[i] = models.Book{
        AuthorID: uint(i+1),
				Title: "Book title" + strconv.Itoa(i+1),
				Content: "Content Content Content Content Content Content Content Content" + strconv.Itoa(i+1),
			}
	}
	for i := range books {
		if err := ctrl.BookService.CreateBook(&books[i]); err != nil {
			log.Printf("User controller GetCreateMock error, service method CreateUser. Error: %s", err.Error())
			c.JSON(http.StatusInternalServerError, models.APIResponse[any]{Message: "error", Error: "cannot create mock users"})
			return
		}
	}

	c.JSON(http.StatusOK, models.APIResponse[any]{Message: "successful create"})
}
