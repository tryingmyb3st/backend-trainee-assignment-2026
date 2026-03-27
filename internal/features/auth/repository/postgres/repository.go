package auth_repository

import postgres_pool "backend-assignment-avito/internal/core/repository/postgres"

type AuthRepository struct {
	ConnPool postgres_pool.Pool
}

func NewAuthRepository(conn postgres_pool.Pool) *AuthRepository {
	return &AuthRepository{
		ConnPool: conn,
	}
}
