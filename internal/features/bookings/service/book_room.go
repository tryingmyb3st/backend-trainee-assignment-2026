package bookings_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *BookingsService) BookRoomSlot(ctx context.Context, slotId, userId string) (*domain.Booking, error) {

	newBooking := domain.NewBooking(
		uuid.NewString(),
		slotId,
		userId,
		"active",
		time.Now(),
	)

	if err := newBooking.Validate(); err != nil {
		return nil, fmt.Errorf("new booking validate: %w", domain.INVALID_REQUEST)
	}

	slot, err := s.slotsRepository.GetSlot(ctx, slotId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("slot does not exists: %w", domain.NOT_FOUND)
		}

		return nil, fmt.Errorf("get slot from repo: %w", err)
	}

	if slot.StartTime.Before(time.Now()) {
		return nil, fmt.Errorf("booking in past is not allowed: %w", domain.INVALID_REQUEST)
	}

	createdBooking, err := s.bookingRepository.SaveNewBooking(ctx, *newBooking)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, fmt.Errorf("slot is already booked: %w", domain.SLOT_ALREADY_BOOKED)
			}
		}

		return nil, err
	}

	return createdBooking, nil
}
