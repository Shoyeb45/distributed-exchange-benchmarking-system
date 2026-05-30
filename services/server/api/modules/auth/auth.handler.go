package auth

import (
	"fmt"
	"net/http"

	authmiddleware "github.com/Shoyeb45/server/api/middleware/auth"
	validatormiddleware "github.com/Shoyeb45/server/api/middleware/validator"
	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/config"
	"github.com/Shoyeb45/server/pkg/shared"
)

type AuthHandler struct {
	authRepository *AuthRepository
}

func NewAuthHandler(repo *AuthRepository) *AuthHandler {
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

	user, err := h.authRepository.FindOrCreateUser(r.Context(), *githubUser)
	if err != nil {
		return err
	}

	tokens, err := h.issueTokens(r, user.ID)
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

// RefreshTokens godoc
// @Summary Refresh access token
// @Description Exchange a valid refresh token for new tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} TokenResponse
// @Failure 401 {object} errormiddleware.ErrorResponse
// @Router /auth/refresh [post].
func (h *AuthHandler) RefreshTokens(w http.ResponseWriter, r *http.Request) error {
	body := validatormiddleware.From[RefreshTokenBody](r)

	userID, err := shared.ParseToken(body.RefreshToken)
	if err != nil {
		return err
	}

	storedUserID, err := h.authRepository.ValidateRefreshToken(r.Context(), body.RefreshToken)
	if err != nil {
		return err
	}

	if storedUserID != userID {
		return apierr.NewUnauthorized("refresh token mismatch")
	}

	if err := h.authRepository.RevokeRefreshToken(r.Context(), body.RefreshToken); err != nil {
		return err
	}

	tokens, err := h.issueTokens(r, userID)
	if err != nil {
		return err
	}

	return shared.WriteJSON(w, http.StatusOK, TokenResponse{
		Success:      true,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// Me godoc
// @Summary Get current user
// @Description Returns the authenticated user profile
// @Tags Auth
// @Produce json
// @Success 200 {object} MeResponse
// @Failure 401 {object} errormiddleware.ErrorResponse
// @Router /auth/me [get].
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) error {
	userID, ok := authmiddleware.UserID(r.Context())
	if !ok {
		return apierr.NewUnauthorized("authentication required")
	}

	user, err := h.authRepository.GetUserByID(r.Context(), userID)
	if err != nil {
		return err
	}

	return shared.WriteJSON(w, http.StatusOK, MeResponse{
		Success:        true,
		ID:             user.ID,
		Name:           user.Name,
		Email:          user.Email,
		AvatarUrl:      user.AvatarUrl,
		GithubUsername: user.GithubUsername,
	})
}

func (h *AuthHandler) issueTokens(r *http.Request, userID int32) (*shared.Tokens, error) {
	tokens, err := shared.GenerateTokens(userID)
	if err != nil {
		return nil, err
	}

	if err := h.authRepository.SaveRefreshToken(r.Context(), userID, tokens.RefreshToken); err != nil {
		return nil, err
	}

	return tokens, nil
}
