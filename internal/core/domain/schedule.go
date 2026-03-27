package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

type Schedule struct {
	ID         string    `validate:"required,uuid"`
	RoomID     string    `validate:"required,uuid"`
	DaysOfWeek []int     `validate:"required"`
	StartTime  time.Time `validate:"required"`
	EndTime    time.Time `validate:"required"`
}

func NewSchedule(id, roomID string, daysOfWeek []int, startTime, endTime time.Time) *Schedule {
	return &Schedule{
		ID:         id,
		RoomID:     roomID,
		DaysOfWeek: daysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}
}

func (s *Schedule) Validate() error {
	scheduleValidator := validator.New()

	if err := scheduleValidator.Struct(s); err != nil {
		return fmt.Errorf("schedule validator struct: %w", err)
	}

	for _, day := range s.DaysOfWeek {
		if day < 1 || day > 7 {
			return errors.New("days of week between 1 and 7")
		}
	}

	return nil
}
