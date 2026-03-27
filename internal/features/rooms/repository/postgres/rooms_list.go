package rooms_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *RoomsRepository) GetRoomsList(ctx context.Context) ([]domain.Room, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, name, description, capacity, created_at
	FROM rooms
	`

	rows, err := r.ConnPool.Query(ctxTimeout, query)
	if err != nil {
		return nil, fmt.Errorf("pool query: %w", err)
	}
	defer rows.Close()

	var models []RoomModel

	for rows.Next() {
		var model RoomModel

		err := rows.Scan(
			&model.ID,
			&model.Name,
			&model.Description,
			&model.Capacity,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("rows scan: %w", err)
		}

		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows after scan: %w", err)
	}

	return modelsToDomain(models), nil
}

func modelsToDomain(models []RoomModel) []domain.Room {
	result := make([]domain.Room, 0, len(models))
	for _, model := range models {
		room := domain.NewRoom(
			model.ID,
			model.Name,
			model.Description,
			model.Capacity,
			model.CreatedAt,
		)

		result = append(result, room)
	}
	return result
}
