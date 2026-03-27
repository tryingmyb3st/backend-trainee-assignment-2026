package bookings_service_test

import (
	"backend-assignment-avito/internal/core/domain"
	bookings_service "backend-assignment-avito/internal/features/bookings/service"
	bookings_service_mock "backend-assignment-avito/mocks/bookings_service"
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBookRoomSlot(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	ctx := context.Background()
	slotId := "12772f07-5bc0-4eb1-a916-37e0dccf0962"
	userId := "f745545e-9ab9-4b33-8a15-739615fd9f60"

	wantSlot := domain.Slot{
		ID:        slotId,
		RoomID:    "4",
		StartTime: time.Now().Add(30 * time.Minute),
		EndTime:   time.Now().Add(60 * time.Minute),
	}

	wantBooking := domain.Booking{
		SlotID: slotId,
		UserID: userId,
		Status: "active",
	}

	slotsRepoMock.
		On("GetSlot", ctx, slotId).
		Return(&wantSlot, nil).
		Once()

	bookingsRepoMock.
		On("SaveNewBooking", ctx, mock.Anything).
		Return(&wantBooking, nil).
		Once()

	booking, err := bookingsService.BookRoomSlot(ctx, slotId, userId)

	require.NoError(t, err)
	assert.Equal(t, wantBooking.SlotID, booking.SlotID)
	assert.Equal(t, wantBooking.UserID, booking.UserID)
	assert.Equal(t, wantBooking.Status, "active")
}

func TestBookRoomSlotFail1(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	ctx := context.Background()
	slotId := "12772f07-5bc0-4eb1-a916-37e0dccf0962"
	userId := "f745545e-9ab9-4b33-8a15-739615fd9f60"

	wantSlot := domain.Slot{
		ID:        slotId,
		RoomID:    "4",
		StartTime: time.Date(2025, time.April, 17, 0, 0, 0, 0, time.UTC),
		EndTime:   time.Date(2025, time.April, 17, 0, 30, 0, 0, time.UTC),
	}

	slotsRepoMock.
		On("GetSlot", ctx, slotId).
		Return(&wantSlot, nil).
		Once()

	_, err := bookingsService.BookRoomSlot(ctx, slotId, userId)

	require.ErrorIs(t, err, domain.INVALID_REQUEST)
}

func TestBookRoomSlotFail2(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	ctx := context.Background()
	slotId := "12772f07-5bc0-4eb1-a916-37e0dccf0962"
	userId := "invalid"

	_, err := bookingsService.BookRoomSlot(ctx, slotId, userId)

	require.ErrorIs(t, err, domain.INVALID_REQUEST)
}

func TestBookRoomSlotFail3(t *testing.T) {
	bookingsRepoMock := bookings_service_mock.NewMockBookingsRepository(t)
	slotsRepoMock := bookings_service_mock.NewMockSlotsRepository(t)

	bookingsService := bookings_service.NewBookingsService(
		bookingsRepoMock,
		slotsRepoMock,
	)

	ctx := context.Background()
	slotId := "12772f07-5bc0-4eb1-a916-37e0dccf0962"
	userId := "f745545e-9ab9-4b33-8a15-739615fd9f60"

	slotsRepoMock.
		On("GetSlot", ctx, slotId).
		Return(nil, pgx.ErrNoRows).
		Once()

	_, err := bookingsService.BookRoomSlot(ctx, slotId, userId)

	require.ErrorIs(t, err, domain.NOT_FOUND)
}
