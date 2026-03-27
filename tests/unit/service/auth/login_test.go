package auth_service

import (
	"backend-assignment-avito/internal/core/domain"
	auth_service "backend-assignment-avito/internal/features/auth/service"
	hash_utils "backend-assignment-avito/internal/utils/hash"
	"backend-assignment-avito/internal/utils/jwt_utils"
	auth_service_mock "backend-assignment-avito/mocks/auth_service"
	"context"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginUser(t *testing.T) {
	m := auth_service_mock.NewMockAuthRepository(t)
	authService := auth_service.NewAuthService(m)
	ctx := context.Background()

	email := "user@example.com"
	password := "correctPassword123"

	hashedPassword, err := hash_utils.HashPassword(password)
	require.NoError(t, err)

	expectedUser := &domain.User{
		ID:       "user-123",
		Email:    email,
		Password: hashedPassword,
		Role:     "user",
	}

	m.On("GetUser", ctx, email).Return(expectedUser, nil)

	jwtToken, err := authService.LoginUser(ctx, email, password)

	require.NoError(t, err)
	assert.NotNil(t, jwtToken)
	assert.NotEmpty(t, *jwtToken)

	claims := jwt_utils.CustomClaims{}
	token, err := jwt.ParseWithClaims(*jwtToken, &claims, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	require.NoError(t, err)
	require.True(t, token.Valid)
	assert.Equal(t, expectedUser.ID, claims.UserId)
	assert.Equal(t, expectedUser.Role, claims.Role)
}

func TestLoginUserFail(t *testing.T) {
	m := auth_service_mock.NewMockAuthRepository(t)
	authService := auth_service.NewAuthService(m)
	ctx := context.Background()

	email := "nonexistent@example.com"
	password := "anyPassword"

	m.On("GetUser", ctx, email).Return(nil, pgx.ErrNoRows)

	jwtToken, err := authService.LoginUser(ctx, email, password)

	assert.Nil(t, jwtToken)
	require.Error(t, err)
	assert.ErrorIs(t, err, domain.UNAUTHORIZED)
	assert.Contains(t, err.Error(), "no user")
}
