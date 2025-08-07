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