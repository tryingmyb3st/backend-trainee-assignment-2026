package bookings_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/middleware"
	"backend-assignment-avito/internal/core/transport/server"
	"context"
	"net/http"
)

type BookingsHTTPHandler struct {
	bookingService BookingsService
}

type BookingsService interface {
	BookRoomSlot(ctx context.Context, slotId, userId string) (*domain.Booking, error)
	GetUserBookings(ctx context.Context, userId string) ([]domain.Booking, error)
	GetAllBookings(ctx context.Context, page, pageSize int) ([]domain.Booking, *domain.Pagination, error)
	CancelBooking(ctx context.Context, userId, bookingId string) (*domain.Booking, error)
}

func NewBookingsHandler(serv BookingsService) *BookingsHTTPHandler {
	return &BookingsHTTPHandler{
		bookingService: serv,
	}
}

func (h *BookingsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			URL:     "/bookings/create",
			Handler: h.CreateNewBooking,
			AdditionalMiddleware: []middleware.Middleware{
				middleware.AuthMiddleware(),
				middleware.OnlyUserMiddleware(),
			},
		},
		{
			Method:  http.MethodGet,
			URL:     "/bookings/my",
			Handler: h.GetUserBooking,
			AdditionalMiddleware: []middleware.Middleware{
				middleware.AuthMiddleware(),
				middleware.OnlyUserMiddleware(),
			},
		},
		{
			Method:  http.MethodPost,
			URL:     "/bookings/{bookingId}/cancel",
			Handler: h.CancelUserBooking,
			AdditionalMiddleware: []middleware.Middleware{
				middleware.AuthMiddleware(),
				middleware.OnlyUserMiddleware(),
			},
		},
		{
			Method:  http.MethodGet,
			URL:     "/bookings/list",
			Handler: h.GetBookingsList,
			AdditionalMiddleware: []middleware.Middleware{
				middleware.AuthMiddleware(),
				middleware.AdminRightsMiddleware(),
			},
		},
	}
}
