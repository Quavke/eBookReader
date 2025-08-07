package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetAllUsers()       									  (*[]models.UserDB, error)
	GetUserByID(id uint) 									  (*models.UserDB, error)
	CreateUser(user *models.UserDB)           error
	UpdateUser(user *models.UserDB, id uint)   error
	DeleteUser(id uint)                      error
}

type UserServiceImpl struct {
	repo repositories.UserRepo
}

func NewUserService(repo repositories.UserRepo) *UserServiceImpl{
	return &UserServiceImpl{repo: repo}
}

func (s UserServiceImpl) GetAllUsers() (*[]models.UserDB, error){
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (s *UserServiceImpl) GetUserByID(id uint) (*models.UserDB, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) CreateUser(user *models.UserDB) error{
	var userDB models.UserDB
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	userDB.Password = hash
	return s.repo.Create(&userDB)
}

func (s *UserServiceImpl) LoginUser(user *models.UserDB) error {
	userDB, err := s.repo.GetByID(user.ID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword(userDB.Password, []byte(user.Password)); err != nil {
		return err
	}
	return nil
}

func (s *UserServiceImpl) UpdateUser(user *models.UserDB, id uint) error{
	return s.repo.Update(user, id)
}

func (s *UserServiceImpl) DeleteUser(id uint) error{
	return s.repo.Delete(id)
}