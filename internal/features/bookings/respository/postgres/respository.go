package bookings_repository

import postgres_pool "backend-assignment-avito/internal/core/repository/postgres"

type BookingsRepository struct {
	ConnPool postgres_pool.Pool
}

func NewBookingsRepository(conn postgres_pool.Pool) *BookingsRepository {
	return &BookingsRepository{
		ConnPool: conn,
	}
}
