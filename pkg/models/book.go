package models

import (
	"time"
)

type Book struct {
	ID        int    `gorm:"primaryKey;not null"`
	Title     string `json:"title" gorm:"not null;unique"`
	Content   string `json:"content" gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Author    Author `json:"author" gorm:"embedded;embeddedPrefix:author_"`
}