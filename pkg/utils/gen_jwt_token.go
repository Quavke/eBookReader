package utils

import (
	"ebookr/pkg/models"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user *models.User, jwtSecretKey []byte) (string, error) {
	expTime := time.Now().Add(24 * time.Hour)

	claims := &models.Claims{
		UserID: user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Subject: fmt.Sprintf("%d", user.ID),
			Issuer: "eBookReader",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}