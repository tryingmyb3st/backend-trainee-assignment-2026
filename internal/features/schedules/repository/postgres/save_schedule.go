package schedules_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *ScheduleRepository) SaveSchedule(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO schedules(id, room_id, days_of_week, start_time, end_time)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id, room_id, days_of_week, start_time, end_time
	`

	row := r.ConnPool.QueryRow(
		ctxTimeout,
		query,
		schedule.ID,
		schedule.RoomID,
		schedule.DaysOfWeek,
		schedule.StartTime,
		schedule.EndTime,
	)

	var model ScheduleModel
	err := row.Scan(
		&model.ID,
		&model.RoomID,
		&model.DaysOfWeek,
		&model.StartTime,
		&model.EndTime,
	)
	if err != nil {
		return nil, fmt.Errorf("scan returning schedule: %w", err)
	}

	return domain.NewSchedule(
		model.ID,
		model.RoomID,
		model.DaysOfWeek,
		model.StartTime,
		model.EndTime,
	), nil
}
