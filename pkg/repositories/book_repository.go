package repositories

import (
	"ebookr/pkg/models"

	"gorm.io/gorm"
)

type BookRepo interface {
    Create(book *models.Book) error
    GetByID(id int) (*models.Book, error)
    GetAll() ([]models.Book, error)
    Update(book *models.Book, id int) error
    Delete(id int) error
}

type GormBookRepo struct {
    db *gorm.DB
}

var _ BookRepo = (*GormBookRepo)(nil)

func NewGormBookRepo(db *gorm.DB) *GormBookRepo{
	return &GormBookRepo{db: db}
}

func (r *GormBookRepo) Create(book *models.Book) error{
    var author models.Author
    err := r.db.Where("user_id = ?", book.AuthorID).First(&author).Error
    if err != nil {
        return err
    }
    result := r.db.Create(book)
    if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
	return result.Error
}

func (r *GormBookRepo) GetByID(id int) (*models.Book, error) {
	var book models.Book
    result := r.db.Where("id = ?", id).First(&book)
    if result.RowsAffected == 0 {
        return nil, gorm.ErrRecordNotFound
    }
	if err := result.Error; err != nil{
		return nil, err
	}
	return &book, nil
}



func (r *GormBookRepo) GetAll() ([]models.Book, error){
	var book []models.Book
    result := r.db.Find(&book)
    if result.RowsAffected == 0 {
        return nil, gorm.ErrRecordNotFound
    }
	if err := result.Error; err != nil{
		return nil, err
	}
	return book, nil
}

func (r *GormBookRepo) Update(book *models.Book, id int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
        var existing models.Book
        result := tx.Where("id = ?", id).First(&existing)
        if result.RowsAffected == 0 {
            return gorm.ErrRecordNotFound
        }
        if err := result.Error; err != nil {
            return err
        }
        updates := models.Book{
            Title:   book.Title,
            Content: book.Content,
        }

        result = tx.Model(&existing).Updates(updates)
        if result.RowsAffected == 0 {
            return gorm.ErrRecordNotFound
        }
        return result.Error
    })
}

func (r *GormBookRepo) Delete(id int) error{
	var book models.Book
	result := r.db.Where("id = ?", id).Delete(&book)
	if result.RowsAffected == 0 {
        return gorm.ErrRecordNotFound
    }
	return result.Error
}