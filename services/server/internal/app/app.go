package app

import (
	"time"

	requestLogger "github.com/Shoyeb45/server/api/middleware/request-logger"
	"github.com/Shoyeb45/server/pkg/config"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	const timeout = 10 * time.Second

	// cors config
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			config.Cfg.OriginURL,
		},
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
		},
		ExposedHeaders: []string{
			"Link",
		},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Timeout(timeout))
	r.Use(requestLogger.RequestLogger) // logs request details
	r.Use(middleware.Recoverer)        // recovers from panic

	return r
}
