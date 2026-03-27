package auth_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
)

type AuthService struct {
	authRepo AuthRepository
}

type AuthRepository interface {
	SaveNewUser(ctx context.Context, user domain.User) (*domain.User, error)
	GetUser(ctx context.Context, email string) (*domain.User, error)
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{
		authRepo: repo,
	}
}
