package slots_repository

import postgres_pool "backend-assignment-avito/internal/core/repository/postgres"

type SlotsRepository struct {
	ConnPool postgres_pool.Pool
}

func NewSlotsRepository(conn postgres_pool.Pool) *SlotsRepository {
	return &SlotsRepository{
		ConnPool: conn,
	}
}
