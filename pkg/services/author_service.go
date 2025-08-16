package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
)

type AuthorService interface {
	GetAllAuthors()       									     (*[]models.AuthorResp, error)
	GetAuthorByID(id uint) 									     (*models.AuthorResp, error)
	CreateAuthor(author *models.Author)          error
	UpdateAuthor(author *models.Author, id uint)  error
	DeleteAuthor(id uint)                         error
}

type AuthorServiceImpl struct {
	repo repositories.AuthorRepo
}

func NewAuthorService(repo repositories.AuthorRepo) *AuthorServiceImpl{
	return &AuthorServiceImpl{repo: repo}
}

func (s *AuthorServiceImpl) GetAllAuthors() (*[]models.AuthorResp, error){
	authorsDB, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	authors := make([]models.AuthorResp, 0, len(authorsDB))
	for _, a := range authorsDB {
		authors = append(authors, models.AuthorResp{
			UserID: a.UserID,
			Firstname: a.Firstname,
			Lastname: a.Lastname,
			Birthday: a.Birthday,
		})
	}
	return &authors, nil
}

func (s *AuthorServiceImpl) GetAuthorByID(id uint) (*models.AuthorResp, error){
	authorBD, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	author := &models.AuthorResp{
		UserID: authorBD.UserID,
		Firstname: authorBD.Firstname,
		Lastname: authorBD.Lastname,
		Birthday: authorBD.Birthday,
	}
	return author, nil
}

func (s *AuthorServiceImpl) CreateAuthor(author *models.Author) error{
	return s.repo.Create(author)
}

func (s *AuthorServiceImpl) UpdateAuthor(author *models.Author, id uint) error{
	return s.repo.Update(author, id)
}

func (s *AuthorServiceImpl) DeleteAuthor(id uint) error{
	return s.repo.Delete(id)
}