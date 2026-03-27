package slots_repository

import (
	"context"
	"fmt"
	"time"
)

func (r *SlotsRepository) IsSlotsBusy(ctx context.Context, roomId string, date time.Weekday) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT slots.id, slots.room_id, slots.start_timestamp, slots.end_timestamp
	FROM slots
	WHERE room_id=$1 AND EXTRACT(DOW FROM start_timestamp) = $2;
	`

	rows, err := r.ConnPool.Query(
		ctxTimeout,
		query,
		roomId,
		int(date),
	)

	if err != nil {
		return false, fmt.Errorf("pool query: %w", err)
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
			return false, fmt.Errorf("rows scan: %w", err)
		}

		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("rows after scan: %w", err)
	}

	return len(models) > 0, nil
}
