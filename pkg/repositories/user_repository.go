package repositories

import (
	"ebookr/pkg/models"

	"gorm.io/gorm"
)

type UserRepo interface {
    Create(user *models.UserDB) error
    GetByID(id uint) (*models.UserDB, error)
    GetAll() ([]models.UserDB, error)
    Update(user *models.UserDB, id uint) error
    Delete(id uint) error
    GetByUsername(username string) (*models.UserDB, error)
}

type GormUserRepo struct {
    db *gorm.DB
}


var _ UserRepo = (*GormUserRepo)(nil)

func NewGormUserRepo(db *gorm.DB) *GormUserRepo{
	return &GormUserRepo{db: db}
}

func (r *GormUserRepo) Create(user *models.UserDB) error{
	result := r.db.Create(user)
	if result.RowsAffected == 0 {
    return gorm.ErrRecordNotFound
  }
	return result.Error
}

func (r *GormUserRepo) GetByID(id uint) (*models.UserDB, error){
	var user models.UserDB
	result := r.db.First(&user, id)
	if result.RowsAffected == 0 {
    return nil, gorm.ErrRecordNotFound
  }
	if err := result.Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepo) GetAll() ([]models.UserDB, error){
	var user []models.UserDB
	result := r.db.Find(&user)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if err := result.Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormUserRepo) Update(user *models.UserDB, id uint) error{
	return r.db.Transaction(func(tx *gorm.DB) error {
        var existing models.UserDB
        result := tx.First(&existing, id)
        if result.RowsAffected == 0 {
            return gorm.ErrRecordNotFound
        }
        if err := result.Error; err != nil {
            return err
        }
        updates := models.UserDB{
            Username: user.Username,
        }

        result = tx.Model(&existing).Updates(updates)
        if result.RowsAffected == 0 {
            return gorm.ErrRecordNotFound
        }
        return result.Error
    })
}

func (r *GormUserRepo) Delete(id uint) error{
	var user models.UserDB
	result := r.db.Delete(&user, id)
	if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
	return result.Error
}

func (r *GormUserRepo) GetByUsername(username string) (*models.UserDB, error) {
    var user models.UserDB
	result := r.db.Where("username = ?", username).First(&user)
	if result.RowsAffected == 0 {
    return nil, gorm.ErrRecordNotFound
    }
	if err := result.Error; err != nil {
		return nil, err
	}
	return &user, nil
}