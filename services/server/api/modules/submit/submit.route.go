package submit

import (
	authmiddleware "github.com/Shoyeb45/server/api/middleware/auth"
	errormiddleware "github.com/Shoyeb45/server/api/middleware/error-middleware"
	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
	sqlcv1 "github.com/Shoyeb45/server/pkg/repository/gen-queries"
	"github.com/go-chi/chi/v5"
)

func Mount(query *sqlcv1.Queries, r chi.Router) {
	repo := NewSubmitRepository(query)
	handler := NewSubmitHandler(repo)

	r.With(authmiddleware.RequireAuth).
		Route("/submit", func(r chi.Router) {
			r.With(validatormiddleware.Bind(validatormiddleware.FromBody[CreateSubmitRequest]())).
				Post("/", errormiddleware.ErrorHandler(handler.CreateSubmission))
		})
}
