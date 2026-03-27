package bookings_transport

import (
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"net/http"

	_ "backend-assignment-avito/internal/core/domain"

	"go.uber.org/zap"
)

// CancelUserBooking godoc
// @Summary Отменить бронь (только своя бронь, только user)
// @Description Доступно только роли user. Пользователь может отменить только свою бронь. Операция идемпотентна: повторный вызов на уже отменённой брони не является ошибкой и возвращает 200 с актуальным состоянием брони (status: cancelled). При попытке отменить чужую бронь — 403. При отсутствии брони с указанным ID — 404.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param bookingId path string true "Идентификатор брони"
// @Success 200 {object} BookingDTOResponse "Бронь отменена (или уже была отменена ранее)"
// @Failure 401 {object} domain.CustomError "Не авторизован"
// @Failure 403 {object} domain.CustomError "Не своя бронь или роль не user"
// @Failure 404 {object} domain.CustomError "Бронь не найдена
// @Failure 500 {object} domain.CustomError "Внутренняя ошибка сервера"
// @Router /bookings/{bookingId}/cancel [post]
func (h *BookingsHTTPHandler) CancelUserBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	userId := ctx.Value("userId").(string)
	bookingId := r.PathValue("bookingId")

	cancelledBooking, err := h.bookingService.CancelBooking(ctx, userId, bookingId)
	if err != nil {
		log.Debug("error happened", zap.Error(err))

		respWriter.MapError(err)
		return
	}

	resp := BookingDTOResponse{
		Booking: BookingDTO{
			ID:        cancelledBooking.ID,
			SlotID:    cancelledBooking.SlotID,
			UserID:    cancelledBooking.UserID,
			Status:    cancelledBooking.Status,
			CreatedAt: cancelledBooking.CreatedAt,
		},
	}

	respWriter.JSONResponse(resp, http.StatusOK)
}
