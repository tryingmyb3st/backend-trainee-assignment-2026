package schedules_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
)

type ScheduleService struct {
	scheduleRepository ScheduleRepository
}

type ScheduleRepository interface {
	SaveSchedule(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error)
}

func NewScheduleService(
	scheduleRepo ScheduleRepository,
) *ScheduleService {
	return &ScheduleService{
		scheduleRepository: scheduleRepo,
	}
}
