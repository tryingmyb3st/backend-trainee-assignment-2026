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

func TestGetUserBookings(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	ctx := context.Background()
	userId := "some_id_1"

	want := []domain.Booking{
		{
			ID:        "1",
			SlotID:    "2",
			UserID:    userId,
			Status:    "active",
			CreatedAt: time.Now(),
		},
	}

	bookingsRepoMock.
		On("GetUserBookings", ctx, userId).
		Return(want, nil).
		Once()

	result, err := bookingsService.GetUserBookings(ctx, userId)

	require.NoError(t, err)
	assert.Equal(t, want, result)
}

func TestGetUserBookingsFail(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	expectedErr := errors.New("repo error")

	ctx := context.Background()
	userId := "some_id_1"

	bookingsRepoMock.
		On("GetUserBookings", ctx, userId).
		Return(nil, expectedErr).
		Once()

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	_, err := bookingsService.GetUserBookings(ctx, userId)

	require.ErrorIs(t, err, expectedErr)
}
