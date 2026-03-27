package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var AllowedRoles = []string{"admin", "user"}

type User struct {
	ID        string `validate:"required,uuid"`
	Email     string `validate:"email"`
	Password  string
	Role      string `validate:"oneof=admin user"`
	CreatedAt time.Time
}

func (u *User) Validate() error {

	validator := validator.New()

	if err := validator.Struct(u); err != nil {
		return fmt.Errorf("validator struct: %w", err)
	}

	return nil
}
