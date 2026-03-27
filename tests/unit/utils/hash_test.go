package utils_test

import (
	hash_utils "backend-assignment-avito/internal/utils/hash"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := hash_utils.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	assert.NoError(t, err)
}

func TestHashPasswordEmpty(t *testing.T) {
	password := ""

	hash, err := hash_utils.HashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	assert.NoError(t, err)
}

func TestCheckPasswordHash_CorrectPassword(t *testing.T) {
	password := "correctPassword123"

	hash, err := hash_utils.HashPassword(password)
	require.NoError(t, err)

	result := hash_utils.CheckPasswordHash(password, hash)

	assert.True(t, result)
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	correctPassword := "correctPassword123"
	wrongPassword := "wrongPassword456"

	hash, err := hash_utils.HashPassword(correctPassword)
	require.NoError(t, err)

	result := hash_utils.CheckPasswordHash(wrongPassword, hash)

	assert.False(t, result)
}
