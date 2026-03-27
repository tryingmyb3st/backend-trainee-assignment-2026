package auth_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
)

type DummyLoginDTO struct {
	Role string `json:"role" example:"user"`
}

type DummyLoginDTOResponse struct {
	Token string `json:"token"`
}

// GetDummyLogin godoc
// @Summary Получить тестовый JWT по роли
// @Description Выдаёт тестовый JWT для указанной роли (admin / user). Для каждой роли возвращается фиксированный UUID пользователя: один и тот же UUID для всех запросов с ролью admin и один и тот же UUID для роли user. Это обеспечивает стабильность при тестировании сценариев, требующих проверки владельца брони.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body DummyLoginDTO true "DummyLogin тело запроса"
// @Success 200 {object} DummyLoginDTOResponse "Тестовый токен"
// @Failure 400 {object} domain.CustomError "Неверный запрос (недопустимое значение роли)"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /dummyLogin [post]
func (h *AuthHTTPHandler) GetDummyLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	var req DummyLoginDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("dummy login decoding", zap.Error(err))

		respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
		return
	}

	jwt, err := h.AuthService.GetTestJWTByRole(domain.User{
		Role: req.Role,
	})

	if err != nil {
		log.Error("get test jwt", zap.Error(err))

		if errors.Is(err, domain.INVALID_REQUEST) {
			respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
			return
		}

		respWriter.ErrorResponse(domain.INTERNAL_ERROR, http.StatusInternalServerError)
		return
	}

	resp := DummyLoginDTOResponse{
		Token: *jwt,
	}

	respWriter.JSONResponse(resp, http.StatusOK)
}
