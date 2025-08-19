package repositories

import (
	"ebookr/pkg/models"

	"gorm.io/gorm"
)

type AuthorRepo interface {
    Create(author *models.Author) error
    GetByID(id uint) (*models.Author, error)
    GetAll(p *models.Pagination) (*models.Pagination, error)
    Update(author *models.UpdateAuthorReq, id uint) error
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

func (r GormAuthorRepo) GetAll(p *models.Pagination) (*models.Pagination, error){
	var authors []models.Author
	result := r.db.Scopes(models.Paginate(authors, p, r.db)).Find(&authors)

	p.Rows = authors

	if len(p.Rows.([]models.Author)) == 0{
		return nil, gorm.ErrRecordNotFound
	}
	if err := result.Error; err != nil{
		return nil, err
	}
	return p, nil
}

func (r GormAuthorRepo) Update(author *models.UpdateAuthorReq, id uint) error{
	return r.db.Transaction(func(tx *gorm.DB) error {
        var existing models.Author
				result := tx.Where("user_id = ?", id).First(&existing)
				if result.RowsAffected == 0{
					return gorm.ErrRecordNotFound
				}
        if err := result.Error; err != nil {
            return err
        }

				
        updates := make(map[string]interface{})
        if author.Firstname != "" && author.Firstname != existing.Firstname {
            updates["firstname"] = author.Firstname
        }
        if author.Lastname != "" && author.Lastname != existing.Lastname {
            updates["lastname"] = author.Lastname
        }
        if !author.Birthday.IsZero() && !author.Birthday.Equal(existing.Birthday.Time) {
            updates["birthday"] = author.Birthday
        }
        
        if len(updates) > 0 {
            return tx.Model(&existing).Updates(updates).Error
        }
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