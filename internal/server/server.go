package server

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	"todo-app/internal/api/middleware"
	authDomain "todo-app/internal/auth/domain"
	authRepo "todo-app/internal/auth/infrastructure/repository"
	authHttp "todo-app/internal/auth/interface/http"
	authUsecase "todo-app/internal/auth/usecase"
	todoDomain "todo-app/internal/todo/domain"
	todoHttp "todo-app/internal/todo/interface/http"
	todoUsecase "todo-app/internal/todo/usecase"
)

type Deps struct {
	JWTSecret string
	AuthRepo  authDomain.AuthRepository
	TokenGen  *authRepo.JWTTokenGenerator
	TodoRepo  todoDomain.TodoRepository
}

// NewRouter creates http.Handler with routes registered.
func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()
	cfg := huma.DefaultConfig("Todo API", "1.0.0")
	cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"myAuth": {Type: "http", Scheme: "bearer", BearerFormat: "JWT"},
	}
	api := humachi.New(r, cfg)
	Register(api, d)
	return r
}

// Register wires middleware & handlers onto an existing huma.API (for tests or custom adapters).
func Register(api huma.API, d Deps) {
	api.UseMiddleware(middleware.NewAuthMiddleware(api, []byte(d.JWTSecret)))
	todoUC := todoUsecase.NewTodoUseCase(d.TodoRepo)
	loginUC := authUsecase.NewLoginUsecase(d.AuthRepo, d.TokenGen)
	registerUC := authUsecase.NewRegisterUsecase(d.AuthRepo)
	todoHttp.NewTodoHandler(api, todoUC)
	authHttp.NewHandler(api, registerUC, loginUC)
}
