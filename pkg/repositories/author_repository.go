package repositories

import (
	"ebookr/pkg/models"

	"gorm.io/gorm"
)

type AuthorRepo interface {
    Create(author *models.Author) error
    GetByID(id uint) (*models.Author, error)
    GetAll() ([]models.Author, error)
    Update(author *models.Author, id uint) error
    Delete(id uint) error
}

type GormAuthorRepo struct {
    db *gorm.DB
}

var _ AuthorRepo = (*GormAuthorRepo)(nil)


func NewGormAuthorRepo(db *gorm.DB) *GormAuthorRepo{
	return &GormAuthorRepo{db: db}
}

func (r GormAuthorRepo) Create(author *models.Author) error{
	result := r.db.Create(&author)
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

func (r GormAuthorRepo) GetByID(id uint) (*models.Author, error){
	var author models.Author
	result := r.db.Preload("Books").Where("user_id = ?", id).First(&author)
	if result.RowsAffected == 0{
		return nil, gorm.ErrRecordNotFound
	}
	if err := result.Error; err != nil{
		return nil, err
	}
	return &author, nil
}

func (r GormAuthorRepo) GetAll() ([]models.Author, error){
	var authors []models.Author
	result := r.db.Preload("Books").Find(&authors)
	if result.RowsAffected == 0{
		return nil, gorm.ErrRecordNotFound
	}
	if err := result.Error; err != nil{
		return nil, err
	}
	return authors, nil
}

func (r GormAuthorRepo) Update(author *models.Author, id uint) error{
	return r.db.Transaction(func(tx *gorm.DB) error {
        var existing models.Author
				result := tx.Where("user_id = ?", id).First(&existing)
				if result.RowsAffected == 0{
					return gorm.ErrRecordNotFound
				}
        if err := result.Error; err != nil {
            return err
        }
        
        updates := models.Author{
					Firstname: author.Firstname,
					Lastname: author.Lastname,
					Birthday: author.Birthday,
        }
        result = tx.Model(&existing).Updates(updates)
				if result.RowsAffected == 0{
					return gorm.ErrRecordNotFound
				}
        return result.Error
    })
}

func (r GormAuthorRepo) Delete(id uint) error{
	author := models.Author{UserID: uint(id)}
	result := r.db.Select("Books").Delete(&author)
	if result.RowsAffected == 0 {
    return gorm.ErrRecordNotFound
  }
	return result.Error
}