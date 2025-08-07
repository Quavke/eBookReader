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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "Cannot get all authors"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful", "error": nil, "authors": authors})
}

func (ctrl *AuthorController) GetByID(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't create integer id"})
		log.Println(err.Error())
		return
	}
	author, err := ctrl.AuthorService.GetAuthorByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't get author by this id"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful create", "error": nil, "author": author})
}

func (ctrl *AuthorController) Create(c *gin.Context){
	var author models.Author

	if err := c.ShouldBindBodyWithJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "something wrong with author"})
		log.Println(err.Error())
		return
	}

	if err := ctrl.AuthorService.CreateAuthor(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "can't create author"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful create", "error": nil})
}

func (ctrl *AuthorController) Update(c *gin.Context){
	var author models.Author

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't create integer id"})
		log.Println(err.Error())
		return
	}

	if err := c.ShouldBindBodyWithJSON(&author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "something wrong with author"})
		log.Println(err.Error())
		return
	}
	if err := ctrl.AuthorService.UpdateAuthor(&author, id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error", "error": "can't update author"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful update", "error": nil})
}

func (ctrl *AuthorController) Delete(c *gin.Context){
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't cast id to integer"})
		log.Println(err.Error())
		return
	}

	if err := ctrl.AuthorService.DeleteAuthor(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error", "error": "can't delete author by this id"})
		log.Println(err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "successful delete", "error": nil})
}