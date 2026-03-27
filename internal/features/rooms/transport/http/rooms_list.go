package rooms_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"net/http"

	"go.uber.org/zap"
)

// GetRoomsList godoc
// @Summary Список переговорок (admin и user)
// @Tags Rooms
// @Accept json
// @Produce json
// @Success 200 {object} RoomsDTOResponse "Список переговорок"
// @Failure 401 {object} domain.CustomError "Не авторизован"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /rooms/list [get]
func (h *RoomsHTTPHandler) GetRoomsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	rooms, err := h.roomsService.GetRoomsList(ctx)
	if err != nil {
		log.Debug("get rooms list", zap.Error(err))

		respWriter.ErrorResponse(domain.INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	resp := RoomsDTOResponse{
		Rooms: domainsToDTO(rooms),
	}

	respWriter.JSONResponse(resp, http.StatusOK)
}

func domainsToDTO(rooms []domain.Room) []RoomDTO {
	result := make([]RoomDTO, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, RoomDTO{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Capacity:    room.Capacity,
			CreatedAt:   room.CreatedAt,
		})
	}
	return result
}
