package e2e_test

import (
	"backend-assignment-avito/internal/core/logger"
	postgres_pool "backend-assignment-avito/internal/core/repository/postgres"
	"backend-assignment-avito/internal/core/transport/server"
	auth_repository "backend-assignment-avito/internal/features/auth/repository/postgres"
	auth_service "backend-assignment-avito/internal/features/auth/service"
	auth_transport "backend-assignment-avito/internal/features/auth/transport/http"
	bookings_repository "backend-assignment-avito/internal/features/bookings/respository/postgres"
	bookings_service "backend-assignment-avito/internal/features/bookings/service"
	bookings_transport "backend-assignment-avito/internal/features/bookings/transport/http"
	rooms_repository "backend-assignment-avito/internal/features/rooms/repository/postgres"
	rooms_service "backend-assignment-avito/internal/features/rooms/service"
	rooms_transport "backend-assignment-avito/internal/features/rooms/transport/http"
	schedules_repository "backend-assignment-avito/internal/features/schedules/repository/postgres"
	schedules_service "backend-assignment-avito/internal/features/schedules/service"
	schedules_transport "backend-assignment-avito/internal/features/schedules/transport/http"
	slots_repository "backend-assignment-avito/internal/features/slots/repository/postgres"
	slots_service "backend-assignment-avito/internal/features/slots/service"
	slots_transport "backend-assignment-avito/internal/features/slots/transport/http"
	"backend-assignment-avito/internal/utils/jwt_utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Setenv("POSTGRES_USER", "test_user")
	os.Setenv("POSTGRES_PASSWORD", "test_pass")
	os.Setenv("POSTGRES_DB", "test_db")
	os.Setenv("POSTGRES_HOST", "localhost:5432")
	os.Setenv("POSTGRES_TIMEOUT", "10s")
	os.Setenv("DATABASE_URL", "postgres://test_user:test_pass@localhost:5432/test_db?sslmode=disable")
	os.Setenv("HTTP_ADDR", ":8080")
	os.Setenv("HTTP_TIMEOUT", "40s")
	os.Setenv("LOG_LEVEL", "ERROR")
	os.Setenv("LOG_FOLDER", "out/logs")
	os.Setenv("TOKEN_SECRET", "test_secret")

	log, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		panic(err)
	}
	defer log.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	pool, err := postgres_pool.NewConnectionPool(ctx, postgres_pool.NewConfigMust())
	if err != nil {
		panic(err)
	}

	serv := server.NewHTTPServer(server.NewConfigMust(), log)

	authRepo := auth_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepo)
	authHandler := auth_transport.NewAuthHandler(authService)
	serv.RegisterRoutes(authHandler.Routes()...)

	roomsRepo := rooms_repository.NewRoomsRepository(pool)
	roomsService := rooms_service.NewRoomsService(roomsRepo)
	roomsHandler := rooms_transport.NewRoomsHandler(roomsService)
	serv.RegisterRoutes(roomsHandler.Routes()...)

	scheduleRepo := schedules_repository.NewScheduleRepository(pool)
	scheduleService := schedules_service.NewScheduleService(scheduleRepo)
	schedulesHandler := schedules_transport.NewScheduleHandler(scheduleService)
	serv.RegisterRoutes(schedulesHandler.Routes()...)

	slotsRepository := slots_repository.NewSlotsRepository(pool)
	slotsService := slots_service.NewSlotsService(slotsRepository, scheduleRepo, roomsRepo)
	slotsHandler := slots_transport.NewSlotsHandler(slotsService)
	serv.RegisterRoutes(slotsHandler.Routes()...)

	bookingsRepository := bookings_repository.NewBookingsRepository(pool)
	bookingsService := bookings_service.NewBookingsService(bookingsRepository, slotsRepository)
	bookingsHandler := bookings_transport.NewBookingsHandler(bookingsService)
	serv.RegisterRoutes(bookingsHandler.Routes()...)

	go func() {
		if err = serv.Run(ctx); err != nil {
			panic("serv run panic")
		}
	}()

	migrations, err := migrate.New(
		"file://../../migrations",
		os.Getenv("DATABASE_URL"),
	)

	if err != nil {
		panic(err)
	}

	if err := migrations.Migrate(1); err != nil {
		panic(err)
	}

	m.Run()

}

