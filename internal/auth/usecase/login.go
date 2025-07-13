package usecase

import (
	"errors"
	"todo-app/internal/auth/domain"

	"golang.org/x/crypto/bcrypt"
)

type LoginResult struct {
	User  domain.AuthUser
	Token string
}

type LoginUsecase interface {
	Login(username, password string) (LoginResult, error)
}

type loginUsecase struct {
	repo     domain.AuthRepository
	tokenGen domain.TokenGenerator
}

func NewLoginUsecase(repo domain.AuthRepository, tokenGen domain.TokenGenerator) LoginUsecase {
	return &loginUsecase{repo: repo, tokenGen: tokenGen}
}

func (uc *loginUsecase) Login(username, password string) (LoginResult, error) {
	user, err := uc.repo.GetUserByUsername(username)
	if err != nil {
		return LoginResult{}, errors.New("invalid credentials")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return LoginResult{}, errors.New("invalid credentials")
	}
	tokenString, err := uc.tokenGen.Generate(user)
	if err != nil {
		return LoginResult{}, errors.New("failed to generate token")
	}
	return LoginResult{User: user, Token: tokenString}, nil
}
