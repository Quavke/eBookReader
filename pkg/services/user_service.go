package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
)

type UserService interface {
	GetAllUsers()       									  (*[]models.User, error)
	GetUserByID(id int) 									  (*models.User, error)
	CreateUser(user *models.User)           error
	UpdateUser(user *models.User, id int)   error
	DeleteUser(id int)                      error
}

type UserServiceImpl struct {
	repo repositories.UserRepo
}

func NewUserService(repo repositories.UserRepo) *UserServiceImpl{
	return &UserServiceImpl{repo: repo}
}

func (s UserServiceImpl) GetAllUsers() (*[]models.User, error){
	users, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (s *UserServiceImpl) GetUserByID(id int) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserServiceImpl) CreateUser(user *models.User) error{
	return s.repo.Create(user)
}

func (s *UserServiceImpl) UpdateUser(user *models.User, id int) error{
	return s.repo.Update(user, id)
}

func (s *UserServiceImpl) DeleteUser(id int) error{
	return s.repo.Delete(id)
}