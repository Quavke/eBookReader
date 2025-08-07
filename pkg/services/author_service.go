package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
)

type AuthorService interface {
	GetAllAuthors()       									  (*[]models.Author, error)
	GetAuthorByID(id int) 									  (*models.Author, error)
	CreateAuthor(author *models.Author)           error
	UpdateAuthor(author *models.Author, id int)   error
	DeleteAuthor(id int)                      error
}

type AuthorServiceImpl struct {
	repo repositories.AuthorRepo
}

func NewAuthorService(repo repositories.AuthorRepo) *AuthorServiceImpl{
	return &AuthorServiceImpl{repo: repo}
}