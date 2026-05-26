package config

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	// Stage or environemnt of the application. dev or prod
	Stage              string
	DatabaseURL        string
	Port               string
	GithubClientID     string
	GithubClientSecret string
	// Frontend application url
	OriginURL               string
	GithubRedirectURL       string
	JwtPrivateKey           string
	JwtPublicKey            string
	TokenIssuer             string
	TokenAudience           string
	AccessTokenValiditySec  string
	RefreshTokenValiditySec string
	// Directory where logging should happen
	LogDirectory   string
	TimeoutSeconds string
}

var Cfg *Config

func LoadEnvironmentVariables() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	Cfg = &Config{
		DatabaseURL:             getEnv("DATABASE_URL", ""),
		Port:                    getEnv("PORT", "8080"),
		GithubClientID:          getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret:      getEnv("GITHUB_CLIENT_SECRET", ""),
		OriginURL:               getEnv("ORIGIN_URL", "http://localhost:3000"),
		GithubRedirectURL:       getEnv("GIHUB_REDIRECT_URL", "https://github.com/login/oauth/authorize"),
		JwtPrivateKey:           getEnv("JWT_PRIVATE_KEY", ""),
		JwtPublicKey:            getEnv("JWT_PUBLIC_KEY", ""),
		TokenIssuer:             getEnv("TOKEN_ISSUER", ""),
		TokenAudience:           getEnv("TOKEN_AUDIENCE", ""),
		AccessTokenValiditySec:  getEnv("ACCESS_TOKEN_VALIDITY_SEC", "3600"),
		RefreshTokenValiditySec: getEnv("REFRESH_TOKEN_VALIDITY_SEC", "864000"),
		Stage:                   getEnv("STAGE", "dev"),
		LogDirectory:            getEnv("LOG_DIRECTORY", "logs"),
		TimeoutSeconds:          getEnv("TIMEOUT_SECONDS", "12"),
	}

	return validateEnvironmentVariables()
}

func validateEnvironmentVariables() error {
	var err error

	if err = validatePort(); err != nil {
		return err
	}

	if err = validateDatabaseURL(); err != nil {
		return err
	}

	if Cfg.GithubClientID == "" {
		return errors.New("github Client ID should not be empty")
	}

	if Cfg.GithubClientSecret == "" {
		return errors.New("github Client Secret should not be empty")
	}

	if Cfg.OriginURL == "" ||
		(!strings.HasPrefix(Cfg.OriginURL, "http://") && !strings.HasPrefix(Cfg.OriginURL, "https://")) {
		return errors.New("origin URL must be a given for CORS and it must start with http:// or https://")
	}

	if Cfg.GithubRedirectURL == "" {
		return errors.New("github Redirect URL must be provided")
	}

	if Cfg.JwtPrivateKey == "" || Cfg.JwtPublicKey == "" {
		return errors.New("JWT Keys must be provided for secure authentication")
	}

	if Cfg.Stage == "" || (Cfg.Stage != "dev" && Cfg.Stage != "prod") {
		return errors.New("STAGE must be either dev or prod")
	}

	_, er := strconv.Atoi(Cfg.TimeoutSeconds)
	if er != nil {
		return er
	}

	return err
}

func validateDatabaseURL() error {
	dbURL := Cfg.DatabaseURL

	if dbURL == "" {
		return errors.New("database url cannot be empty")
	}

	if !strings.HasPrefix(dbURL, "postgresql://") && !strings.HasPrefix(dbURL, "potgres://") {
		return errors.New("you must provide postgres database url which starts with postgresql:// or postgres://")
	}
	return nil
}

func validatePort() error {
	port, err := strconv.Atoi(Cfg.Port)

	if err != nil {
		return errors.New("PORT must be a number, PORT: " + Cfg.Port)
	}
	if port < 0 || port > 65535 {
		return errors.New("PORT out of valid range, expected between min : 1, max: 65535, got " + Cfg.Port)
	}
	return nil
}

func getEnv(key string, defaultVal string) string {
	val, ok := os.LookupEnv(key)

	if ok {
		return val
	}
	if val == "" {
		return defaultVal
	}

	return defaultVal
}
