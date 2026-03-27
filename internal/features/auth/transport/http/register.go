package auth_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type RegisterDTO struct {
	Email    string `json:"email" example:"useremail@gmail.com" extensions:"x-order=0"`
	Password string `json:"password" example:"supersecretpass" extensions:"x-order=1"`
	Role     string `json:"role" example:"admin" extensions:"x-order=2"`
}

type RegisterDTOResponse struct {
	ID        string    `json:"id"  example:"3fa85f64-5717-4562-b3fc-2c963f66afa6" extensions:"x-order=0"`
	Email     string    `json:"email" example:"useremail@gmail.com" extensions:"x-order=1"`
	Password  string    `json:"password,omitempty" extensions:"x-order=2"`
	Role      string    `json:"role" example:"admin" extensions:"x-order=3"`
	CreatedAt time.Time `json:"createdAt" example:"2026-03-25T12:00:41.267Z" extensions:"x-order=4"`
}

// RegisterUser godoc
// @Summary Регистрация пользователя
// @Description Создаёт нового пользователя и возвращает его данные. Реализация этого эндпоинта является дополнительным заданием. Для авторизации в рамках обязательной части используйте /dummyLogin. Доступен без авторизации.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterDTO true "тело запроса"
// @Success 201 {object} RegisterDTOResponse "Пользователь создан"
// @Failure 400 {object} domain.CustomError "Неверный запрос или email уже занят"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /register [post]
func (h *AuthHTTPHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	var req RegisterDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("error reguster decoding", zap.Error(err))

		respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
		return
	}

	user, err := h.AuthService.RegisterUser(
		ctx,
		domain.User{
			Email: req.Email,
			Role:  req.Role,
		},
		req.Password,
	)

	if err != nil {
		log.Error("get user after register", zap.Error(err))

		respWriter.MapError(err)
		return
	}

	resp := RegisterDTOResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}

	respWriter.JSONResponse(resp, http.StatusCreated)
}
