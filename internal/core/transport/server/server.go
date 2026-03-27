package server

import (
	"backend-assignment-avito/docs"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/middleware"
	"context"
	"errors"
	"fmt"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type HTTPServer struct {
	mux *http.ServeMux
	cfg Config
	log *logger.Logger
}

func NewHTTPServer(config Config, logger *logger.Logger) *HTTPServer {
	return &HTTPServer{
		mux: http.NewServeMux(),
		cfg: config,
		log: logger,
	}
}

func (s *HTTPServer) RegisterRoutes(routes ...Route) {
	for _, route := range routes {

		pattern := fmt.Sprintf("%s %s", route.Method, route.URL)

		middlewares := []middleware.Middleware{
			middleware.PanicMiddleware(s.log),
			middleware.LogMiddleware(s.log),
			middleware.TraceMiddleware(),
		}

		middlewares = append(middlewares, route.AdditionalMiddleware...)

		handlerWithMiddleware := middleware.ChainMiddleware(
			route.Handler,
			middlewares...,
		)

		s.mux.Handle(pattern, handlerWithMiddleware)
	}
}

func (s *HTTPServer) RegisterSwagger() {
	s.mux.Handle(
		"/swagger/",
		httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")),
	)

	s.mux.HandleFunc("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(docs.SwaggerInfo.ReadDoc()))
	})
}

func (s *HTTPServer) Run(ctx context.Context) error {
	server := &http.Server{
		Addr:    s.cfg.Addr,
		Handler: s.mux,
	}

	s.mux.HandleFunc("GET /_info", HandleInfo)

	ch := make(chan error, 1)

	go func() {
		defer close(ch)

		s.log.Info("starting http server", zap.String("addr", s.cfg.Addr))
		err := server.ListenAndServe()

		if !errors.Is(http.ErrServerClosed, err) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("error from server: %w", err)
		}
	case <-ctx.Done():
		s.log.Warn("shutdown htttp server...")

		shutdownCtx, cancel := context.WithTimeout(ctx, s.cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}

		s.log.Warn("http server stopped")
	}

	return nil
}

func HandleInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
