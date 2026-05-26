package auth

import (
	"net/http"

	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
	"github.com/Shoyeb45/server/pkg/apiresponse"
)

type AuthHandler struct {
	authRepository AuthRepository
}

func NewAuthHandler(repo AuthRepository) *AuthHandler {
	return &AuthHandler{authRepository: repo}
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body      RequestLogIn            true  "Login credentials"
// @Success      200   {object}  apiresponse.SwaggerResponse{data=ResponseLogIn}
// @Failure      400   {object}  errormiddleware.ErrorResponse
// @Failure      401   {object}  errormiddleware.ErrorResponse
// @Failure      422   {object}  errormiddleware.ErrorResponse
// @Router       /api/auth [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	login := validatormiddleware.From[RequestLogIn](r)

	// return writeJSON(w, http.StatusOK, user)
	return apiresponse.ResponseWriter(w, http.StatusCreated, "Logged in", ResponseLogIn{
		ID: "550e8400-e29b-41d4-a716-446655440000",
		Name: login.Email,
		Email: login.Email,
		AccessToken: "eyerfdfdj9dufjdfjd9fjpdfd",
	})
}
