package slots_service_test

import (
	"backend-assignment-avito/internal/core/domain"
	slots_service "backend-assignment-avito/internal/features/slots/service"
	slots_service_mock "backend-assignment-avito/mocks/slots_service"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetActiveSlots(t *testing.T) {
	scheduleRepoMock := slots_service_mock.NewMockScheduleRepository(t)
	slotsRepoMock := slots_service_mock.NewMockSlotsRepository(t)
	roomsRepoMock := slots_service_mock.NewMockRoomsRepository(t)

	slotsService := slots_service.NewSlotsService(
		slotsRepoMock,
		scheduleRepoMock,
		roomsRepoMock,
	)

	ctx := context.Background()
	date, _ := time.Parse("2014-09-12", "2026-03-25")

	wantSchedule := domain.NewSchedule(
		"1",
		"2",
		[]int{1, 3, 5},
		time.Now(),
		time.Now().Add(2*time.Hour+3*time.Minute),
	)

	scheduleRepoMock.
		On("GetScheduleByRoomId", ctx, "2").
		Return(wantSchedule, nil).
		Once()

	slotsRepoMock.
		On("GetActiveSlots", ctx, "2", mock.Anything).
		Return([]domain.Slot{}, nil).
		Once()

	slotsRepoMock.
		On("SaveSlots", ctx, mock.MatchedBy(func(arg []domain.Slot) bool {
			return true
		})).
		Return(func(ctx context.Context, slots []domain.Slot) ([]domain.Slot, error) {
			return slots, nil
		})

	slotsRepoMock.
		On("IsSlotsBusy", ctx, mock.Anything, mock.Anything).
		Return(false, nil)

	slots, err := slotsService.GetActiveSlots(ctx, "2", date)

	require.NoError(t, err)
	assert.Equal(t, 4, len(slots))
	assert.Equal(t, "2", slots[0].RoomID)
}
