package apiresponse

import (
	"net/http"

	"github.com/Shoyeb45/server/pkg/shared"
)

type Response[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data,omitempty"`
}

func ResponseWriter[T any](w http.ResponseWriter, status int, message string, data T) error {
	return shared.WriteJSON(w, status, Response[T]{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// SwaggerResponse is the envelope shape for all success responses.
// @Description Standard success response.
type SwaggerResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"operation successful"`
	Data    any    `json:"data"`
}
