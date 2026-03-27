package rooms_service

import (
	"backend-assignment-avito/internal/core/domain"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (s *RoomsService) CreateNewRoom(ctx context.Context, room domain.Room) (*domain.Room, error) {
	room.ID = uuid.New().String()
	room.CreatedAt = time.Now()

	if err := room.Validate(); err != nil {
		return nil, fmt.Errorf("error validating room: %w(%v)", domain.INVALID_REQUEST, err)
	}

	createdRoom, err := s.roomsRepository.SaveNewRoom(ctx, room)
	if err != nil {
		return nil, fmt.Errorf("error saving to repository: %v", err)
	}

	return createdRoom, nil
}
