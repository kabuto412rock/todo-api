package middleware_test

import (
	"context"
	"net/http"
	"testing"
	"todo-app/internal/api/middleware"
	"todo-app/internal/auth/domain"
	"todo-app/internal/auth/infrastructure/repository"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/require"
)

func TestNewAuthMiddleware(t *testing.T) {
	config := huma.DefaultConfig("Todo API", "1.0.0")
	config.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"myAuth": {
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}
	_, api := humatest.New(t, config)
	secret := []byte("your-256")
	authMiddleware := middleware.NewAuthMiddleware(api, secret)
	api.UseMiddleware(authMiddleware)

	grp := huma.NewGroup(api, "/todos")
	myAuthSecurity := []map[string][]string{
		{"myAuth": {}},
	}

	type demoResp struct {
		Message string `json:"message"`
	}
	huma.Register(grp, huma.Operation{
		OperationID: "create-todo",
		Summary:     "Create a new todo item",
		Method:      http.MethodPost,
		Path:        "",
		Security:    myAuthSecurity,
	}, func(ctx context.Context, i *struct{ Name string }) (*demoResp, error) {
		return &demoResp{Message: "Todo created"}, nil
	})

	// Test without Authorization header
	resp := api.Post("/todos", map[string]any{
		"name": "World",
	})
	require.Equal(t, 401, resp.Code)

	// Test with invalid token
	resp = api.Post("/todos", "Authorization: Bearer invalidtoken", map[string]any{
		"name": "World",
	})
	require.Equal(t, 401, resp.Code)

	// Test with valid token
	jwtTokenGenerator := repository.JWTTokenGenerator{Secret: string(secret)}
	token, err := jwtTokenGenerator.Generate(domain.AuthUser{Username: "testuser"})
	require.NoError(t, err)
	resp = api.Post("/todos", "Authorization: Bearer "+token, map[string]any{
		"name": "World",
	})
	require.Equal(t, 204, resp.Code)
}
