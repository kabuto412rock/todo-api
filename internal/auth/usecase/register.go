package usecase

import (
	"todo-app/internal/auth/domain"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUsecase interface {
	Register(username, password string) error
}

type registerUsecase struct {
	repo domain.AuthRepository
}

func NewRegisterUsecase(repo domain.AuthRepository) RegisterUsecase {
	return &registerUsecase{repo}
}

func (uc *registerUsecase) Register(username, password string) error {
	_, err := uc.repo.GetUserByUsername(username)
	if err == nil {
		return ErrUserExists
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return uc.repo.CreateUser(domain.AuthUser{Username: username, PasswordHash: string(hash)})
}
