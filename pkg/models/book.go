package models

import (
	"time"
)

type Book struct {
	ID        int    `gorm:"primaryKey;not null"`
	Title     string `json:"title" gorm:"not null;unique" binding:"required,min=3"`
	Content   string `json:"content" gorm:"not null;unique" binding:"required,min=50"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Author    Author `json:"author" gorm:"embedded;embeddedPrefix:author_" binding:"required"`
}