package usecase_test

import (
	"errors"
	"testing"

	"todo-app/internal/auth/domain"
	"todo-app/internal/auth/usecase"
)

// Reuse simple mocks (separate from login tests to keep files independent)
type regMockRepo struct {
	getUser func(username string) (domain.AuthUser, error)
	create  func(user domain.AuthUser) error
}

func (m *regMockRepo) GetUserByUsername(username string) (domain.AuthUser, error) {
	if m.getUser != nil {
		return m.getUser(username)
	}
	return domain.AuthUser{}, errors.New("not found")
}

func (m *regMockRepo) CreateUser(user domain.AuthUser) error {
	if m.create != nil {
		return m.create(user)
	}
	return nil
}

func TestRegister_Success(t *testing.T) {
	var created domain.AuthUser
	repo := &regMockRepo{
		getUser: func(username string) (domain.AuthUser, error) { // simulate not found
			return domain.AuthUser{}, errors.New("not found")
		},
		create: func(user domain.AuthUser) error {
			created = user
			if user.Username != "newuser" {
				return errors.New("unexpected username")
			}
			if user.PasswordHash == "plain" || user.PasswordHash == "" { // should be hashed
				return errors.New("password not hashed")
			}
			return nil
		},
	}
	uc := usecase.NewRegisterUsecase(repo)
	if err := uc.Register("newuser", "plaintext"); err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if created.Username == "" {
		t.Fatalf("expected user to be created")
	}
}

func TestRegister_UserExists(t *testing.T) {
	repo := &regMockRepo{
		getUser: func(username string) (domain.AuthUser, error) { // simulate existing
			return domain.AuthUser{Username: "taken", PasswordHash: "hash"}, nil
		},
	}
	uc := usecase.NewRegisterUsecase(repo)
	err := uc.Register("taken", "x")
	if err == nil || !errors.Is(err, usecase.ErrUserExists) {
		t.Fatalf("expected ErrUserExists, got %v", err)
	}
}

func TestRegister_CreateUserFailure(t *testing.T) {
	repo := &regMockRepo{
		getUser: func(username string) (domain.AuthUser, error) { // not found
			return domain.AuthUser{}, errors.New("not found")
		},
		create: func(user domain.AuthUser) error { return errors.New("insert failed") },
	}
	uc := usecase.NewRegisterUsecase(repo)
	err := uc.Register("another", "pwd")
	if err == nil || err.Error() != "insert failed" {
		t.Fatalf("expected insert failed error, got %v", err)
	}
}

// Note: Simulating bcrypt.GenerateFromPassword error directly is non-trivial without abstraction; skipped.
