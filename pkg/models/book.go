package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title     string `json:"title" gorm:"not null;unique" binding:"required,min=3"`
	Content   string `json:"content" gorm:"not null;unique" binding:"required,min=50"`
	AuthorID  uint64
}