package models

import (
	"time"
)

type Author struct {
	Firstname string `json:"Firstname"`
	Lastname  string `json:"Lastname"`
	Birthday  time.Time `json:"Birthday"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}