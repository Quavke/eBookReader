package services

import (
	"ebookr/pkg/models"
	"ebookr/pkg/repositories"
)

type AuthorService interface {
	GetAllAuthors(limit, page int, sort string)  (*models.Pagination, error)
	GetAuthorByID(id uint) 									     (*models.AuthorResp, error)
	CreateAuthor(author *models.Author)          error
	UpdateAuthor(author *models.UpdateAuthorReq, id uint)  error
	DeleteAuthor(id uint)                         error
}

type AuthorServiceImpl struct {
	repo repositories.AuthorRepo
}

func NewAuthorService(repo repositories.AuthorRepo) *AuthorServiceImpl{
	return &AuthorServiceImpl{repo: repo}
}

var _ AuthorService = (*AuthorServiceImpl)(nil)

func (s *AuthorServiceImpl) GetAllAuthors(limit, page int, sort string) (*models.Pagination, error){
	p := &models.Pagination{
		Limit: limit,
		Page: page,
		Sort: sort,
	}
	p, err := s.repo.GetAll(p)
	if err != nil {
		return nil, err
	}
	rows := p.Rows.([]models.Author)
	authors := make([]models.AuthorResp, 0, len(rows))
	for _, a := range rows {
		authors = append(authors, models.AuthorResp{
			UserID: a.UserID,
			Firstname: a.Firstname,
			Lastname: a.Lastname,
			Birthday: a.Birthday,
		})
	}
	p.Rows = authors
	return p, nil
}

func (s *AuthorServiceImpl) GetAuthorByID(id uint) (*models.AuthorResp, error){
	authorBD, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	author := &models.AuthorResp{
		UserID: authorBD.UserID,
		Firstname: authorBD.Firstname,
		Lastname: authorBD.Lastname,
		Birthday: authorBD.Birthday,
	}
	return author, nil
}

func (s *AuthorServiceImpl) CreateAuthor(author *models.Author) error{
	return s.repo.Create(author)
}

func (s *AuthorServiceImpl) UpdateAuthor(author *models.UpdateAuthorReq, id uint) error{
	if author.Firstname == "" && author.Lastname == "" && author.Birthday.IsZero() {
		return nil
	}

	// existing, err := s.repo.GetByID(id)
  //   if err != nil {
  //       return err
  //   }
    
  //   hasChanges := false
  //   if author.Firstname != "" && author.Firstname != existing.Firstname {
  //       hasChanges = true
  //   }
  //   if author.Lastname != "" && author.Lastname != existing.Lastname {
  //       hasChanges = true
  //   }
  //   if !author.Birthday.IsZero() && !author.Birthday.Equal(existing.Birthday.Time) {
  //       hasChanges = true
  //   }
    
  //   if !hasChanges {
  //       return nil // ничего не обновляем
  //   }
	
	return s.repo.Update(author, id)
}

func (s *AuthorServiceImpl) DeleteAuthor(id uint) error{
	return s.repo.Delete(id)
}