package bookings_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookingsRepository) SaveNewBooking(ctx context.Context, newBooking domain.Booking) (*domain.Booking, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO bookings(id, slot_id, user_id, status, created_at)
	VALUES($1,$2,$3,$4,$5)
	RETURNING id, slot_id, user_id, status, created_at
	`

	row := r.ConnPool.QueryRow(
		ctxTimeout,
		query,
		newBooking.ID,
		newBooking.SlotID,
		newBooking.UserID,
		newBooking.Status,
		newBooking.CreatedAt,
	)

	var model BookingModel
	err := row.Scan(
		&model.ID,
		&model.SlotID,
		&model.UserID,
		&model.Status,
		&model.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("scan returning schedule: %w", err)
	}

	return domain.NewBooking(
		model.ID,
		model.SlotID,
		model.UserID,
		model.Status,
		model.CreatedAt,
	), nil
}
