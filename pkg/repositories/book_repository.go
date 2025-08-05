package repositories

import (
	"ebookr/pkg/models"

	"gorm.io/gorm"
)

type BookRepo interface {
    Create(book *models.Book) error
    GetByID(id int) (*models.Book, error)
    GetAll() ([]models.Book, error)
    Update(bookNew *models.Book, id int) error
    Delete(id int) error
}

type GormBookRepo struct {
    db *gorm.DB
}

var _ BookRepo = (*GormBookRepo)(nil)

func NewGormBookRepo(db *gorm.DB) *GormBookRepo{
	return &GormBookRepo{db: db}
}

func (r GormBookRepo) Create(book *models.Book) error{
	return r.db.Create(book).Error
}

func (r GormBookRepo) GetByID(id int) (*models.Book, error) {
	var book models.Book
	if err := r.db.First(&book, id).Error; err != nil{
		return nil, err
	}
	return &book, nil
}



func (r GormBookRepo) GetAll() ([]models.Book, error){
	var book []models.Book
	if err := r.db.Find(&book).Error; err != nil{
		return nil, err
	}
	return book, nil
}

func (r GormBookRepo) Update(bookNew *models.Book, id int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
        var existing models.Book
        if err := tx.First(&existing, id).Error; err != nil {
            return err
        }
        
        // Обновляем только необходимые поля
        updates := models.Book{
            Title:   bookNew.Title,
            Content: bookNew.Content,
            Author: models.Author{
                Firstname: bookNew.Author.Firstname,
                Lastname:  bookNew.Author.Lastname,
                Birthday:  bookNew.Author.Birthday,
            },
        }
        
        return tx.Model(&existing).Updates(updates).Error
    })
}

func (r GormBookRepo) Delete(id int) error{
	var book models.Book
	result := r.db.Delete(&book, id)
	if result.RowsAffected == 0 {
    return gorm.ErrRecordNotFound
  }
	return result.Error
}