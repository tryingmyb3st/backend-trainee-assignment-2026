package rooms_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
)

type RoomsService struct {
	roomsRepository RoomsRepository
}

type RoomsRepository interface {
	SaveNewRoom(ctx context.Context, room domain.Room) (*domain.Room, error)
	GetRoomsList(ctx context.Context) ([]domain.Room, error)
}

func NewRoomsService(roomsRepo RoomsRepository) *RoomsService {
	return &RoomsService{
		roomsRepository: roomsRepo,
	}
}
