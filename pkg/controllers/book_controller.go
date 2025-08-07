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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't get all books"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful", "error": nil, "books": books})
}

func (ctrl *BookController) GetByID(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't create integer id"})
		log.Println(err.Error())
		return
	}
	book, err := ctrl.BookService.GetBookByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't get book by this id"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful create", "error": nil, "book": book})
}

func (ctrl *BookController) Create(c *gin.Context){
	var book models.Book

	if err := c.ShouldBindBodyWithJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "something wrong with book"})
		log.Println(err.Error())
		return
	}

	if err := ctrl.BookService.CreateBook(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "can't create book"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful create", "error": nil})
}

func (ctrl *BookController) Update(c *gin.Context){
	var book models.Book

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't create integer id"})
		log.Println(err.Error())
		return
	}

	if err := c.ShouldBindBodyWithJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "something wrong with book"})
		log.Println(err.Error())
		return
	}
	if err := ctrl.BookService.UpdateBook(&book, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "can't update book"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful update", "error": nil})
}

func (ctrl *BookController) Delete(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't cast id to integer"})
		log.Println(err.Error())
		return
	}

	if err := ctrl.BookService.DeleteBook(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't delete book by this id"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful delete", "error": nil})
}
