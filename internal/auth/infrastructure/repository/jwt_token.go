package repository

import (
	"time"
	"todo-app/internal/auth/domain"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenGenerator struct {
	Secret string
}

func (j *JWTTokenGenerator) Generate(user domain.AuthUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Username,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte(j.Secret))
}
