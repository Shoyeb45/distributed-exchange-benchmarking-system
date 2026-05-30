package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Shoyeb45/server/pkg/apierr"
	"github.com/Shoyeb45/server/pkg/config"
	"github.com/Shoyeb45/server/pkg/logger"
)

type TokenResp struct {
	AccessToken      string `json:"access_token"`
	Scope            string `json:"scope"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}
type GithubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type GithubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

var client = &http.Client{}

func GetGithubAccessToken(code string) (string, error) {
	payload := map[string]string{
		"client_id":     config.Cfg.GithubClientID,
		"client_secret": config.Cfg.GithubClientSecret,
		"code":          code,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github returned status %d", resp.StatusCode)
	}

	var tokenResp TokenResp

	if err = json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	// GitHub returns 200 with an error field on bad codes
	if tokenResp.AccessToken == "" {
		logger.Log.Info("", slog.Any("tokenresp", tokenResp))
		return "", apierr.NewUnauthorized("github returned empty access token — code may be expired or already used")
	}

	return tokenResp.AccessToken, nil
}

func GetGithubUser(accessToken string) (*GithubUser, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	logger.Log.Info("access token", slog.String("access token", accessToken))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, apierr.NewUnauthorized(fmt.Sprintf("github returned status %d", resp.StatusCode))
	}

	var user GithubUser

	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	// GitHub may return null email if private
	if user.Email == "" {
		email, err := getGithubPrimaryEmail(client, accessToken)
		if err != nil {
			return nil, err
		}

		user.Email = email
	}

	if user.Email == "" {
		return nil, apierr.NewUnauthorized("no verified github email found")
	}

	return &user, nil
}

func getGithubPrimaryEmail(
	client *http.Client,
	accessToken string,
) (string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github returned status %d", resp.StatusCode)
	}

	var emails []GithubEmail

	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	// primary verified email first
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	// fallback: any verified email
	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	return "", nil
}
