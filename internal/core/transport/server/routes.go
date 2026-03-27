package server

import (
	"backend-assignment-avito/internal/core/middleware"
	"net/http"
)

type Route struct {
	Method               string
	URL                  string
	Handler              http.HandlerFunc
	AdditionalMiddleware []middleware.Middleware
}
