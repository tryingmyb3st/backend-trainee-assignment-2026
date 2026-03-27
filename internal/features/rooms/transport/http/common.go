package rooms_transport

import "time"

type RoomDTO struct {
	ID          string    `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6" extensions:"x-order=0"`
	Name        string    `json:"name" example:"work_room" extensions:"x-order=1"`
	Description string    `json:"description" example:"room for work..." extensions:"x-order=2"`
	Capacity    int       `json:"capacity" example:"3" extensions:"x-order=3"`
	CreatedAt   time.Time `json:"createdAt" example:"2026-03-25T12:00:41.267Z" extensions:"x-order=4"`
}

type RoomsDTOResponse struct {
	Rooms []RoomDTO `json:"rooms"`
}

type RoomDTOResponse struct {
	Room RoomDTO `json:"room"`
}
