package usecase

import (
	"errors"
	"time"
	"todo-app/internal/auth/domain"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginUsecase interface {
	Login(username, password string) (domain.AuthUser, error)
}

type loginUsecase struct {
	repo domain.AuthRepository
}

func NewLoginUsecase(repo domain.AuthRepository) LoginUsecase {
	return &loginUsecase{repo}
}

func (uc *loginUsecase) Login(username, password string) (domain.AuthUser, error) {
	user, err := uc.repo.GetUserByUsername(username)
	if err != nil {
		return domain.AuthUser{}, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return domain.AuthUser{}, errors.New("invalid credentials")
	}
	tokenString, err := uc.generateJWT(user)
	if err != nil {
		return domain.AuthUser{}, errors.New("failed to generate token")
	}
	user.Token = tokenString
	return user, nil
}

func (uc *loginUsecase) generateJWT(user domain.AuthUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.Username,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString([]byte("your-256-bit-secret"))
}
