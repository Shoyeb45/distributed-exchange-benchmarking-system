package auth

type GithubCallbackQuery struct {
    Code string `query:"code" validate:"required"`
}

type GithubCallback struct {
	ID          string `json:"id"          doc:"User UUID"  example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name"        doc:"Full name"  example:"John Doe"`
	Email       string `json:"email"       doc:"User email" example:"john@example.com"`
	AccessToken string `json:"accessToken" doc:"JWT token"  example:"eyJhbGci..."`
}
