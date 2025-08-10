package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
)

type BookService interface {
	GetAllBooks()       									  (*[]models.BookResp, error)
	GetBookByID(id int) 									  (*models.BookResp, error)
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

func (s *BookServiceImpl) GetAllBooks() (*[]models.BookResp, error){
	booksDB, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	books := make([]models.BookResp, 0, len(booksDB))
	for _, b := range booksDB {
		books = append(books, models.BookResp{
			ID: b.ID,
			Title: b.Title,
			Content: b.Content,
			AuthorID: b.AuthorID,
		})
	}
	return &books, nil
}

func (s *BookServiceImpl) GetBookByID(id int) (*models.BookResp, error) {
	bookDB, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	book := &models.BookResp{
		ID: bookDB.ID,
		Title: bookDB.Title,
		Content: bookDB.Content,
		AuthorID: bookDB.AuthorID,
	}
	return book, nil
}

func (s *BookServiceImpl) CreateBook(book *models.Book) error{
	return s.repo.Create(book)
}

func (s *BookServiceImpl) UpdateBook(book *models.Book, id int) error {
	return s.repo.Update(book, id)
}

func (s *BookServiceImpl) DeleteBook(id int) error {
	return s.repo.Delete(id)
}