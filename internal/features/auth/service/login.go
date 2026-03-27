package auth_service

import (
	"backend-assignment-avito/internal/core/domain"
	hash_utils "backend-assignment-avito/internal/utils/hash"
	"backend-assignment-avito/internal/utils/jwt_utils"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (*string, error) {
	user, err := s.authRepo.GetUser(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("no user: %w", domain.UNAUTHORIZED)
		}

		return nil, fmt.Errorf("get from database: %w", domain.INTERNAL_ERROR)
	}

	if !hash_utils.CheckPasswordHash(password, user.Password) {
		return nil, fmt.Errorf("wrong password: %w", domain.UNAUTHORIZED)
	}

	jwt, err := jwt_utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generating jwt: %w", domain.INTERNAL_ERROR)
	}

	return jwt, nil
}
