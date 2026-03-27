package auth_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/transport/server"
	"context"
)

type AuthHTTPHandler struct {
	AuthService AuthService
}

type AuthService interface {
	GetTestJWTByRole(user domain.User) (*string, error)
	RegisterUser(ctx context.Context, user domain.User, password string) (*domain.User, error)
	LoginUser(ctx context.Context, email, password string) (*string, error)
}

func NewAuthHandler(authServ AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		AuthService: authServ,
	}
}

func (h *AuthHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  "POST",
			URL:     "/dummyLogin",
			Handler: h.GetDummyLogin,
		},
		{
			Method:  "POST",
			URL:     "/register",
			Handler: h.RegisterUser,
		},
		{
			Method:  "POST",
			URL:     "/login",
			Handler: h.LoginUser,
		},
	}
}
