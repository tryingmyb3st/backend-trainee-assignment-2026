package schedules_transport

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ScheduleDTO struct {
	DaysOfWeek []int     `json:"daysOfWeek" example:"1" extensions:"x-order=0"`
	StartTime  time.Time `json:"startTime" example:"15:13" extensions:"x-order=1"`
	EndTime    time.Time `json:"endTime" example:"15:47" extensions:"x-order=2"`
}

func (s *ScheduleDTO) UnmarshalJSON(data []byte) error {
	type Schedule struct {
		DaysOfWeek []int  `json:"daysOfWeek"`
		StartTime  string `json:"startTime"`
		EndTime    string `json:"endTime"`
	}

	var res Schedule
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	parsedStartTime, err := time.Parse("15:04", res.StartTime)
	if err != nil {
		return err
	}

	parsedEndTime, err := time.Parse("15:04", res.EndTime)
	if err != nil {
		return err
	}

	s.DaysOfWeek = res.DaysOfWeek
	s.StartTime = parsedStartTime
	s.EndTime = parsedEndTime
	return nil
}

type ScheduleDTOResponse struct {
	ID         string    `json:"id"`
	RoomID     string    `json:"roomId"`
	DaysOfWeek []int     `json:"daysOfWeek"`
	StartTime  time.Time `json:"-"`
	EndTime    time.Time `json:"-"`
}

func (s ScheduleDTOResponse) MarshalJSON() ([]byte, error) {
	type Alias ScheduleDTOResponse

	startTimeString := s.StartTime.Format("15:04")
	endTimeString := s.EndTime.Format("15:04")

	return json.Marshal(&struct {
		*Alias
		StartTimeString string `json:"start_time"`
		EndTimeString   string `json:"end_time"`
	}{
		Alias:           (*Alias)(&s),
		StartTimeString: startTimeString,
		EndTimeString:   endTimeString,
	})
}

// CreateSchedule godoc
// @Summary Создать расписание переговорки (только admin, только один раз). Длительность слота 30 мин. После создания расписание изменить нельзя.
// @Description Доступно только роли admin. Расписание можно создать только один раз для каждой переговорки. Поле daysOfWeek должно содержать значения от 1 до 7 (1=Пн, 7=Вс). При передаче значений вне этого диапазона возвращается 400.
// @Tags Schedules
// @Accept json
// @Produce json
// @Param roomId path string true "Идентификатор переговорки"
// @@Param request body ScheduleDTO true "тело запроса"
// @Success 201 {object} ScheduleDTOResponse "Расписание сохранено"
// @Failure 400 {object} domain.CustomError "Неверный запрос (в т.ч. недопустимые значения daysOfWeek)"
// @Failure 401 {object} domain.CustomError "Не авторизован"
// @Failure 403 {object} domain.CustomError "Доступ запрещён (требуется роль admin)"
// @Failure 404 {object} domain.CustomError "Переговорка не найдена"
// @Failure 409 {object} domain.CustomError "Расписание для переговорки уже создано, изменение не допускается"
// @Failure 500 {object} domain.InternalError "Внутренняя ошибка сервера"
// @Router /rooms/{roomId}/schedule/create [post]
func (h *ScheduleHTTPHandler) CreateSchedule(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	log := ctx.Value("log").(*logger.Logger)
	respWriter := response.NewResponseHandler(log, w)

	roomID := r.PathValue("roomId")

	var req ScheduleDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Debug("decoding create room", zap.Error(err))

		respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
		return
	}

	createdSchedule, err := h.scheduleService.CreateSchedule(ctx, domain.Schedule{
		RoomID:     roomID,
		DaysOfWeek: req.DaysOfWeek,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
	})

	if err != nil {
		log.Debug("service create schedule", zap.Error(err))

		respWriter.MapError(err)
		return
	}

	dtoSchedule := ScheduleDTOResponse{
		ID:         createdSchedule.ID,
		RoomID:     createdSchedule.RoomID,
		DaysOfWeek: createdSchedule.DaysOfWeek,
		StartTime:  createdSchedule.StartTime,
		EndTime:    createdSchedule.EndTime,
	}
	respWriter.JSONResponse(dtoSchedule, http.StatusCreated)
}
