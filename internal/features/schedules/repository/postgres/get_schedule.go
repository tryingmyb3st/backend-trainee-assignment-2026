package schedules_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *ScheduleRepository) GetScheduleByRoomId(ctx context.Context, roomID string) (*domain.Schedule, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, room_id, days_of_week, start_time, end_time
	FROM schedules
	WHERE room_id=$1
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, roomID)

	var model ScheduleModel
	err := row.Scan(
		&model.ID,
		&model.RoomID,
		&model.DaysOfWeek,
		&model.StartTime,
		&model.EndTime,
	)
	if err != nil {
		return nil, fmt.Errorf("row scan model: %w", err)
	}

	return domain.NewSchedule(
		model.ID,
		model.RoomID,
		model.DaysOfWeek,
		model.StartTime,
		model.EndTime,
	), nil
}
