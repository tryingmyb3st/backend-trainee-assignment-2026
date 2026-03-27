package middleware

import (
	"backend-assignment-avito/internal/core/domain"
	"backend-assignment-avito/internal/core/logger"
	"backend-assignment-avito/internal/core/transport/response"
	"backend-assignment-avito/internal/utils/jwt_utils"
	"context"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, middlewares ...Middleware) http.Handler {
	if len(middlewares) == 0 {
		return h
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}

	return h
}

func LogMiddleware(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			log := log.With(zap.String("URL", r.URL.String()))

			ctx := context.WithValue(r.Context(), "log", log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func PanicMiddleware(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			handler := response.NewResponseHandler(log, w)

			defer func() {
				if err := recover(); err != nil {
					log.Error("panic happened", zap.Any("err", err))
					handler.PanicResponse(err, "panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func TraceMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := r.Context().Value("log").(*logger.Logger)

			before := time.Now()
			log.Info(">>> incoming request", zap.String("METHOD", r.Method))

			next.ServeHTTP(w, r)

			after := time.Now()

			log.Info("<<< request handled", zap.Duration("latency", time.Duration(after.Sub(before))))
		})
	}
}

func AuthMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := r.Context().Value("log").(*logger.Logger)
			respWriter := response.NewResponseHandler(log, w)

			token := r.Header.Get("Authorization")

			if token == "" || !strings.HasPrefix(token, "Bearer") {
				log.Debug("unathourized user", zap.String("Authorization", r.Header.Get("Authorization")))

				respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusUnauthorized)
				return
			}

			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := jwt_utils.VerifyJWTtoken(token)
			if err != nil {
				log.Debug("invalid token", zap.Any("claims", claims))

				respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
				return
			}
			log.Debug("user claims", zap.String("id", claims.UserId), zap.String("role", claims.Role))

			if claims.ExpiresAt.Before(time.Now()) {
				log.Debug("token expired", zap.Time("expiresAt", claims.ExpiresAt.Time))

				respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), "userId", claims.UserId)
			ctx = context.WithValue(ctx, "role", claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AdminRightsMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := r.Context().Value("log").(*logger.Logger)
			respWriter := response.NewResponseHandler(log, w)

			rights := r.Context().Value("role").(string)

			if !strings.EqualFold(rights, "admin") {
				log.Debug("no admin rights", zap.String("role", rights))

				respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func OnlyUserMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := r.Context().Value("log").(*logger.Logger)
			respWriter := response.NewResponseHandler(log, w)

			rights := r.Context().Value("role").(string)

			if !strings.EqualFold(rights, "user") {
				log.Debug("admins is not allowed", zap.String("role", rights))

				respWriter.ErrorResponse(domain.INVALID_REQUEST, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
