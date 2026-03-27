package schedules_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *ScheduleService) CreateSchedule(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error) {
	schedule.ID = uuid.New().String()

	if err := schedule.Validate(); err != nil {
		return nil, domain.INVALID_REQUEST
	}

	createdSchedule, err := s.scheduleRepository.SaveSchedule(ctx, schedule)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				return nil, fmt.Errorf("save schedule foreign key: %w", domain.NOT_FOUND)
			case "23505":
				return nil, fmt.Errorf("save schedule unique key: %w", domain.SCHEDULE_EXISTS)
			}
		}

		return nil, fmt.Errorf("save schedule to repo: %w", err)
	}

	return createdSchedule, nil
}
