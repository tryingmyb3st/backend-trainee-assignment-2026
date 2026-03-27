package schedules_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/middleware"
	"backend-assignment-avito/internal/core/transport/server"
	"context"
	"net/http"
)

type ScheduleHTTPHandler struct {
	scheduleService ScheduleService
}

type ScheduleService interface {
	CreateSchedule(ctx context.Context, schedule domain.Schedule) (*domain.Schedule, error)
}

func NewScheduleHandler(serv ScheduleService) *ScheduleHTTPHandler {
	return &ScheduleHTTPHandler{
		scheduleService: serv,
	}
}

func (h *ScheduleHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  http.MethodPost,
			URL:     "/rooms/{roomId}/schedule/create",
			Handler: h.CreateSchedule,
			AdditionalMiddleware: []middleware.Middleware{
				middleware.AuthMiddleware(),
				middleware.AdminRightsMiddleware(),
			},
		},
	}
}
