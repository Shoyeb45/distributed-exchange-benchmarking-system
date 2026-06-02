package submit

import (

	sqlcv1 "github.com/Shoyeb45/server/pkg/repository/gen-queries"
)

type CreateSubmitRequest struct {
	SourceCode string                   `json:"sourceCode" validate:"required"`
	Language   sqlcv1.SupportedLanguage `json:"language"   validate:"required"`
}

type SubmitResponse struct {
	ID int32 `json:"id"`
}

type KafkaMessage struct {
	UserID       int32 `json:"userId"`
	SubmissionID int32 `json:"submissionId"`
	CreatedTime  int64 `json:"createdTime"`
}