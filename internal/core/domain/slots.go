package domain

import "time"

type Slot struct {
	ID        string `validate:"required,uuid"`
	RoomID    string `validate:"required,uuid"`
	StartTime time.Time
	EndTime   time.Time
}

func NewSlot(id, roomID string, startTime, endTime time.Time) *Slot {
	return &Slot{
		ID:        id,
		RoomID:    roomID,
		StartTime: startTime,
		EndTime:   endTime,
	}
}
