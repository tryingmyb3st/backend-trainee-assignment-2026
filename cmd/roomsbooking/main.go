package main

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
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	_ "backend-assignment-avito/docs"
)

// @title Rooms Booking Service
// @host 127.0.0.1:8080
func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	log, err := logger.NewLogger(logger.NewConfigMust())
	if err != nil {
		fmt.Println("error with initializing new logger")
		os.Exit(1)
	}
	defer log.Close()
	log.Debug("logger initialized successfully")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	pool, err := postgres_pool.NewConnectionPool(ctx, postgres_pool.NewConfigMust())
	if err != nil {
		log.Fatal("creating new pgxPool", zap.Error(err))
	}
	serv := server.NewHTTPServer(server.NewConfigMust(), log)

	log.Debug("initializing auth service")
	authRepo := auth_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepo)
	authHandler := auth_transport.NewAuthHandler(authService)
	serv.RegisterRoutes(authHandler.Routes()...)

	log.Debug("initializing rooms service")
	roomsRepo := rooms_repository.NewRoomsRepository(pool)
	roomsService := rooms_service.NewRoomsService(roomsRepo)
	roomsHandler := rooms_transport.NewRoomsHandler(roomsService)
	serv.RegisterRoutes(roomsHandler.Routes()...)

	log.Debug("initializing schedules service")
	scheduleRepo := schedules_repository.NewScheduleRepository(pool)
	scheduleService := schedules_service.NewScheduleService(scheduleRepo)
	schedulesHandler := schedules_transport.NewScheduleHandler(scheduleService)
	serv.RegisterRoutes(schedulesHandler.Routes()...)

	log.Debug("initializing slots service")
	slotsRepository := slots_repository.NewSlotsRepository(pool)
	slotsService := slots_service.NewSlotsService(slotsRepository, scheduleRepo, roomsRepo)
	slotsHandler := slots_transport.NewSlotsHandler(slotsService)
	serv.RegisterRoutes(slotsHandler.Routes()...)

	log.Debug("initializing bookings service")
	bookingsRepository := bookings_repository.NewBookingsRepository(pool)
	bookingsService := bookings_service.NewBookingsService(bookingsRepository, slotsRepository)
	bookingsHandler := bookings_transport.NewBookingsHandler(bookingsService)
	serv.RegisterRoutes(bookingsHandler.Routes()...)

	serv.RegisterSwagger()

	if err = serv.Run(ctx); err != nil {
		log.Error("error occured in server", zap.Error(err))
	}
}
