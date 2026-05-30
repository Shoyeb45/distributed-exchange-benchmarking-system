package auth

type GithubCallbackQuery struct {
	Code string `query:"code" validate:"required"`
}

type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type TokenResponse struct {
	Success      bool   `json:"success"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type MeResponse struct {
	Success        bool   `json:"success"`
	ID             int32  `json:"id"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	AvatarUrl      string `json:"avatarUrl"`
	GithubUsername string `json:"githubUsername"`
}

type GithubCallback struct {
	ID          string `json:"id"          doc:"User UUID"  example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name"        doc:"Full name"  example:"John Doe"`
	Email       string `json:"email"       doc:"User email" example:"john@example.com"`
	AccessToken string `json:"accessToken" doc:"JWT token"  example:"eyJhbGci..."`
}
