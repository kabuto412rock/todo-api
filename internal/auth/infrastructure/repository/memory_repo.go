package repository

import (
	"errors"
	"todo-app/internal/auth/domain"
)

var users = map[string]domain.AuthUser{}

type memoryRepo struct{}

func NewMemoryRepo() domain.AuthRepository {
	return &memoryRepo{}
}

func (r *memoryRepo) CreateUser(user domain.AuthUser) error {
	if _, ok := users[user.Username]; ok {
		return errors.New("user exists")
	}
	users[user.Username] = user
	return nil
}
func (r *memoryRepo) GetUserByUsername(username string) (domain.AuthUser, error) {
	u, ok := users[username]
	if !ok {
		return domain.AuthUser{}, errors.New("user not found")
	}
	return u, nil
}
