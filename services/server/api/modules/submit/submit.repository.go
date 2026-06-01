package submit

import (
	"context"

	queries "github.com/Shoyeb45/server/pkg/repository/gen-queries"
	"github.com/jackc/pgx/v5/pgtype"
)

type SubmitRepository struct {
	q *queries.Queries
}

func NewSubmitRepository(q *queries.Queries) *SubmitRepository {
	return &SubmitRepository{q: q}
}

func (r *SubmitRepository) CreateSubmission(
	ctx context.Context,
	userID int32,
	data *CreateSubmitRequest,
) (*queries.Submission, error) {
	return r.q.CreateSubmission(ctx, queries.CreateSubmissionParams{
		SourceCode: pgtype.Text{
			Valid:  true,
			String: data.SourceCode,
		},
		Language: data.Language,
		UserID:   userID,
	})
}
