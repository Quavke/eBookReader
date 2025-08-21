package models

import (
	"math"

	"gorm.io/gorm"
)

type Pagination struct {
	Limit      uint         `json:"limit,omitempty" query:"limit"`
	Page       uint         `json:"page,omitempty" query:"page"`
	Sort       string      `json:"sort,omitempty" query:"sort"`
	TotalRows  uint64       `json:"total_rows"`
	TotalPages uint         `json:"total_pages"`
	Rows       any         `json:"rows"`
}

func (p *Pagination) GetOffset() uint {
	return (p.GetPage() - 1) * p.GetLimit()
}

func (p *Pagination) GetLimit() uint {
	if p.Limit == 0 {
		p.Limit = 10
	}
	return p.Limit
}

func (p *Pagination) GetPage() uint {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "Id desc"
	}
	return p.Sort
}

func Paginate(value any, pagination *Pagination, db *gorm.DB) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	db.Model(value).Count(&totalRows)

	pagination.TotalRows = uint64(totalRows)
	totalPages := uint(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(int(pagination.GetOffset())).Limit(int(pagination.GetLimit())).Order(pagination.GetSort())
	}
}