package slots_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"time"
)

type SlotsService struct {
	slotsRepository    SlotsRepository
	scheduleRepository ScheduleRepository
	roomsRepository    RoomsRepository
}

type SlotsRepository interface {
	SaveSlots(ctx context.Context, slots []domain.Slot) ([]domain.Slot, error)
	GetActiveSlots(ctx context.Context, roomId string, date time.Weekday) ([]domain.Slot, error)
	IsSlotsBusy(ctx context.Context, roomId string, date time.Weekday) (bool, error)
}

type ScheduleRepository interface {
	GetScheduleByRoomId(ctx context.Context, roomID string) (*domain.Schedule, error)
}

type RoomsRepository interface {
	GetRoom(ctx context.Context, roomId string) (*domain.Room, error)
}

func NewSlotsService(
	slotsRepo SlotsRepository,
	scheduleRepo ScheduleRepository,
	roomsRepo RoomsRepository,
) *SlotsService {
	return &SlotsService{
		slotsRepository:    slotsRepo,
		scheduleRepository: scheduleRepo,
		roomsRepository:    roomsRepo,
	}
}
