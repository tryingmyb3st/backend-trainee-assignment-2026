package bookings_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (s *BookingsService) CancelBooking(ctx context.Context, userId, bookingId string) (*domain.Booking, error) {
	booking, err := s.bookingRepository.CancelBooking(ctx, userId, bookingId)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("booking not found: %w", domain.NOT_FOUND)
		}

		if errors.Is(err, domain.FORBIDDEN) {
			return nil, fmt.Errorf("another user booking: %w", err)
		}

		return nil, err
	}

	return booking, nil
}
