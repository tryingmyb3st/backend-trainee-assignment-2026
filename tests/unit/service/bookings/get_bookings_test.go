package bookings_service_test

import (
	"backend-assignment-avito/internal/core/domain"
	bookings_service "backend-assignment-avito/internal/features/bookings/service"
	bookings_service_mock "backend-assignment-avito/mocks/bookings_service"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAllBookings(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	ctx := context.Background()
	page := 2
	pageSize := 10

	wantBooking := []domain.Booking{
		{
			ID:        "1",
			SlotID:    "2",
			UserID:    "3",
			Status:    "active",
			CreatedAt: time.Now(),
		},
	}

	wantPagination := domain.Pagination{
		Page:     page,
		PageSize: pageSize,
		Total:    5,
	}

	bookingsRepoMock.
		On("GetBookingsWithPagination", ctx, domain.Pagination{
			Page:     page,
			PageSize: pageSize,
		}).
		Return(wantBooking, &wantPagination, nil).
		Once()

	bookings, pag, err := bookingsService.GetAllBookings(ctx, page, pageSize)

	require.NoError(t, err)
	assert.Equal(t, wantBooking, bookings)
	assert.Equal(t, wantPagination, *pag)
}

func TestGetAllBookingsFail1(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	ctx := context.Background()
	page := -1
	pageSize := 10

	_, _, err := bookingsService.GetAllBookings(ctx, page, pageSize)
	require.ErrorIs(t, err, domain.INVALID_REQUEST)
}

func TestGetAllBookingsFail2(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	expectedErr := errors.New("repo err")

	ctx := context.Background()
	page := 2
	pageSize := 10

	bookingsRepoMock.
		On("GetBookingsWithPagination", ctx, domain.Pagination{
			Page:     page,
			PageSize: pageSize,
		}).
		Return(nil, nil, expectedErr).
		Once()

	_, _, err := bookingsService.GetAllBookings(ctx, page, pageSize)

	require.ErrorIs(t, err, expectedErr)
}
