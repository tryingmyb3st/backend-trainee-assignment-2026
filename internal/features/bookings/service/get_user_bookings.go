package bookings_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookingsService) GetUserBookings(ctx context.Context, userId string) ([]domain.Booking, error) {
	bookings, err := s.bookingRepository.GetUserBookings(ctx, userId)

	if err != nil {
		return nil, fmt.Errorf("get user bookings from repo: %w", err)
	}

	return bookings, nil
}
