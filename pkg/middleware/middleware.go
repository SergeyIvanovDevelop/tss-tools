package middleware

import (
	"net/http"

	log "github.com/SergeyIvanovDevelop/tss-tools/pkg/logger"
)

type Middleware func(http.Handler) http.Handler
type ServiceHandlerFunc func(http.ResponseWriter, *http.Request)
type contextKey string

type MiddlewareConfig struct {
	Key string
}

var pkgLog log.Log

const pkgName = "internal/middleware"

func BuildConveyorMiddleware(cfg *MiddlewareConfig) func(handlerFunc ServiceHandlerFunc) http.HandlerFunc {
	var middlewares []Middleware
	middlewares = append(middlewares, logMiddleware)
	middlewares = append(middlewares, corsMiddleware)
	middlewares = append(middlewares, requestIDMiddleware)
	middlewares = append(middlewares, requestTimeoutMiddleware)
	middlewares = append(middlewares, gzipMiddleware)

	if cfg.Key != "" {
		hashMiddleware := buildHashMiddleware(cfg.Key)
		middlewares = append(middlewares, hashMiddleware)
	}

	// Inline ConveyorMiddleware
	return func(handlerFunc ServiceHandlerFunc) http.HandlerFunc {
		handlerWithMiddleware := conveyor(http.HandlerFunc(handlerFunc), middlewares...)
		return handlerWithMiddleware.ServeHTTP
	}
}

func conveyor(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}
