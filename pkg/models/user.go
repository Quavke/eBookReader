package models

import (
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required" gorm:"min=5, not null"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	UserID uint64 `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}