package auth

import (
	errormiddleware "github.com/Shoyeb45/server/api/middleware/error-middleware"
	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
	sqlcv1 "github.com/Shoyeb45/server/pkg/repository/gen-queries"
	chi "github.com/go-chi/chi/v5"
)

func Mount(query *sqlcv1.Queries, r chi.Router) {
	authRepo := NewAuthRepository(query)
	authHandler := NewAuthHandler(*authRepo)

	r.Route("/auth", func(r chi.Router) {
		r.Route("/github", func(r chi.Router) {
			r.Get("/", errormiddleware.ErrorHandler(authHandler.RedirectGithub))
			r.With(validatormiddleware.Bind(validatormiddleware.FromQuery[GithubCallbackQuery]())).
				Get("/callback", errormiddleware.ErrorHandler(authHandler.GithubCallback))
		})
	})
}
