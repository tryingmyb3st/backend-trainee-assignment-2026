package bookings_repository

import "time"

type BookingModel struct {
	ID        string
	SlotID    string
	UserID    string
	Status    string
	CreatedAt time.Time
}
