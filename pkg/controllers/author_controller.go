package controllers

import (
	"ebookr/pkg/services"
)

type AuthorController struct {
	BookService services.BookService
}

func NewAuthorController(service services.BookService) *BookController {
	return &BookController{BookService: service}
}