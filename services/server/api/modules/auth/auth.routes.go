package auth

import (
	errormiddleware "github.com/Shoyeb45/server/api/middleware/error-middleware"
	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
	chi "github.com/go-chi/chi/v5"
	sqlcv1 "github.com/Shoyeb45/server/pkg/repository/gen-queries"
)

func Mount(query *sqlcv1.Queries, r chi.Router) {
	authRepo := NewAuthRepository(query)
	authHandler := NewAuthHandler(*authRepo)

	r.Route("/auth", func(r chi.Router) {
		r.With(
			validatormiddleware.Bind(validatormiddleware.FromBody[RequestLogIn]()),
		).Post("/", errormiddleware.ErrorHandler(authHandler.Login))
	})
}
