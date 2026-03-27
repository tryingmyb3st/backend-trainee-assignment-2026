package bookings_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type SlotDTO struct {
	SlotID string `json:"slotId"`
}

// CreateNewBooking godoc
// @Summary Создать бронь на слот (только user). Опционально — запросить ссылку на конференцию.
// @Description Доступно только роли user. Администратор не может создавать брони (403). userId берётся из JWT-токена, а не из тела запроса. Если временной слот находится в прошлом (start < now), возвращается 400 (INVALID_REQUEST).
// @Tags Bookings
// @Accept json
// @Produce json
// @Param request body SlotDTO true "тело запроса"
// @Success 201 {object} BookingDTOResponse "Бронь создана"
// @Failure 400 {object} domain.CustomError "Неверный запрос"
// @Failure 401 {object} domain.CustomError "Не авторизован"
// @Failure 403 {object} domain.CustomError "Доступ запрещён (бронирование доступно только роли user)"
// @Failure 404 {object} domain.CustomError "Слот не найден"
// @Failure 409 {object} domain.CustomError "Слот уже занят"
// @Failure 500 {object} domain.CustomError "Внутренняя ошибка сервера"
// @Router /bookings/create [post]
func (h *BookingsHTTPHandler) CreateNewBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	userId := ctx.Value("userId").(string)

	var req SlotDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("invalid time query", zap.Error(err))

		respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
		return
	}

	createdBooking, err := h.bookingService.BookRoomSlot(ctx, req.SlotID, userId)
	if err != nil {
		log.Debug("error happened", zap.Error(err))

		respWriter.MapError(err)
		return
	}

	resp := BookingDTOResponse{
		Booking: BookingDTO{
			ID:        createdBooking.ID,
			SlotID:    createdBooking.SlotID,
			UserID:    createdBooking.UserID,
			Status:    createdBooking.Status,
			CreatedAt: createdBooking.CreatedAt,
		},
	}

	respWriter.JSONResponse(resp, http.StatusCreated)
}
