package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
)

type BookService interface {
	GetAllBooks()       									  (*[]models.BookResp, error)
	GetBookByID(bookID int) 									  (*models.BookResp, error)
	CreateBook(book *models.Book)           error
	UpdateBook(book *models.Book, bookID int, userID uint)   error
	DeleteBook(bookID int, userID uint)                      error
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

func (s *BookServiceImpl) GetBookByID(bookID int) (*models.BookResp, error) {
	bookDB, err := s.repo.GetByID(bookID)
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

func (s *BookServiceImpl) UpdateBook(book *models.Book, bookID int, userID uint) error {
	result, err := s.repo.IsBelongsTo(bookID, userID)
	if err != nil && !result{
		return err
	}
	return s.repo.Update(book, bookID)
}

func (s *BookServiceImpl) DeleteBook(bookID int, userID uint) error {
	result, err := s.repo.IsBelongsTo(bookID, userID)
	if err != nil && !result{
		return err
	}
	return s.repo.Delete(bookID)
}