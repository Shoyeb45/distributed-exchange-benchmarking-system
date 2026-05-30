package auth

import (
	"fmt"
	"net/http"

	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/config"
	"github.com/Shoyeb45/server/pkg/shared"
)

type AuthHandler struct {
	authRepository AuthRepository
}

func NewAuthHandler(repo AuthRepository) *AuthHandler {
	return &AuthHandler{authRepository: repo}
}

// OAuth github redirect godoc
// @Summary      RedirectGithub
// @Description  Redirect to github oauth page
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      307 
// @Failure      500   {object}  errormiddleware.ErrorResponse
// @Router       /auth/github [get].
func (h *AuthHandler) RedirectGithub(w http.ResponseWriter, r *http.Request) error {
	redirectURI := fmt.Sprintf(
        "%s?client_id=%s&scope=user:email",
        config.Cfg.GithubRedirectURL,
        config.Cfg.GithubClientID,
    )

	http.Redirect(w, r, redirectURI, http.StatusFound)
	return nil
}

// OAuth Github Verify
// @Summary Verify github code
// @Description Verify github code
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200
// @Failure 500
// @Router /auth/github/callback [get].
func (h *AuthHandler) GithubCallback(w http.ResponseWriter, r *http.Request) error {
	query := validatormiddleware.From[GithubCallbackQuery](r)

	accessToken, err := shared.GetGithubAccessToken(query.Code)
	if err != nil {
		return err
	}

	githubUser, err := shared.GetGithubUser(accessToken)
	if err != nil {
		return err
	}

	_, err = h.authRepository.GetUserByGithubId(r.Context(), int32(githubUser.ID))
	if err == nil {
		return apierr.NewConflict("user already exists")
	}

	user, err := h.authRepository.CreateUser(r.Context(), *githubUser)

	if err != nil {
		return err
	}

	tokens, err := shared.GenerateTokens(user.ID)

	if err != nil {
		return err
	}

	http.Redirect(
		w,
		r,
		fmt.Sprintf(
			"%s/auth/callback?accessToken=%s&refreshToken=%s",
			config.Cfg.OriginURL,
			tokens.AccessToken,
			tokens.RefreshToken,
		),
		http.StatusFound,
	)

	return nil
}
