package rooms_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

type RoomModel domain.Room

func (r *RoomsRepository) SaveNewRoom(ctx context.Context, room domain.Room) (*domain.Room, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO rooms(id, name, description, capacity, created_at)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id, name, description, capacity, created_at
	`

	row := r.ConnPool.QueryRow(
		ctxTimeout,
		query,
		room.ID,
		room.Name,
		room.Description,
		room.Capacity,
		room.CreatedAt,
	)

	var model RoomModel
	err := row.Scan(
		&model.ID,
		&model.Name,
		&model.Description,
		&model.Capacity,
		&model.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("error returning row scan: %w", err)
	}

	roomDomain := domain.NewRoom(
		model.ID,
		model.Name,
		model.Description,
		model.Capacity,
		model.CreatedAt,
	)

	return &roomDomain, nil
}
