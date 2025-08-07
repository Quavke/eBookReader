package repositories

import (
	"ebookr/pkg/models"

	"gorm.io/gorm"
)

type AuthorRepo interface {
    Create(author *models.Author) error
    GetByID(id int) (*models.Author, error)
    GetAll() ([]models.Author, error)
    Update(author *models.Author, id int) error
    Delete(id int) error
}

type GormAuthorRepo struct {
    db *gorm.DB
}

var _ AuthorRepo = (*GormAuthorRepo)(nil)

func NewGormAuthorRepo(db *gorm.DB) *GormAuthorRepo{
	return &GormAuthorRepo{db: db}
}

func (r GormAuthorRepo) Create(author *models.Author) error{
	return r.db.Create(&author).Error
}

func (r GormAuthorRepo) GetByID(id int) (*models.Author, error){
	var author *models.Author
	if err := r.db.First(author, id).Error; err != nil{
		return nil, err
	}
	return author, nil
}

func (r GormAuthorRepo) GetAll() ([]models.Author, error){
	var authors []models.Author
	if err := r.db.Find(&authors).Error; err != nil{
		return nil, err
	}
	return authors, nil
}

func (r GormAuthorRepo) Update(author *models.Author, id int) error{
	return r.db.Transaction(func(tx *gorm.DB) error {
        var existing models.Author
        if err := tx.First(&existing, id).Error; err != nil {
            return err
        }
        
        updates := models.Author{
					Firstname: author.Firstname,
					Lastname: author.Lastname,
					Birthday: author.Birthday,
        }
        
        return tx.Model(&existing).Updates(updates).Error
    })
}

func (r GormAuthorRepo) Delete(id int) error{
	var author *models.Author
	result := r.db.Delete(author, id)
	if result.RowsAffected == 0 {
    return gorm.ErrRecordNotFound
  }
	return result.Error
}