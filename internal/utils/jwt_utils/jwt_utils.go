package jwt_utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserId string
	Role   string
}

func GenerateJWT(userId string, role string) (*string, error) {

	expirationTime := time.Now().Add(12 * time.Hour)

	claims := &CustomClaims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil, fmt.Errorf("signing with key: %w", err)
	}

	return &tokenString, nil
}

func VerifyJWTtoken(tokenString string) (*CustomClaims, error) {
	claims := CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing claims: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("token invalid: %w", err)
	}

	return token.Claims.(*CustomClaims), nil
}
