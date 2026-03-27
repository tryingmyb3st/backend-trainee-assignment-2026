package domain

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Booking struct {
	ID     string `validate:"required,uuid"`
	SlotID string `validate:"required,uuid"`
	UserID string `validate:"required,uuid"`
	Status string `validate:"oneof=active cancelled"`
	// conference link
	CreatedAt time.Time
}

func NewBooking(id, slotId, userId, status string, createdAt time.Time) *Booking {
	return &Booking{
		ID:        id,
		SlotID:    slotId,
		UserID:    userId,
		Status:    status,
		CreatedAt: createdAt,
	}
}

func (b *Booking) Validate() error {
	bookingValidator := validator.New()
	if err := bookingValidator.Struct(b); err != nil {
		return fmt.Errorf("validate struct: %w", err)
	}

	return nil
}
