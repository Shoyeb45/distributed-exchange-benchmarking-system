package app

import (
	"time"

	requestLogger "github.com/Shoyeb45/fast-docs/api/middleware/request-logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)



func New() *chi.Mux {
	r := chi.NewRouter()

	const timeout = 10 * time.Second

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(timeout))
	r.Use(requestLogger.RequestLogger) // logs request details
	r.Use(middleware.Recoverer)        // recovers from panic

	return r
}
