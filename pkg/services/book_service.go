package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Quavke/eBookReader/pkg/models"
	"github.com/Quavke/eBookReader/pkg/repositories"
	"github.com/redis/go-redis/v9"
)

type BookService interface {
	GetAllBooks(limit, page uint, sort string)       									  	(*models.Pagination, error)
	GetBookByID(id uint) 									  						(*models.BookResp, error)
	CreateBook(book *models.Book)           								 error
	UpdateBook(book *models.Book, id, userID uint)   error
	DeleteBook(id uint, userID uint)                      error
}

type BookServiceImpl struct {
	repo repositories.BookRepo
	context context.Context
	redisClient *redis.Client
}

func NewBookService(repo repositories.BookRepo, context context.Context, redisClient *redis.Client) *BookServiceImpl{
	return &BookServiceImpl{
		repo: repo,
		context: context,
		redisClient: redisClient,
	}
}

var _ BookService = (*BookServiceImpl)(nil)

func (s *BookServiceImpl) GetAllBooks(limit, page uint, sort string) (*models.Pagination, error){
	cacheKey := fmt.Sprintf("books:limit=%d,page=%d,sort=%s", limit, page, sort)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
	if err == nil && cachedData != "" {
		var p models.Pagination
		if err := json.Unmarshal([]byte(cachedData), &p); err == nil {
			return &p, nil
		}
	}
	p := &models.Pagination{
		Limit: limit,
		Page: page,
		Sort: sort,
	}
	p, err = s.repo.GetAll(p)
	if err != nil {
		return nil, err
	}
	rows := p.Rows.([]models.Book)
	books := make([]models.BookResp, 0, len(rows))
	for _, b := range rows {
		books = append(books, models.BookResp{
			ID: b.ID,
			Title: b.Title,
			Content: b.Content,
			AuthorID: b.AuthorID,
		})
	}
	p.Rows = books

	data, err := json.Marshal(p)
  if err == nil {
      s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
			log.Print("Cached books data")
  }

	return p, nil
}

func (s *BookServiceImpl) GetBookByID(id uint) (*models.BookResp, error) {
	cacheKey := fmt.Sprintf("book:%d", id)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
	if err == nil && cachedData != "" {
		var book models.BookResp
		if err := json.Unmarshal([]byte(cachedData), &book); err == nil {
			return &book, nil
		}
	}
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
	data, err := json.Marshal(book)
  if err == nil {
      s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
			log.Print("Cached book data")
  }
	return book, nil
}

func (s *BookServiceImpl) CreateBook(book *models.Book) error{
	if len(book.Title) < 3 || len(book.Title) > 400 && len(book.Content) < 10 {
		return errors.New("title must be between 3 and 400 characters and content must be at least 10 characters")
	}
	cacheKey := fmt.Sprintf("create_book:%d,title=%s", book.AuthorID, book.Title)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
	if err == nil && cachedData != "" {
		var result string
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			if result == "success" {
				return nil
			} else {
				return fmt.Errorf("cached error: %s", result)
			}
		}
	}
	createResult := s.repo.Create(book)
	var result string
	if createResult != nil {
		result = createResult.Error()
	} else {
		result = "success"
	}
	data, err := json.Marshal(result)
	if err == nil {
      s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
			log.Print("Cached create book data")
  }
	return createResult
}

func (s *BookServiceImpl) UpdateBook(book *models.Book, id, userID uint) error {
	cacheKey := fmt.Sprintf("update_book:%d,title=%s", book.AuthorID, book.Title)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
	if err == nil && cachedData != "" {
		var result string
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			if result == "success" {
				return nil
			} else {
				return fmt.Errorf("cached error: %s", result)
			}
		}
	}

	isBelongs, err := s.repo.IsBelongsTo(id, userID)
	if err != nil && !isBelongs{
		return err
	}


	updateResult :=  s.repo.Update(book, id)
	var result string
	if updateResult != nil {
		result = updateResult.Error()
	} else {
		result = "success"
	}
	data, err := json.Marshal(result)
	if err == nil {
      s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
			log.Print("Cached update author data")
  }
	return updateResult
}

func (s *BookServiceImpl) DeleteBook(id uint, userID uint) error {
	cacheKey := fmt.Sprintf("delete_book:%d", id)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
	if err == nil && cachedData != "" {
		var result string
		if err := json.Unmarshal([]byte(cachedData), &result); err == nil {
			if result == "success" {
				return nil
			} else {
				return fmt.Errorf("cached error: %s", result)
			}
		}
	}

	isBelongs, err := s.repo.IsBelongsTo(id, userID)
	if err != nil && !isBelongs{
		return err
	}

	deleteResult := s.repo.Delete(id)
	var result string
	if deleteResult != nil {
		result = deleteResult.Error()
	} else {
		result = "success"
	}
	data, err := json.Marshal(result)
	if err == nil {
      s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
			log.Print("Cached delete author data")
  }
	return deleteResult
}