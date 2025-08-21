package repositories

import (
	"fmt"

	"github.com/Quavke/eBookReader/pkg/models"

	"gorm.io/gorm"
)

type BookRepo interface {
    Create(book *models.Book) error
    GetByID(id uint) (*models.Book, error)
    GetAll(p *models.Pagination) (*models.Pagination, error)
    IsBelongsTo(id uint, authorID uint) (bool, error)
    Update(book *models.Book, id uint) error
    Delete(id uint) error
}

type GormBookRepo struct {
    db *gorm.DB
}

var _ BookRepo = (*GormBookRepo)(nil)

func NewGormBookRepo(db *gorm.DB) *GormBookRepo{
	return &GormBookRepo{db: db}
}

func (r *GormBookRepo) Create(book *models.Book) error{
    result := r.db.Create(book)
    if result.RowsAffected == 0 {
        return result.Error
    }
	return result.Error
}

func (r *GormBookRepo) GetByID(id uint) (*models.Book, error) {
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

func (r *GormBookRepo) IsBelongsTo(id uint, authorID uint) (bool, error){
    var book models.Book
    result := r.db.Where("id = ? AND author_id = ?", id, authorID).First(&book)
    if result.RowsAffected == 0 {
        return false, result.Error
    }
    if err := result.Error; err != nil {
        return false, err
    }
    return true, nil
}

func (r *GormBookRepo) GetAll(p *models.Pagination) (*models.Pagination, error){
	var books []models.Book
    result := r.db.Scopes(models.Paginate(books, p, r.db)).Find(&books)

    p.Rows = books

    if len(p.Rows.([]models.Book)) == 0 {
        return nil, gorm.ErrRecordNotFound
    }
	if err := result.Error; err != nil{
		return nil, err
	}
	return p, nil
}

func (r *GormBookRepo) Update(book *models.Book, id uint) error {
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
            return fmt.Errorf("no book found with id %d. Error: %v", id, result.Error)
        }
        return result.Error
    })
}

func (r *GormBookRepo) Delete(id uint) error{
	var book models.Book
	result := r.db.Where("id = ?", id).Delete(&book)
	if result.RowsAffected == 0 {
        return fmt.Errorf("no book found with id %d. Error: %v", id, result.Error)
    }
	return result.Error
}