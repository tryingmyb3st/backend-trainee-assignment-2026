package rooms_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
)

func (s *RoomsService) GetRoomsList(ctx context.Context) ([]domain.Room, error) {

	rooms, err := s.roomsRepository.GetRoomsList(ctx)
	if err != nil {
		return nil, fmt.Errorf("get rooms list: %w", err)
	}

	return rooms, nil
}
