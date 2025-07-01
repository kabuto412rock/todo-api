package http

import (
	"context"
	"todo-app/internal/auth/usecase"

	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

type handler struct {
	RegisterUC usecase.RegisterUsecase
	LoginUC    usecase.LoginUsecase
	JWTSecret  []byte
}

func NewHandler(api huma.API, reg usecase.RegisterUsecase, login usecase.LoginUsecase, jwtSecret []byte) {
	h := handler{RegisterUC: reg, LoginUC: login, JWTSecret: jwtSecret}

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

	token, err := generateJWT(user.Username, h.JWTSecret)
	if err != nil {
		return nil, err
	}
	result := &loginOutput{}
	result.Body.Token = token
	return result, nil
}

// generateJWT generates a JWT token for the given username and secret.
func generateJWT(username string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"username": username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
