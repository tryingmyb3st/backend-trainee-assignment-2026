package auth_service

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/utils/jwt_utils"
	"fmt"

	"github.com/go-playground/validator/v10"
)

const (
	TEST_USER_UUID  = "8794e589-0ddb-43ce-9f92-16faafcf4ee4"
	TEST_ADMIN_UUID = "249be7cf-d419-4c54-97f2-d04107806e36"
)

func (s *AuthService) GetTestJWTByRole(user domain.User) (*string, error) {
	if err := s.ValidateDummy(user); err != nil {
		return nil, fmt.Errorf("invalid user role: %w", domain.INVALID_REQUEST)
	}

	var jwt *string
	var err error

	switch user.Role {
	case "admin":
		jwt, err = jwt_utils.GenerateJWT(TEST_ADMIN_UUID, user.Role)
	case "user":
		jwt, err = jwt_utils.GenerateJWT(TEST_USER_UUID, user.Role)
	}

	if err != nil {
		return nil, fmt.Errorf("generating jwt: %w", err)
	}

	return jwt, nil
}

func (s *AuthService) ValidateDummy(user domain.User) error {
	userValidator := validator.New()

	if err := userValidator.StructPartial(user, "Role"); err != nil {
		return err
	}

	return nil
}
