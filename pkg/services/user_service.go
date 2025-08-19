package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers(limit, page int, sort string)       									  (*[]models.UserResp, error)
	GetUserByID(id uint) 									  (*models.UserResp, error)
	CreateUser(username string, pwd []byte)           error
	LoginUser(user *models.RegisterReq) (*models.Claims, error)
	UpdateUser(user *models.UpdateReq, id uint)   error
	DeleteUser(id uint)                      error
}

type UserServiceImpl struct {
	repo repositories.UserRepo
}

func NewUserService(repo repositories.UserRepo) *UserServiceImpl{
	return &UserServiceImpl{repo: repo}
}

var _ UserService = (*UserServiceImpl)(nil)

func (s UserServiceImpl) GetAllUsers(limit, page int, sort string) (*[]models.UserResp, error){
		p := &models.Pagination{
			Limit: limit,
			Page: page,
			Sort: sort,
		}
		p, err := s.repo.GetAll(p)
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
	return &users, nil
}

func (s *UserServiceImpl) GetUserByID(id uint) (*models.UserResp, error) {
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
	return user, nil
}

func (s *UserServiceImpl) CreateUser(username string, pwd []byte) error{
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
	return s.repo.Create(&userDB)
}

func (s *UserServiceImpl) LoginUser(user *models.RegisterReq) (*models.Claims, error) {
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
	return claims, nil
}

func (s *UserServiceImpl) UpdateUser(user *models.UpdateReq, id uint) error{
	return s.repo.Update(user, id)
}

func (s *UserServiceImpl) DeleteUser(id uint) error{
	return s.repo.Delete(id)
}