package errormiddleware

import (
	"log/slog"
	"net/http"

	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/logger"
	"github.com/Shoyeb45/server/pkg/shared"
)

// ErrorResponse is the shape of all error responses
// @Description Standard error response.
type ErrorResponse struct {
	Success bool                `json:"success"`
	Code    apierr.ErrorCode    `json:"code"`
	Message string              `json:"message"`
	Details []apierr.FieldError `json:"details,omitempty"`
}

// Middleware to pass the handler function and the error will be
// handled by this middleware.
func ErrorHandler(h shared.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		HandleError(w, r, err)
	}
}

func logAPIError(r *http.Request, err *apierr.APIError) {
	attrs := []any{
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("error_code", string(err.Code)),
		slog.String("error_message", err.Message),
	}

	if cause := err.Cause(); cause != nil {
		attrs = append(attrs, slog.Any("cause", cause))
	}

	switch {
	case err.StatusCode >= http.StatusInternalServerError:
		logger.Log.Error("server error", attrs...)
	case err.StatusCode == http.StatusUnauthorized || err.StatusCode == http.StatusNotFound:
		logger.Log.Info("client error", attrs...)
	default:
		logger.Log.Warn("client error", attrs...)
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *apierr.APIError

	switch {
	case apierr.As(err, &apiErr):
		logAPIError(r, apiErr)
		shared.WriteJSON(w, apiErr.StatusCode, ErrorResponse{
			Success: false,
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Details: apiErr.Details,
		})
	default:
		logger.Log.Error("unhandled error",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Any("error", err),
		)
		shared.WriteJSON(w, http.StatusInternalServerError, ErrorResponse{
			Success: false,
			Code:    apierr.CodeInternal,
			Message: "an unexpected error occurred",
		})
	}
}
