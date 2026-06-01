package submit

import sqlcv1 "github.com/Shoyeb45/server/pkg/repository/gen-queries"

type CreateSubmitRequest struct {
	SourceCode string                   `json:"sourceCode" validate:"required"`
	Language   sqlcv1.SupportedLanguage `json:"language"   validate:"required"`
}

type SubmitResponse struct {
	ID int32 `json:"id"`
}
