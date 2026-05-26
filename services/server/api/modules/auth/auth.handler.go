package auth

import (
	"net/http"

	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
)

type AuthHandler struct {
	authRepository AuthRepository
}

func NewAuthHandler(repo AuthRepository) *AuthHandler {
	return &AuthHandler{authRepository: repo}
}

// Login godoc
// @Summary      Login
// @Description  Authenticate user and return token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      LogIn    true  "Login credentials"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      401   {object}  map[string]string
// @Router       /api/auth/ [post]   ← must match exactly what chi mounts
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) error {
	login := validatormiddleware.From[LogIn](r)

	// return writeJSON(w, http.StatusOK, user)
	w.Write([]byte("Hi " + login.Name));
	return nil;
}
