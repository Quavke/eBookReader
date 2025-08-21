package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Quavke/eBookReader/pkg/models"
	"github.com/Quavke/eBookReader/pkg/repositories"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers(limit, page uint, sort string)      (*models.Pagination, error)
	GetUserByID(id uint) 									  				(*models.UserResp, error)
	CreateUser(username string, pwd []byte)         error
	LoginUser(user *models.RegisterReq) 						(*models.Claims, error)
	UpdateUser(user *models.UpdateReq, id uint)   	error
	DeleteUser(id uint)                      				error
}

type UserServiceImpl struct {
	repo repositories.UserRepo
	context context.Context
	redisClient *redis.Client
}

func NewUserService(repo repositories.UserRepo, context context.Context, redisClient *redis.Client) *UserServiceImpl{
	return &UserServiceImpl{
		repo: repo,
		context: context,
		redisClient: redisClient,
	}
}

var _ UserService = (*UserServiceImpl)(nil)

func (s UserServiceImpl) GetAllUsers(limit, page uint, sort string) (*models.Pagination, error){
	cacheKey := fmt.Sprintf("users:limit=%d,page=%d,sort=%s", limit, page, sort)
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
  users := make([]models.UserResp, 0, len(p.Rows.([]models.UserDB)))

  ids := make([]uint, 0, len(p.Rows.([]models.UserDB)))
  for _, u := range p.Rows.([]models.UserDB) {
      ids = append(ids, u.ID)
  }

  isAuthor, err := s.repo.IsAuthors(ids)
  if err != nil {
      return nil, err
  }
  for _, u := range p.Rows.([]models.UserDB) {
    users = append(users, models.UserResp{
        Username: u.Username,
        IsAuthor: isAuthor[u.ID],
    })
  }

	p.Rows = users

	data, err := json.Marshal(p)
	if err == nil {
		s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
		log.Print("Cached users data")
	}

	return p, nil
}

func (s *UserServiceImpl) GetUserByID(id uint) (*models.UserResp, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
		if err == nil && cachedData != "" {
			var user models.UserResp
			if err := json.Unmarshal([]byte(cachedData), &user); err == nil {
				return &user, nil
			}
	}
	userDB, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	isAuthor, err := s.repo.IsAuthor(userDB.ID)
	if err != nil {
		return nil, err
	}

	user := &models.UserResp{
		Username: userDB.Username,
		IsAuthor: isAuthor,
	}

	data, err := json.Marshal(user)
	if err == nil {
		s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
		log.Print("Cached users data")
	}

	return user, nil
}

func (s *UserServiceImpl) CreateUser(username string, pwd []byte) error{
	cacheKey := fmt.Sprintf("user:username=%s", username)
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
	var userDB models.UserDB
	hash, err := bcrypt.GenerateFromPassword(pwd, 12)
	if err != nil {
		return err
	}
	for i := range pwd {
		pwd[i] = 0
	}
	userDB.Username = username
	userDB.PasswordHash = hash

	createResult := s.repo.Create(&userDB)
	var result string
	if createResult != nil {
		result = createResult.Error()
	} else {
		result = "success"
	}
	data, err := json.Marshal(result)
	if err == nil {
		s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
		log.Print("Cached create user data")
	}

	return createResult
}

func (s *UserServiceImpl) LoginUser(user *models.RegisterReq) (*models.Claims, error) {
	cacheKey := fmt.Sprintf("login_user:username=%s", user.Username)
	cachedData, err := s.redisClient.Get(s.context, cacheKey).Result()
		if err == nil && cachedData != "" {
			var claims models.Claims
			if err := json.Unmarshal([]byte(cachedData), &claims); err == nil {
				return &claims, nil
			}
	}
	userDB, err := s.repo.GetByUsername(user.Username)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword(userDB.PasswordHash, []byte(user.Password)); err != nil {
		return nil, err
	}
	user.Password = ""
	claims := &models.Claims{
    UserID:   uint(userDB.ID),
    Username: userDB.Username,
    RegisteredClaims: jwt.RegisteredClaims{
      ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
      IssuedAt:  jwt.NewNumericDate(time.Now()),
      Issuer:    "eBookReader",
      Subject:   fmt.Sprintf("%d", userDB.ID),
    },
  }
	data, err := json.Marshal(claims)
	if err == nil {
		s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
		log.Print("Cached login user data")
	}
	return claims, nil
}

func (s *UserServiceImpl) UpdateUser(user *models.UpdateReq, id uint) error{
	cacheKey := fmt.Sprintf("update_user:%d,username=%s", id, user.Username)
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
	updateResult := s.repo.Update(user, id)
	var result string
	if updateResult != nil {
		result = updateResult.Error()
	} else {
		result = "success"
	}
	data, err := json.Marshal(result)
	if err == nil {
		s.redisClient.Set(s.context, cacheKey, data, 5 * time.Minute)
		log.Print("Cached update user data")
	}
	return updateResult
}

func (s *UserServiceImpl) DeleteUser(id uint) error{
	cacheKey := fmt.Sprintf("delete_user:%d", id)
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
		log.Print("Cached delete user data")
	}
	return deleteResult
}