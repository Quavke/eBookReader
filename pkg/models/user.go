package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID uint64 `json:"id"`
	Username string `json:"username"`
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