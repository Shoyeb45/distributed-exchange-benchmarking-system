package modules

import (
	"github.com/Shoyeb45/server/api/modules/auth"
	"github.com/Shoyeb45/server/pkg/database"
	sqlcv1 "github.com/Shoyeb45/server/pkg/repository/gen-queries"
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
	query := sqlcv1.New(database.DB)

	authRepo := auth.NewAuthRepository(query)
	authHandler := auth.NewAuthHandler(*authRepo)

	r.Route("/api/", func(r chi.Router) {
		auth.Mount(r, authHandler)
	});
}
