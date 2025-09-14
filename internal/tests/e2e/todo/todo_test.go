package todo_test

import (
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"

	authDomain "todo-app/internal/auth/domain"
	authRepo "todo-app/internal/auth/infrastructure/repository"
	"todo-app/internal/server"
	todoRepo "todo-app/internal/todo/infrastructure/repository"
)

// Basic smoke test ensuring routes register & open list endpoint (adjust path as needed)
func TestTodoAPI_ListEmpty(t *testing.T) {
	config := huma.DefaultConfig("Todo API", "1.0.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"myAuth": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
	}
	_, api := humatest.New(t, config)

	deps := server.Deps{
		JWTSecret: "test-secret",
		AuthRepo:  authRepo.NewMemoryRepo(),
		TokenGen:  &authRepo.JWTTokenGenerator{Secret: "test-secret"},
		TodoRepo:  todoRepo.NewMemoryTodoRepository(),
	}
	server.Register(api, deps)

	// create auth token for header
	token, err := deps.TokenGen.Generate(authDomain.AuthUser{Username: "tester"})
	if err != nil {
		t.Fatalf("token gen: %v", err)
	}
	resp := api.Get("/todos", "Authorization: Bearer "+token)
	if resp.Code != 200 {
		t.Fatalf("expected 200 got %d", resp.Code)
	}
}
