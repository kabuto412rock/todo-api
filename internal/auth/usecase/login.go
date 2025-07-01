package usecase

import (
	"errors"
	"todo-app/internal/auth/domain"

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
	return user, nil
}
