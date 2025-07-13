package http

import (
	"context"
	"todo-app/internal/auth/usecase"

	"github.com/danielgtaylor/huma/v2"
)

type handler struct {
	RegisterUC usecase.RegisterUsecase
	LoginUC    usecase.LoginUsecase
}

func NewHandler(api huma.API, reg usecase.RegisterUsecase, login usecase.LoginUsecase) {
	h := handler{RegisterUC: reg, LoginUC: login}

	huma.Post(api, "/auth/register", h.Register)
	huma.Post(api, "/auth/login", h.Login)
}

type authInput struct {
	Body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
}
type loginOutput struct {
	Body struct {
		Token string `json:"token"`
	}
}
type registerOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func (h *handler) Register(ctx context.Context, in *authInput) (*registerOutput, error) {
	userName, password := in.Body.Username, in.Body.Password
	if err := h.RegisterUC.Register(userName, password); err != nil {
		return nil, err
	}
	result := &registerOutput{}
	result.Body.Message = "User registered successfully"
	return result, nil
}
func (h *handler) Login(ctx context.Context, in *authInput) (*loginOutput, error) {
	userName, password := in.Body.Username, in.Body.Password
	user, err := h.LoginUC.Login(userName, password)
	if err != nil {
		return nil, err
	}
	result := &loginOutput{}
	result.Body.Token = user.Token
	return result, nil
}
