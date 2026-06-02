package submit

import (
	"encoding/base64"
	"net/http"
	"time"

	authmiddleware "github.com/Shoyeb45/server/api/middleware/auth"
	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/apiresponse"
	kafkaservice "github.com/Shoyeb45/server/pkg/core"
)

type SubmitHandler struct {
	submitRepository *SubmitRepository
}

func NewSubmitHandler(repo *SubmitRepository) *SubmitHandler {
	return &SubmitHandler{submitRepository: repo}
}

// Create submission godoc
// @Summary      Create Submission
// @Tags         Submit
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        body  body      CreateSubmitRequest   true  "source code and lang"
// @Success      200   {object}  apiresponse.SwaggerResponse{data=SubmitResponse}
// @Failure      500   {object}  errormiddleware.ErrorResponse
// @Router       /submit/ [post].
func (h *SubmitHandler) CreateSubmission(res http.ResponseWriter, req *http.Request) error {
	ctx := req.Context()
	body := validatormiddleware.From[CreateSubmitRequest](req)

	userID, ok := authmiddleware.UserID(ctx)

	if !ok {
		return apierr.NewUnauthorized("authentication required")
	}
	body.SourceCode = base64.StdEncoding.EncodeToString([]byte(body.SourceCode))

	createdSubmission, err := h.submitRepository.CreateSubmission(ctx, userID, &body)

	if err != nil {
		return err
	}

	kafkaservice.Produce(userID, KafkaMessage{
		UserID:       userID,
		SubmissionID: createdSubmission.ID,
		CreatedTime:  time.Now().Unix(),
	})

	return apiresponse.ResponseWriter(res, http.StatusCreated, "submission created", SubmitResponse{
		ID: createdSubmission.ID,
	})
}
