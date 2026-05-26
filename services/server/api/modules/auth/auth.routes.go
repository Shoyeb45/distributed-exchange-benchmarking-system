package auth

import (
	errormiddleware "github.com/Shoyeb45/server/api/middleware/error-middleware"
	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
	chi "github.com/go-chi/chi/v5"
)

func Mount(r chi.Router, h *AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.With(
			validatormiddleware.Bind(validatormiddleware.FromBody[LogIn]()),
		).Post("/", errormiddleware.ErrorHandler(h.Login))
	})
}
