package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	authDomain "todo-app/internal/auth/domain"
	authRepo "todo-app/internal/auth/infrastructure/repository"
	"todo-app/internal/server"
	todoRepo "todo-app/internal/todo/infrastructure/repository"
)

func TestTodoList_HTTP(t *testing.T) {
	deps := server.Deps{
		JWTSecret: "test-secret",
		AuthRepo:  authRepo.NewMemoryRepo(),
		TokenGen:  &authRepo.JWTTokenGenerator{Secret: "test-secret"},
		TodoRepo:  todoRepo.NewMemoryTodoRepository(),
	}
	h := server.NewRouter(deps)
	ts := httptest.NewServer(h)
	defer ts.Close()
	// need auth token header
	token, err := deps.TokenGen.Generate(authDomain.AuthUser{Username: "tester"})
	if err != nil {
		t.Fatalf("token gen: %v", err)
	}
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/todos", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http get: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("expected 200 got %d", res.StatusCode)
	}
}
