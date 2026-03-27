package auth_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *AuthRepository) GetUser(ctx context.Context, email string) (*domain.User, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, email, password, role, created_at
	FROM users
	WHERE email=$1
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, email)

	var model UserModel
	err := row.Scan(
		&model.ID,
		&model.Email,
		&model.Password,
		&model.Role,
		&model.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan returning user: %w", err)
	}

	return &domain.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
	}, nil
}
