package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Quavke/eBookReader/pkg/models"
	"github.com/Quavke/eBookReader/pkg/repositories"
	"github.com/redis/go-redis/v9"
)

type AuthorService interface {
	GetAllAuthors(limit, page uint, sort string)  (*models.Pagination, error)
	GetAuthorByID(id uint) 									     (*models.AuthorResp, error)
	CreateAuthor(author *models.Author)          error
	UpdateAuthor(author *models.UpdateAuthorReq, id uint)  error
	DeleteAuthor(id uint)                         error
}

type AuthorServiceImpl struct {
	repo repositories.AuthorRepo
	context context.Context
	redisClient *redis.Client
}

func NewAuthorService(repo repositories.AuthorRepo, context context.Context, 	redisClient *redis.Client) *AuthorServiceImpl{
	return &AuthorServiceImpl{
		repo: repo,
		context: context,
		redisClient: redisClient,
	}
}

var _ AuthorService = (*AuthorServiceImpl)(nil)

func (s *AuthorServiceImpl) GetAllAuthors(limit, page uint, sort string) (*models.Pagination, error){
	cacheKey := fmt.Sprintf("authors:limit=%d,page=%d,sort=%s", limit, page, sort)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
	if err == nil && cachedData != "" {
		var p models.Pagination
		if err := json.Unmarshal([]byte(cachedData), &p); err == nil {
			return &p, nil
		}
	}

	p := &models.Pagination{
		Limit: uint(limit),
		Page: uint(page),
		Sort: sort,
	}
	p, err = s.repo.GetAll(p)
	if err != nil {
		return nil, err
	}
	rows := p.Rows.([]models.Author)
	authors := make([]models.AuthorResp, 0, len(rows))
	for _, a := range rows {
		authors = append(authors, models.AuthorResp{
			UserID: a.UserID,
			Firstname: a.Firstname,
			Lastname: a.Lastname,
			Birthday: a.Birthday,
		})
	}
	p.Rows = authors

	data, err := json.Marshal(p)
  if err == nil {
      s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
			log.Print("Cached authors data")
  }

	return p, nil
}

func (s *AuthorServiceImpl) GetAuthorByID(id uint) (*models.AuthorResp, error){
	cacheKey := fmt.Sprintf("author:%d", id)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
	if err == nil && cachedData != "" {
		var author models.AuthorResp
		if err := json.Unmarshal([]byte(cachedData), &author); err == nil {
			return &author, nil
		}
	}
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

	data, err := json.Marshal(author)
  if err == nil {
      s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
			log.Print("Cached author data")
  }

	return author, nil
}

func (s *AuthorServiceImpl) CreateAuthor(author *models.Author) error{
	if author.Firstname == "" && author.Lastname == "" && author.Birthday.IsZero() {
		return nil
	}
	cacheKey := fmt.Sprintf("create_author:%d", author.UserID)
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
	createResult := s.repo.Create(author)
	var result string
	if createResult != nil {
		result = createResult.Error()
	} else {
		result = "success"
	}
	data, err := json.Marshal(result)
	if err == nil {
      s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
			log.Print("Cached create author data")
  }
	return createResult
}

func (s *AuthorServiceImpl) UpdateAuthor(author *models.UpdateAuthorReq, id uint) error{
	cacheKey := fmt.Sprintf("update_author:firstname=%s,lastname=%s,birthday=%s", author.Firstname, author.Lastname, author.Birthday.Format("2006-01-02"))
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
	
	updateResult := s.repo.Update(author, id)
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

func (s *AuthorServiceImpl) DeleteAuthor(id uint) error{
	cacheKey := fmt.Sprintf("delete_author:%d", id)
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