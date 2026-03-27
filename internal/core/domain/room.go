package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Room struct {
	ID          string `validate:"required,uuid"`
	Name        string `validate:"required,max=100"`
	Description string
	Capacity    int `validate:"gte=0"`
	CreatedAt   time.Time
}

func NewRoom(id, name, description string, capacity int, createdAt time.Time) Room {
	return Room{
		ID:          id,
		Name:        name,
		Description: description,
		Capacity:    capacity,
		CreatedAt:   createdAt,
	}
}

func (r *Room) Validate() error {

	roomValidator := validator.New()
	if err := roomValidator.Struct(r); err != nil {
		return fmt.Errorf("room validate struct: %w", err)
	}

	return nil
}
