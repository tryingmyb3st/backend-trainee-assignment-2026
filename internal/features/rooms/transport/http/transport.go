package rooms_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/middleware"
	"backend-assignment-avito/internal/core/transport/server"
	"context"
	"net/http"
)

type RoomsHTTPHandler struct {
	roomsService RoomsService
}

type RoomsService interface {
	CreateNewRoom(ctx context.Context, room domain.Room) (*domain.Room, error)
	GetRoomsList(ctx context.Context) ([]domain.Room, error)
}

func NewRoomsHandler(roomsServ RoomsService) *RoomsHTTPHandler {
	return &RoomsHTTPHandler{
		roomsService: roomsServ,
	}
}

func (h *RoomsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			URL:     "/rooms/create",
			Handler: h.CreateRoom,
			AdditionalMiddleware: []middleware.Middleware{
				middleware.AuthMiddleware(),
				middleware.AdminRightsMiddleware(),
			},
		},
		{
			Method:  http.MethodGet,
			URL:     "/rooms/list",
			Handler: h.GetRoomsList,
			AdditionalMiddleware: []middleware.Middleware{
				middleware.AuthMiddleware(),
			},
		},
	}
}
