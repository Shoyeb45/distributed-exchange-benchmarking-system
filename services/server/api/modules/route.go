package modules

import (
	"net/http"

	errormiddleware "github.com/Shoyeb45/server/api/middleware/error-middleware"
	"github.com/Shoyeb45/server/api/modules/auth"
	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/database"
	sqlcv1 "github.com/Shoyeb45/server/pkg/repository/gen-queries"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Mount all the application routes
func MountRoutes(r chi.Router) {
	query := sqlcv1.New(database.DB)

	r.NotFound(errormiddleware.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		return apierr.NewNotFound("route " + r.Method + " " + r.URL.Path + " does not exist")
	}))

	r.MethodNotAllowed(errormiddleware.ErrorHandler(func(w http.ResponseWriter, r *http.Request) error {
		return apierr.NewMethodNotAllowed("method " + r.Method + " is not allowed on " + r.URL.Path)
	}))

	// add swagger
	r.Get("/api-docs/*", httpSwagger.WrapHandler)

	r.Route("/api/", func(r chi.Router) {
		auth.Mount(query, r)
	})
}
