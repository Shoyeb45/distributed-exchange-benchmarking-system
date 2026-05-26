package modules

import (
	"encoding/json"
	"net/http"

	"github.com/Shoyeb45/server/api/modules/auth"
	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/database"
	sqlcv1 "github.com/Shoyeb45/server/pkg/repository/gen-queries"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func MountRoutes(r chi.Router) {
	query := sqlcv1.New(database.DB)

	authRepo := auth.NewAuthRepository(query)
	authHandler := auth.NewAuthHandler(*authRepo)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(apierr.APIError{
			Code:    apierr.CodeNotFound,
			Message: "route " + r.Method + " " + r.URL.Path + " does not exist",
		})
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(apierr.APIError{
			Code:    apierr.CodeBadRequest,
			Message: "method " + r.Method + " is not allowed on " + r.URL.Path,
		})
	})
	
	// add swagger
	r.Get("/docs/*", httpSwagger.WrapHandler)

	r.Route("/api/", func(r chi.Router) {
		auth.Mount(r, authHandler)
	})
}
