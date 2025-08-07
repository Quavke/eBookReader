package repositories

import (
	"ebookr/pkg/models"

	"gorm.io/gorm"
)

type UserRepo interface {
    Create(user *models.User) error
    GetByID(id int) (*models.User, error)
    GetAll() ([]models.User, error)
    Update(user *models.User, id int) error
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
	result := r.db.Create(user)
	if result.RowsAffected == 0 {
    return gorm.ErrRecordNotFound
  }
	return result.Error
}

func (r *GormUserRepo) GetByID(id int) (*models.User, error){
	var user *models.User
	result := r.db.First(user, id)
	if result.RowsAffected == 0 {
    return nil, gorm.ErrRecordNotFound
  }
	if err := result.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormUserRepo) GetAll() ([]models.User, error){
	var user []models.User
	result := r.db.Find(&user)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err := result.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormUserRepo) Update(user *models.User, id int) error{
	return r.db.Transaction(func(tx *gorm.DB) error {
        var existing models.User
        result := tx.First(&existing, id)
        if result.RowsAffected == 0 {
            return gorm.ErrRecordNotFound
        }
        if err := result.Error; err != nil {
            return err
        }
        updates := models.User{
            Username: user.Username,
        }

        result = tx.Model(&existing).Updates(updates)
        if result.RowsAffected == 0 {
            return gorm.ErrRecordNotFound
        }
        return result.Error
    })
}

func (r *GormUserRepo) Delete(id int) error{
	var user models.User
	result := r.db.Delete(&user, id)
	if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
	return result.Error
}