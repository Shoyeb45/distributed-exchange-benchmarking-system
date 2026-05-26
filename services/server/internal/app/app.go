package app

import (
	"time"

	requestLogger "github.com/Shoyeb45/server/api/middleware/request-logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	const timeout = 10 * time.Second

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(timeout))
	r.Use(requestLogger.RequestLogger) // logs request details
	r.Use(middleware.Recoverer)        // recovers from panic

	// add swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
