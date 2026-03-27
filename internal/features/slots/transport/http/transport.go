package slots_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/middleware"
	"backend-assignment-avito/internal/core/transport/server"
	"context"
	"time"
)

type SlotsHTTPHandler struct {
	slotsService SlotsService
}

type SlotsService interface {
	GetActiveSlots(ctx context.Context, roomId string, date time.Time) ([]domain.Slot, error)
}

func NewSlotsHandler(serv SlotsService) *SlotsHTTPHandler {
	return &SlotsHTTPHandler{
		slotsService: serv,
	}
}

func (h *SlotsHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  "GET",
			URL:     "/rooms/{roomId}/slots/list",
			Handler: h.GetActiveSlots,
			AdditionalMiddleware: []middleware.Middleware{
				middleware.AuthMiddleware(),
			},
		},
	}
}
