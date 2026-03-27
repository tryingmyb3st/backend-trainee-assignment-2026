package slots_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *SlotsRepository) GetSlot(ctx context.Context, slotId string) (*domain.Slot, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, room_id, start_timestamp, end_timestamp
	FROM slots
	WHERE id=$1
	`

	row := r.ConnPool.QueryRow(
		ctxTimeout,
		query,
		slotId,
	)

	var model SlotModel
	err := row.Scan(
		&model.ID,
		&model.RoomID,
		&model.StartTime,
		&model.EndTime,
	)
	if err != nil {
		return nil, fmt.Errorf("scan row: %w", err)
	}

	return domain.NewSlot(
		model.ID,
		model.RoomID,
		model.StartTime,
		model.EndTime,
	), nil
}