func DoHttpRequest(t *testing.T, method, pattern string, reqBody interface{}, token string) *http.Response {
	t.Helper()

	jsonBody, err := json.Marshal(reqBody)
	require.NoError(t, err)
	reqBytes := bytes.NewReader(jsonBody)

	req, err := http.NewRequest(method, "http://localhost:8080"+pattern, reqBytes)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err)

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	require.NoError(t, err)

	return resp
}

func TestWorkflow(t *testing.T) {

	t.Log("1) admin dummy login")
	loginBody := map[string]interface{}{
		"role": "admin",
	}

	resp := DoHttpRequest(t, "POST", "/dummyLogin", loginBody, "")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var jwtAdmin auth_transport.DummyLoginDTOResponse
	err := json.NewDecoder(resp.Body).Decode(&jwtAdmin)
	require.NoError(t, err)

	t.Log("2) create room")
	body := map[string]interface{}{
		"name":        "work_room",
		"description": "room for work...",
		"capacity":    5,
	}
	resp = DoHttpRequest(t, "POST", "/rooms/create", body, jwtAdmin.Token)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var room rooms_transport.RoomDTOResponse
	err = json.NewDecoder(resp.Body).Decode(&room)
	require.NoError(t, err)

	assert.Equal(t, "work_room", room.Room.Name)
	assert.Equal(t, "room for work...", room.Room.Description)
	assert.Equal(t, 5, room.Room.Capacity)
	assert.Equal(t, 36, len([]rune(room.Room.ID)))

	t.Log("3) create schedule for room")
	body = map[string]interface{}{
		"daysOfWeek": []int{1, 2, 7},
		"startTime":  "17:13",
		"endTime":    "19:10",
	}

	resp = DoHttpRequest(t, "POST", fmt.Sprintf("/rooms/%s/schedule/create", room.Room.ID), body, jwtAdmin.Token)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var schedule schedules_transport.ScheduleDTOResponse
	err = json.NewDecoder(resp.Body).Decode(&schedule)
	require.NoError(t, err)

	assert.Equal(t, []int{1, 2, 7}, schedule.DaysOfWeek)
	assert.Equal(t, room.Room.ID, schedule.RoomID)
	assert.Equal(t, 36, len([]rune(schedule.ID)))

	t.Log("4) user dummy login")
	userLoginBody := map[string]interface{}{
		"role": "user",
	}

	resp = DoHttpRequest(t, "POST", "/dummyLogin", userLoginBody, "")
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var jwtUser auth_transport.DummyLoginDTOResponse
	err = json.NewDecoder(resp.Body).Decode(&jwtUser)
	require.NoError(t, err)

	claims, err := jwt_utils.VerifyJWTtoken(jwtUser.Token)
	require.NoError(t, err)

	t.Log("4) slots list for user")

	resp = DoHttpRequest(
		t,
		"GET",
		fmt.Sprintf("/rooms/%s/slots/list?date=2026-03-29", room.Room.ID),
		nil,
		jwtUser.Token,
	)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var slots slots_transport.SlotsDTOResponse
	err = json.NewDecoder(resp.Body).Decode(&slots)
	require.NoError(t, err)

	assert.Equal(t, 3, len(slots.Slots))

	t.Log("5) user is booking second slot")
	body = map[string]interface{}{
		"slotId": slots.Slots[1].ID,
	}

	resp = DoHttpRequest(t, "POST", "/bookings/create", body, jwtUser.Token)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var booking bookings_transport.BookingDTOResponse
	err = json.NewDecoder(resp.Body).Decode(&booking)
	require.NoError(t, err)

	assert.Equal(t, 36, len([]rune(booking.Booking.ID)))
	assert.Equal(t, slots.Slots[1].ID, booking.Booking.SlotID)
	assert.Equal(t, claims.UserId, booking.Booking.UserID)
	assert.Equal(t, "active", booking.Booking.Status)

	t.Log("5) user is cancellig his slot")

	resp = DoHttpRequest(t, "POST", fmt.Sprintf("/bookings/%s/cancel", booking.Booking.ID), nil, jwtUser.Token)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	err = json.NewDecoder(resp.Body).Decode(&booking)
	require.NoError(t, err)

	assert.Equal(t, 36, len([]rune(booking.Booking.ID)))
	assert.Equal(t, slots.Slots[1].ID, booking.Booking.SlotID)
	assert.Equal(t, claims.UserId, booking.Booking.UserID)
	assert.Equal(t, "cancelled", booking.Booking.Status)

	t.Log("6) Workflow ended successfully")
}
