package slots_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type SlotDTOResponse struct {
	ID        string    `json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6" extensions:"x-order=0"`
	RoomID    string    `json:"roomId" example:"78a85f64-6717-4562-b3fc-2c963f6abfa6" extensions:"x-order=1"`
	StartTime time.Time `json:"start" example:"2026-03-25T14:01:22.680Z" extensions:"x-order=2"`
	EndTime   time.Time `json:"end" example:"2026-03-25T14:01:22.680Z" extensions:"x-order=3"`
}

type SlotsDTOResponse struct {
	Slots []SlotDTOResponse `json:"slots"`
}

// GetActiveSlots godoc
// @Summary Список доступных для бронирования слотов по переговорке и дате (admin и user). Наиболее нагруженный эндпоинт.
// @Description Возвращает слоты, не занятые активной бронью, для указанной переговорки на указанную дату. Все даты и время передаются и возвращаются в UTC. Параметр date является обязательным; при его отсутствии возвращается 400. Если у переговорки нет расписания — возвращается пустой список (переговорка считается всегда недоступной).
// @Tags Slots
// @Accept json
// @Produce json
// @Param roomId path string true "Идентификатор переговорки"
// @Param date query string true "Дата в формате ISO 8601 (например: 2024-06-10). Обязательный параметр; при отсутствии — 400"
// @Success 200 {object} SlotsDTOResponse "Список доступных слотов (не занятых активной бронью). Пустой список, если у переговорки нет расписания или на эту дату нет слотов."
// @Failure 400 {object} domain.CustomError "Неверный запрос (отсутствует или некорректен параметр date)"
// @Failure 401 {object} domain.CustomError "Не авторизован"
// @Failure 404 {object} domain.CustomError "Переговорка не найдена"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /rooms/{roomId}/slots/list [get]
func (h *SlotsHTTPHandler) GetActiveSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	roomID := r.PathValue("roomId")

	date, err := h.getQueryDate(r.URL.Query().Get("date"))
	if err != nil {
		log.Debug("invalid time query", zap.Error(err))

		respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
		return
	}

	slots, err := h.slotsService.GetActiveSlots(ctx, roomID, *date)
	if err != nil {
		log.Debug("error happend", zap.Error(err))

		respWriter.MapError(err)
		return
	}

	resp := SlotsDTOResponse{
		Slots: domainsToDTO(slots),
	}

	respWriter.JSONResponse(resp, http.StatusOK)
}

func (h *SlotsHTTPHandler) getQueryDate(query string) (*time.Time, error) {
	if query == "" {
		return nil, domain.INVALID_REQUEST
	}

	date, err := time.Parse("2006-01-02", query)
	if err != nil {
		return nil, domain.INVALID_REQUEST
	}

	return &date, nil
}

func domainsToDTO(domains []domain.Slot) []SlotDTOResponse {
	result := make([]SlotDTOResponse, 0, len(domains))

	for _, domain := range domains {
		result = append(result, SlotDTOResponse{
			ID:        domain.ID,
			RoomID:    domain.RoomID,
			StartTime: domain.StartTime,
			EndTime:   domain.EndTime,
		})
	}

	return result
}
