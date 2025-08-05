package models

import (
	"time"
)

type Author struct {
	Firstname string    `json:"Firstname"  binding:"required"`
	Lastname  string    `json:"Lastname"  binding:"required"`
	Birthday  time.Time `json:"Birthday"  binding:"required"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}