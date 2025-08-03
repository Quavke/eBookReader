package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
)

type BookService interface {
	GetAllBooks()       									  (*[]models.Book, error)
	GetBookByID(id int) 									  (*models.Book, error)
	CreateBook(book *models.Book)           error
	UpdateBook(book *models.Book, id int)   error
	DeleteBook(id int)                      error
}

type BookServiceImpl struct {
	repo repositories.BookRepo
}

func NewBookService(repo repositories.BookRepo) *BookServiceImpl{
	return &BookServiceImpl{repo: repo}
}

func (s *BookServiceImpl) GetAllBooks() (*[]models.Book, error){
	books, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return &books, nil
}

func (s *BookServiceImpl) GetBookByID(id int) (*models.Book, error) {
	book, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (s *BookServiceImpl) CreateBook(book *models.Book) error{
	if err := s.repo.Create(book); err != nil {
		return err
	}
	return nil
}

func (s *BookServiceImpl) UpdateBook(book *models.Book, id int) error {
	if err := s.repo.Update(book, id); err != nil {
		return err
	}
	return nil
}

func (s *BookServiceImpl) DeleteBook(id int) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	return nil
}