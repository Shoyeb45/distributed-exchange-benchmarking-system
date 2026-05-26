package app

import (
	"time"

	requestLogger "github.com/Shoyeb45/server/api/middleware/request-logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
