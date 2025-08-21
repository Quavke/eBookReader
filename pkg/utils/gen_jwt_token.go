package utils

import (
	"github.com/Quavke/eBookReader/pkg/models"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(claims *models.Claims, jwtSecretKey []byte) (string, error) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}