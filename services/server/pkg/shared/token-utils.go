package shared

import (
	"strconv"
	"time"

	"github.com/Shoyeb45/server/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

func generateToken(userId int32, expireDuration time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(
			time.Now().Add(expireDuration),
		),
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Issuer:   config.Cfg.TokenIssuer,
		Audience: jwt.ClaimStrings{
			config.Cfg.TokenAudience,
		},
		Subject: strconv.Itoa(int(userId)),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(config.Cfg.JwtSecret))
}

func GenerateTokens(userId int32) (*Tokens, error) {
	accessToken, err := generateToken(userId, time.Duration(config.Cfg.AccessTokenValiditySec)*time.Second)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(userId, time.Duration(config.Cfg.RefreshTokenValiditySec)*time.Second)
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
