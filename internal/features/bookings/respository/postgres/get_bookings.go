package bookings_repository

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (r *BookingsRepository) GetBookingsWithPagination(
	ctx context.Context,
	pagination domain.Pagination,
) ([]domain.Booking, *domain.Pagination, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, r.ConnPool.OpTimeout())
	defer cancel()

	query := `
	SELECT id, slot_id, user_id, status, created_at, COUNT(*) OVER() AS total_count
	FROM bookings
	LIMIT $1
	OFFSET $2
	`

	rows, err := r.ConnPool.Query(
		ctxTimeout,
		query,
		pagination.PageSize,
		pagination.Page*pagination.PageSize-pagination.PageSize,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("get bookings with pagination: %w", err)
	}
	defer rows.Close()

	models := make([]BookingModel, 0)
	var total int
	for rows.Next() {
		var model BookingModel

		err := rows.Scan(
			&model.ID,
			&model.SlotID,
			&model.UserID,
			&model.Status,
			&model.CreatedAt,
			&total,
		)
		if err != nil {
			return nil, nil, fmt.Errorf("scan bookings with pagination: %w", err)
		}

		models = append(models, model)
	}
	if rows.Err() != nil {
		return nil, nil, fmt.Errorf("err after scan rows: %w", err)
	}

	pagination.Total = total
	return modelsToDomain(models), &pagination, nil
}
