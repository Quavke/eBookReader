package models

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title     string `json:"title" gorm:"not null;unique" binding:"required,min=1,max=400"`
	Content   string `json:"content" gorm:"not null;unique" binding:"required,min=10"`
	AuthorID  uint   `json:"-" gorm:"not null;constraint:OnUpdate:CASCADE;"`
	Author    *Author `json:"-" gorm:"foreignKey:AuthorID;references:UserID"`
}

type BookResp struct {
  ID       uint   `json:"id"`
  Title    string `json:"title"`
  Content  string `json:"content"`
  AuthorID uint   `json:"author_id"`
}