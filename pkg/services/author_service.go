package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
)

type AuthorService interface {
	GetAllAuthors()       									     (*[]models.Author, error)
	GetAuthorByID(id int) 									     (*models.Author, error)
	CreateAuthor(author *models.Author)          error
	UpdateAuthor(author *models.Author, id int)  error
	DeleteAuthor(id int)                         error
}

type AuthorServiceImpl struct {
	repo repositories.AuthorRepo
}

func NewAuthorService(repo repositories.AuthorRepo) *AuthorServiceImpl{
	return &AuthorServiceImpl{repo: repo}
}

func (s *AuthorServiceImpl) GetAllAuthors() (*[]models.Author, error){
	authors, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return &authors, nil
}

func (s *AuthorServiceImpl) GetAuthorByID(id int) (*models.Author, error){
	author, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return author, nil
}

func (s *AuthorServiceImpl) CreateAuthor(author *models.Author) error{
	return s.repo.Create(author)
}

func (s *AuthorServiceImpl) UpdateAuthor(author *models.Author, id int) error{
	return s.repo.Update(author, id)
}

func (s *AuthorServiceImpl) DeleteAuthor(id int) error{
	return s.repo.Delete(id)
}