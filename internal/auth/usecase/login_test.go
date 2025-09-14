package usecase_test

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"todo-app/internal/auth/domain"
	"todo-app/internal/auth/usecase"
)

// mockAuthRepo implements domain.AuthRepository
// Behavior is controlled via function fields so tests can override only what they need.
type mockAuthRepo struct {
	getUserFunc func(username string) (domain.AuthUser, error)
	createFunc  func(user domain.AuthUser) error
}

func (m *mockAuthRepo) CreateUser(user domain.AuthUser) error {
	if m.createFunc != nil {
		return m.createFunc(user)
	}
	return nil
}

func (m *mockAuthRepo) GetUserByUsername(username string) (domain.AuthUser, error) {
	if m.getUserFunc != nil {
		return m.getUserFunc(username)
	}
	return domain.AuthUser{}, errors.New("not found")
}

// mockTokenGen implements domain.TokenGenerator
type mockTokenGen struct {
	generateFunc func(user domain.AuthUser) (string, error)
}

func (m *mockTokenGen) Generate(user domain.AuthUser) (string, error) {
	if m.generateFunc != nil {
		return m.generateFunc(user)
	}
	return "", nil
}

func TestLogin_Success(t *testing.T) {
	password := "secret123"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user := domain.AuthUser{Username: "alice", PasswordHash: string(hash)}

	repo := &mockAuthRepo{getUserFunc: func(username string) (domain.AuthUser, error) {
		if username != user.Username {
			return domain.AuthUser{}, errors.New("not found")
		}
		return user, nil
	}}
	tokenValue := "token-abc"
	tokenGen := &mockTokenGen{generateFunc: func(u domain.AuthUser) (string, error) { return tokenValue, nil }}

	uc := usecase.NewLoginUsecase(repo, tokenGen)
	result, err := uc.Login(user.Username, password)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
	if result.User.Username != user.Username {
		t.Errorf("expected username %s, got %s", user.Username, result.User.Username)
	}
	if result.Token != tokenValue {
		t.Errorf("expected token %s, got %s", tokenValue, result.Token)
	}
}

func TestLogin_InvalidUsername(t *testing.T) {
	repo := &mockAuthRepo{getUserFunc: func(username string) (domain.AuthUser, error) {
		return domain.AuthUser{}, errors.New("not found")
	}}
	tokenGen := &mockTokenGen{}
	uc := usecase.NewLoginUsecase(repo, tokenGen)
	_, err := uc.Login("bob", "irrelevant")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	goodPassword := "correct"
	hash, _ := bcrypt.GenerateFromPassword([]byte(goodPassword), bcrypt.DefaultCost)
	storedUser := domain.AuthUser{Username: "carol", PasswordHash: string(hash)}
	repo := &mockAuthRepo{getUserFunc: func(username string) (domain.AuthUser, error) { return storedUser, nil }}
	tokenGen := &mockTokenGen{}

	uc := usecase.NewLoginUsecase(repo, tokenGen)
	_, err := uc.Login("carol", "wrong")
	if err == nil || err.Error() != "invalid credentials" {
		t.Fatalf("expected invalid credentials error, got %v", err)
	}
}

func TestLogin_TokenGenerationFailure(t *testing.T) {
	password := "pass"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	storedUser := domain.AuthUser{Username: "dave", PasswordHash: string(hash)}
	repo := &mockAuthRepo{getUserFunc: func(username string) (domain.AuthUser, error) { return storedUser, nil }}
	tokenGen := &mockTokenGen{generateFunc: func(u domain.AuthUser) (string, error) { return "", errors.New("boom") }}

	uc := usecase.NewLoginUsecase(repo, tokenGen)
	_, err := uc.Login("dave", password)
	if err == nil || err.Error() != "failed to generate token" {
		t.Fatalf("expected failed to generate token error, got %v", err)
	}
}
