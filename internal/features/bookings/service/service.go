package bookings_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
)

type BookingsService struct {
	slotsRepository   SlotsRepository
	bookingRepository BookingsRepository
}

type BookingsRepository interface {
	SaveNewBooking(ctx context.Context, newBooking domain.Booking) (*domain.Booking, error)
	GetUserBookings(ctx context.Context, userId string) ([]domain.Booking, error)
	GetBookingsWithPagination(ctx context.Context, pagination domain.Pagination) ([]domain.Booking, *domain.Pagination, error)
	CancelBooking(ctx context.Context, userId, bookingId string) (*domain.Booking, error)
}

type SlotsRepository interface {
	GetSlot(ctx context.Context, slotId string) (*domain.Slot, error)
}

func NewBookingsService(bookingRepo BookingsRepository, slotsRepo SlotsRepository) *BookingsService {
	return &BookingsService{
		slotsRepository:   slotsRepo,
		bookingRepository: bookingRepo,
	}
}
