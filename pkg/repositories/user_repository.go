package repositories

import (
	"ebookr/pkg/models"

	"gorm.io/gorm"
)

type UserRepo interface {
    Create(user *models.User) error
    GetByID(id int) (*models.User, error)
    GetAll() ([]models.User, error)
    Update(userNew *models.User, id int) error
    Delete(id int) error
}

type GormUserRepo struct {
    db *gorm.DB
}


var _ UserRepo = (*GormUserRepo)(nil)

func NewGormUserRepo(db *gorm.DB) *GormUserRepo{
	return &GormUserRepo{db: db}
}

func (r *GormUserRepo) Create(user *models.User) error{
	return nil
}

func (r *GormUserRepo) GetByID(id int) (*models.User, error){
	return &models.User{}, nil
}

func (r *GormUserRepo) GetAll() ([]models.User, error){
	return []models.User{}, nil
}

func (r *GormUserRepo) Update(userNew *models.User, id int) error{
	return nil
}

func (r *GormUserRepo) Delete(id int) error{
	return nil
}