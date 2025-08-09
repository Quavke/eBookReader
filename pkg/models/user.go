package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
}

type UserDB struct {
	gorm.Model
	Username string     `gorm:"type:varchar(64);not null;uniqueIndex:ux_users_username"`
	PasswordHash []byte `json:"-" gorm:"not null"`
}

type Claims struct {
	UserID uint `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}