package controllers

import (
	"ebookr/pkg/models"
	"ebookr/pkg/services"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookController struct {
	BookService services.BookService
}

func NewBookController(service services.BookService) *BookController {
	return &BookController{BookService: service}
}

func (ctrl *BookController) GetAll(c *gin.Context){
	books, err := ctrl.BookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot get all books"})
		log.Printf("Book controller GetAll error. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful", "error": nil, "books": books})
}

func (ctrl *BookController) GetByID(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot create integer id"})
    log.Printf("Book controller GetByID error, cast id to int. Error: %s", err.Error())
		return
	}
	book, err := ctrl.BookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot get book by this id"})
    log.Printf("Book controller GetByID error. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful create", "error": nil, "book": book})
}

func (ctrl *BookController) Create(c *gin.Context){
	var book models.Book

	if err := c.ShouldBindBodyWithJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "something wrong with your request. You need to sent Title(min 3 chars), Content(min 50 chars), Author"})
    log.Printf("Book controller Create error, bind. Error: %s", err.Error())
		return
	}

	if err := ctrl.BookService.CreateBook(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "cannot create book"})
        log.Printf("Book controller Create error, service method CreateBook. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful create", "error": nil})
}

func (ctrl *BookController) Update(c *gin.Context){
	var book models.Book

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot cast id to integer"})
    log.Printf("Book controller Update error, cast id to int. Error: %s", err.Error())
		return
	}

	if err := c.ShouldBindBodyWithJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "something wrong with your request. You need to sent Title(min 3 chars), Content(min 50 chars), Author"})
    log.Printf("Book controller Update error, bind. Error: %s", err.Error())
		return
	}
	if err := ctrl.BookService.UpdateBook(&book, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "cannot update book"})
    log.Printf("Book controller Update error, service method UpdateBook. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful update", "error": nil})
}

func (ctrl *BookController) Delete(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot cast id to integer"})
    log.Printf("Book controller Delete error, cast id to int. Error: %s", err.Error())
		return
	}

	if err := ctrl.BookService.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "cannot delete book by this id"})
    log.Printf("Book controller Delete error, service method DeleteBook. Error: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful delete", "error": nil})
}
