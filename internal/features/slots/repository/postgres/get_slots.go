package slots_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
	"time"
)

func (r *SlotsRepository) GetActiveSlots(ctx context.Context, roomId string, date time.Weekday) ([]domain.Slot, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT slots.id, slots.room_id, slots.start_timestamp, slots.end_timestamp
	FROM slots
	LEFT JOIN bookings ON bookings.slot_id = slots.id
	WHERE room_id=$1 AND EXTRACT(DOW FROM start_timestamp) = $2 
	AND (slot_id IS NULL OR status = 'cancelled');
	`

	rows, err := r.ConnPool.Query(
		ctxTimeout,
		query,
		roomId,
		int(date),
	)

	if err != nil {
		return nil, fmt.Errorf("pool query: %w", err)
	}
	defer rows.Close()

	var models []SlotModel

	for rows.Next() {
		var model SlotModel

		err := rows.Scan(
			&model.ID,
			&model.RoomID,
			&model.StartTime,
			&model.EndTime,
		)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows after scan: %w", err)
	}

	return r.modelsToDomains(models), nil
}
