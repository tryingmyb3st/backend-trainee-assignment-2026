package rooms_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type CreateRoomDTO struct {
	Name        string `json:"name" example:"work_room" extensions:"x-order=0"`
	Description string `json:"description" example:"some description..." extensions:"x-order=1"`
	Capacity    int    `json:"capacity" example:"3" extensions:"x-order=2"`
}

// CreateRoom godoc
// @Summary Создать переговорку (только admin)
// @Tags Rooms
// @Accept json
// @Produce json
// @Param request body CreateRoomDTO true "тело запроса"
// @Success 201 {object} RoomDTOResponse "Переговорка создана"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Не авторизован"
// @Failure 403 {object} domain.CustomError "Доступ запрещён (требуется роль admin)"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /rooms/create [post]
func (h *RoomsHTTPHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	var req CreateRoomDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("decoding create room", zap.Error(err))

		respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
		return
	}

	room, err := h.roomsService.CreateNewRoom(ctx, domain.Room{
		Name:        req.Name,
		Description: req.Description,
		Capacity:    req.Capacity,
	})

	if err != nil {
		log.Debug("error create room", zap.Error(err))

		if errors.Is(err, domain.INVALID_REQUEST) {
			respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
			return
		}

		respWriter.ErrorResponse(domain.INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	resp := RoomDTOResponse{
		Room: RoomDTO{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Capacity:    room.Capacity,
			CreatedAt:   room.CreatedAt,
		},
	}

	respWriter.JSONResponse(resp, http.StatusCreated)
}
