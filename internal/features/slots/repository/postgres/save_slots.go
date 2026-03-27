package slots_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
	"time"
)

type SlotModel struct {
	ID        string
	RoomID    string
	StartTime time.Time
	EndTime   time.Time
}

func (r *SlotsRepository) SaveSlots(ctx context.Context, slots []domain.Slot) ([]domain.Slot, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO slots(id, room_id, start_timestamp, end_timestamp)
	VALUES($1,$2,$3,$4)
	RETURNING id, room_id, start_timestamp, end_timestamp;
	`

	models := make([]SlotModel, 0, len(slots))

	for _, slot := range slots {
		var model SlotModel

		row := r.ConnPool.QueryRow(
			ctxTimeout,
			query,
			slot.ID,
			slot.RoomID,
			slot.StartTime,
			slot.EndTime,
		)

		err := row.Scan(
			&model.ID,
			&model.RoomID,
			&model.StartTime,
			&model.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		models = append(models, model)
	}

	return r.modelsToDomains(models), nil
}
