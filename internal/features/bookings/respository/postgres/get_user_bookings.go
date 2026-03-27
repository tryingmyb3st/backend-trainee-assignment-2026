package bookings_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookingsRepository) GetUserBookings(ctx context.Context, userId string) ([]domain.Booking, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT bookings.id, bookings.slot_id, bookings.user_id, bookings.status, bookings.created_at
	FROM bookings
	JOIN slots ON slots.id = bookings.slot_id
	WHERE user_id=$1 AND slots.start_timestamp > NOW()
	`

	rows, err := r.ConnPool.Query(ctxTimeout, query, userId)
	if err != nil {
		return nil, fmt.Errorf("query get user bookings: %w", err)
	}
	defer rows.Close()

	bookings := make([]BookingModel, 0)
	for rows.Next() {
		var model BookingModel

		err := rows.Scan(
			&model.ID,
			&model.SlotID,
			&model.UserID,
			&model.Status,
			&model.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan model: %w", err)
		}

		bookings = append(bookings, model)
	}

	return modelsToDomain(bookings), nil
}

func modelsToDomain(models []BookingModel) []domain.Booking {
	result := make([]domain.Booking, 0, len(models))

	for _, model := range models {
		booking := domain.NewBooking(
			model.ID,
			model.SlotID,
			model.UserID,
			model.Status,
			model.CreatedAt,
		)

		result = append(result, *booking)
	}

	return result
}
