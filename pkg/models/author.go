package models

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	gorm.Model
	Firstname string    `json:"Firstname" binding:"required"`
	Lastname  string    `json:"Lastname"  binding:"required"`
	Birthday  time.Time `json:"Birthday"  binding:"required"`
	Books     []Book    `gorm:"foreignKey:AuthorID"`
}