package shared

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/Shoyeb45/server/pkg/apierr"
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

func HashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func ParseToken(raw string) (int32, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(
		raw,
		claims,
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(config.Cfg.JwtSecret), nil
		},
	)
	if err != nil || !token.Valid {
		return 0, apierr.NewUnauthorized("invalid or expired token")
	}

	userID, err := strconv.ParseInt(claims.Subject, 10, 32)
	if err != nil {
		return 0, apierr.NewUnauthorized("invalid token subject")
	}

	return int32(userID), nil
}

