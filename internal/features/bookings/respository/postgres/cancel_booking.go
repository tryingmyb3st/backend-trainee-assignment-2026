package bookings_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookingsRepository) CancelBooking(ctx context.Context, userId, bookingId string) (*domain.Booking, error) {
	ctxTimeoutGet, cancelGet := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancelGet()

	queryGet := `
	SELECT id, slot_id, user_id, status, created_at
	FROM bookings
	WHERE id=$1
	`

	row := r.ConnPool.QueryRow(ctxTimeoutGet, queryGet, bookingId)

	var model BookingModel
	err := row.Scan(
		&model.ID,
		&model.SlotID,
		&model.UserID,
		&model.Status,
		&model.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan booking model: %w", err)
	}

	if model.UserID != userId {
		return nil, domain.FORBIDDEN
	}

	queryUpdate := `
	UPDATE bookings
	SET status='cancelled'
	WHERE id=$1
	RETURNING id, slot_id, user_id, status, created_at
	`

	ctxTimeoutUpdate, cancelUpdate := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancelUpdate()

	row = r.ConnPool.QueryRow(ctxTimeoutUpdate, queryUpdate, bookingId)

	err = row.Scan(
		&model.ID,
		&model.SlotID,
		&model.UserID,
		&model.Status,
		&model.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update booking: %w", err)
	}

	return domain.NewBooking(
		model.ID,
		model.SlotID,
		model.UserID,
		model.Status,
		model.CreatedAt,
	), nil
}
