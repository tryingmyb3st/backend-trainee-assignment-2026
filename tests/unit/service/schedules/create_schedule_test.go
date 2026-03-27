package schedules_schedules_test

import (
	"backend-assignment-avito/internal/core/domain"
	schedules_service "backend-assignment-avito/internal/features/schedules/service"
	schedules_service_mock "backend-assignment-avito/mocks/schedules_service"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateSchedules(t *testing.T) {
	repoMock := schedules_service_mock.NewMockScheduleRepository(t)
	scheduleService := schedules_service.NewScheduleService(repoMock)

	ctx := context.Background()

	want := &domain.Schedule{
		RoomID:     "12772f07-5bc0-4eb1-a916-37e0dccf0962",
		DaysOfWeek: []int{3, 4},
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(2*time.Hour + 3*time.Minute),
	}

	repoMock.
		On("SaveSchedule", ctx, mock.Anything).
		Return(func(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error) {
			return &schedule, nil
		})

	schedule, err := scheduleService.CreateSchedule(ctx, *want)

	require.NoError(t, err)
	assert.Equal(t, 36, len([]rune(schedule.ID)))
	assert.Equal(t, "12772f07-5bc0-4eb1-a916-37e0dccf0962", schedule.RoomID)
}

func TestCreateSchedulesFail(t *testing.T) {
	repoMock := schedules_service_mock.NewMockScheduleRepository(t)
	scheduleService := schedules_service.NewScheduleService(repoMock)

	ctx := context.Background()

	want := &domain.Schedule{
		RoomID:     "1",
		DaysOfWeek: []int{3, 4},
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(2*time.Hour + 3*time.Minute),
	}

	_, err := scheduleService.CreateSchedule(ctx, *want)

	require.ErrorIs(t, err, domain.INVALID_REQUEST)
}
