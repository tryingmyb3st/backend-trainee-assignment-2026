package bookings_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"time"
)

type BookingDTO struct {
	ID        string    `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6" extensions:"x-order=0"`
	SlotID    string    `json:"slotId" example:"4hbdaf64-5717-4562-b3fc-2c963f66aba6" extensions:"x-order=1"`
	UserID    string    `json:"userId" example:"9hbd88864-5717-4562-b3fc-2c963f66af77" extensions:"x-order=3"`
	Status    string    `json:"status" example:"string" extensions:"x-order=4"`
	CreatedAt time.Time `json:"createdAt" example:"2026-03-25T12:00:41.267Z" extensions:"x-order=5"`
}

type PaginationDTO struct {
	Page     int `json:"page" example:"1" extensions:"x-order=0"`
	PageSize int `json:"pageSize" example:"10" extensions:"x-order=1"`
	Total    int `json:"total" example:"15" extensions:"x-order=2"`
}

type BookingDTOResponse struct {
	Booking BookingDTO `json:"booking"`
}

type BookingsDTOResponse struct {
	Bookings []BookingDTO `json:"bookings"`
}

type PaginationDTOResponse struct {
	Pagination PaginationDTO `json:"pagination"`
}

type BookingsDTOWithPagination struct {
	BookingsDTOResponse
	PaginationDTOResponse
}

func domainToDTO(bookings []domain.Booking) []BookingDTO {
	result := make([]BookingDTO, 0, len(bookings))

	for _, booking := range bookings {
		result = append(result, BookingDTO{
			ID:        booking.ID,
			SlotID:    booking.SlotID,
			UserID:    booking.UserID,
			Status:    booking.Status,
			CreatedAt: booking.CreatedAt,
		})
	}

	return result
}
