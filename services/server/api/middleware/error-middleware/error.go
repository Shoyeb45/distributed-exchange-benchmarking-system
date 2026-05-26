package errormiddleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/logger"
	"github.com/Shoyeb45/server/pkg/shared"
	
)

type errorResponse struct {
	Code    apierr.ErrorCode    `json:"code"`
	Message string              `json:"message"`
	Details []apierr.FieldError `json:"details,omitempty"`
}


func ErrorHandler(h shared.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := h(w, r)
		if err == nil {
			return
		}

		var apiErr *apierr.APIError

		switch {
		case apierr.As(err, &apiErr):
			logAPIError(r, apiErr)
			writeJSON(w, apiErr.StatusCode, errorResponse{
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
			writeJSON(w, http.StatusInternalServerError, errorResponse{
				Code:    apierr.CodeInternal,
				Message: "an unexpected error occurred",
			})
		}
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	var apiErr *apierr.APIError
	switch {
	case apierr.As(err, &apiErr):
		logAPIError(r, apiErr)
		writeJSON(w, apiErr.StatusCode, errorResponse{
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Details: apiErr.Details,
		})
	default:
		logger.Log.Error("unhandled error")
		writeJSON(w, http.StatusInternalServerError, errorResponse{
			Code:    apierr.CodeInternal,
			Message: "an unexpected error occurred",
		})
	}
}
