package bookings_transport

import (
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"net/http"

	"go.uber.org/zap"

	_ "backend-assignment-avito/internal/core/domain"
)

// GetUserBooking godoc
// @Summary Список броней текущего пользователя (только user)
// @Description Доступно только роли user. Возвращает брони пользователя, чей user_id содержится в JWT-токене. Возвращаются только брони на будущие слоты (start >= now); брони на уже прошедшие слоты в ответ не включаются.
// @Tags Bookings
// @Accept json
// @Produce json
// @Success 200 {object} BookingsDTOResponse "Список броней текущего пользователя"
// @Failure 401 {object} domain.CustomError "Не авторизован"
// @Failure 403 {object} domain.CustomError "Доступ запрещён (только user)"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /bookings/my [get]
func (h *BookingsHTTPHandler) GetUserBooking(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	userId := ctx.Value("userId").(string)

	bookings, err := h.bookingService.GetUserBookings(ctx, userId)
	if err != nil {
		log.Debug("error happened", zap.Error(err))

		respWriter.MapError(err)
		return
	}

	resp := BookingsDTOResponse{
		Bookings: domainToDTO(bookings),
	}

	respWriter.JSONResponse(resp, http.StatusOK)
}
