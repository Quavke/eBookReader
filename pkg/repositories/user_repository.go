package repositories

import (
	"ebookr/pkg/models"
	"errors"

	"gorm.io/gorm"
)

type UserRepo interface {
    Create(user *models.UserDB) error
    GetByID(id uint) (*models.UserDB, error)
		IsExists(id uint) error
		IsAuthor(id uint) (bool, error)
    IsAuthors(ids []uint) (map[uint]bool, error)
    GetAll() ([]models.UserDB, error)
    Update(user *models.UpdateReq, id uint) error
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
	result := r.db.Where("id = ?", id).First(&user)
	if result.RowsAffected == 0 {
    return nil, gorm.ErrRecordNotFound
  }
	if err := result.Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepo) IsExists(id uint) (error) {
	var user models.UserDB
	result := r.db.Where("id = ?", id).First(&user)
	if result.RowsAffected == 0 {
    return gorm.ErrRecordNotFound
  }
	if err := result.Error; err != nil {
		return err
	}
	return nil
}

func (r *GormUserRepo) IsAuthor(id uint) (bool, error) {
	var author models.Author
	result := r.db.Where("user_id = ?", id).First(&author)
	if result.RowsAffected == 0 {
    return false, gorm.ErrRecordNotFound
  }
	if err := result.Error; err != nil {
		return false, err
	}
	return true, nil
}


func (r *GormUserRepo) IsAuthors(ids []uint) (map[uint]bool, error) {
    authorsSet := make(map[uint]bool)
    if len(ids) == 0 {
        return nil, errors.New("there are no user ids")
    }

    var authorUserIDs []uint
    result := r.db.Model(&models.Author{}).Where("user_id IN ?", ids).Pluck("user_id", &authorUserIDs)
    if err := result.Error; err != nil {
        return nil, err
    }
    for _, id := range authorUserIDs {
        authorsSet[id] = true
    }
    return authorsSet, nil
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

func (r *GormUserRepo) Update(user *models.UpdateReq, id uint) error{
	return r.db.Transaction(func(tx *gorm.DB) error {
        var existing models.UserDB
        result := tx.Where("id = ?", id).First(&existing)
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
    user.ID = id
    
    result := r.db.Select("Author", "Author.Books").Delete(&user)
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