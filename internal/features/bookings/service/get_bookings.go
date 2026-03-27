package bookings_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (s *BookingsService) GetAllBookings(ctx context.Context, page, pageSize int) ([]domain.Booking, *domain.Pagination, error) {
	paginationSettings := domain.NewPagination(page, pageSize, 0)

	if err := paginationSettings.Validate(); err != nil {
		return nil, nil, domain.INVALID_REQUEST
	}

	bookings, pagination, err := s.bookingRepository.GetBookingsWithPagination(ctx, *paginationSettings)
	if err != nil {
		return nil, nil, fmt.Errorf("get bookings with pagination: %w", err)
	}

	return bookings, pagination, nil
}
