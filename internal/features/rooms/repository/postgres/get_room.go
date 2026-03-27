package rooms_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *RoomsRepository) GetRoom(ctx context.Context, roomId string) (*domain.Room, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, name, description, capacity, created_at
	FROM rooms
	WHERE id=$1
	`

	row := r.ConnPool.QueryRow(ctxTimeout, query, roomId)

	var model RoomModel
	err := row.Scan(
		&model.ID,
		&model.Name,
		&model.Description,
		&model.Capacity,
		&model.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("row scan: %w", err)
	}

	return &domain.Room{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Capacity:    model.Capacity,
		CreatedAt:   model.CreatedAt,
	}, nil
}
