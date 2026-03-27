package rooms_repository

import postgres_pool "backend-assignment-avito/internal/core/repository/postgres"

type RoomsRepository struct {
	ConnPool postgres_pool.Pool
}

func NewRoomsRepository(conn postgres_pool.Pool) *RoomsRepository {
	return &RoomsRepository{
		ConnPool: conn,
	}
}
