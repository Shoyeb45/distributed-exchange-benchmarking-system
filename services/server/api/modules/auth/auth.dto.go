package auth

import "github.com/Shoyeb45/server/pkg/apiresponse"

type RequestLogIn struct {
	Email    string `json:"email"    validate:"required,email" doc:"User email"    example:"john@example.com"`
	Password string `json:"password" validate:"required,min=8"  doc:"User password" example:"secret123"`
}
type ResponseLogIn struct {
	ID          string `json:"id"          doc:"User UUID"      example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string `json:"name"        doc:"Full name"      example:"John Doe"`
	Email       string `json:"email"       doc:"User email"     example:"john@example.com"`
	AccessToken string `json:"accessToken" doc:"JWT token"      example:"eyJhbGci..."`
}


type LoginInput struct {
	Body RequestLogIn
}

type LoginOutput struct {
	Body apiresponse.Response[ResponseLogIn]
}