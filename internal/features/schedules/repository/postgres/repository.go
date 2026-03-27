package schedules_repository

import postgres_pool "backend-assignment-avito/internal/core/repository/postgres"

type ScheduleRepository struct {
	ConnPool postgres_pool.Pool
}

func NewScheduleRepository(conn postgres_pool.Pool) *ScheduleRepository {
	return &ScheduleRepository{
		ConnPool: conn,
	}
}
