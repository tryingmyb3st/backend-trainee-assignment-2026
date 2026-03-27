package bookings_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"fmt"
	"net/http"
	"strconv"

	"go.uber.org/zap"
)

// GetBookingsList godoc
// @Summary Список всех броней с пагинацией (только admin)
// @Description Доступно только роли admin. Поддерживает пагинацию через параметры page и pageSize. Оба параметра опциональны; значения по умолчанию: page=1, pageSize=20. Максимальное значение pageSize — 100.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param page query integer false "Номер страницы (начиная с 1). По умолчанию 1."
// @Param pageSize query integer false "Количество записей на странице. По умолчанию 20, максимум 100."
// @Success 200 {object} BookingsDTOWithPagination "Cписок всех броней"
// @Failure 400 {object} domain.CustomError "Неверный запрос (некорректные параметры пагинации)"
// @Failure 401 {object} domain.CustomError "Не авторизован"
// @Failure 403 {object} domain.CustomError "Доступ запрещён (только admin)"
// @Failure 500 {object} domain.CustomError "Внутренняя ошибка сервера"
// @Router /bookings/list [get]
func (h *BookingsHTTPHandler) GetBookingsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	page, err := getQueryInt("page", r, 1)
	if err != nil {
		log.Debug("error get page query", zap.Error(err))

		respWriter.MapError(domain.INVALID_REQUEST)
		return
	}

	pageSize, err := getQueryInt("pageSize", r, 20)
	if err != nil {
		log.Debug("error get pageSize query", zap.Error(err))

		respWriter.MapError(domain.INVALID_REQUEST)
		return
	}

	bookings, pagination, err := h.bookingService.GetAllBookings(ctx, *page, *pageSize)
	if err != nil {
		log.Debug("error get all bookings", zap.Error(err))

		respWriter.MapError(domain.INVALID_REQUEST)
		return
	}

	resp := BookingsDTOWithPagination{
		BookingsDTOResponse: BookingsDTOResponse{
			Bookings: domainToDTO(bookings),
		},
		PaginationDTOResponse: PaginationDTOResponse{
			Pagination: PaginationDTO{
				Page:     pagination.Page,
				PageSize: pagination.PageSize,
				Total:    pagination.Total,
			},
		},
	}

	respWriter.JSONResponse(resp, http.StatusOK)
}

func getQueryInt(queryName string, r *http.Request, defaultValue int) (*int, error) {
	value := r.URL.Query().Get(queryName)

	if value == "" {
		return &defaultValue, nil
	}

	valueInt, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("parsing string to int: %w", err)
	}

	return &valueInt, nil
}
