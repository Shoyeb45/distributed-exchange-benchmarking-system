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
	JwtSecret               string
	TokenIssuer             string
	TokenAudience           string
	AccessTokenValiditySec  int
	RefreshTokenValiditySec int
	// Directory where logging should happen
	LogDirectory   string
	TimeoutSeconds string
	// kafka cred
	KafkaAddress string
	KafkaTopic   string
}

var Cfg *Config

func LoadEnvironmentVariables() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}

	accessTokenValidity, err := strconv.Atoi(getEnv("ACCESS_TOKEN_VALIDITY_SEC", "3600"))
	if err != nil {
		return err
	}
	refreshTokenValidity, err := strconv.Atoi(getEnv("REFRESH_TOKEN_VALIDITY_SEC", "864000"))
	if err != nil {
		return err
	}

	Cfg = &Config{
		DatabaseURL:             getEnv("DATABASE_URL", ""),
		Port:                    getEnv("PORT", "8080"),
		GithubClientID:          getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret:      getEnv("GITHUB_CLIENT_SECRET", ""),
		OriginURL:               getEnv("ORIGIN_URL", "http://localhost:3000"),
		GithubRedirectURL:       getEnv("GITHUB_REDIRECT_URL", "https://github.com/login/oauth/authorize"),
		JwtSecret:               getEnv("JWT_SECRET", ""),
		TokenIssuer:             getEnv("TOKEN_ISSUER", ""),
		TokenAudience:           getEnv("TOKEN_AUDIENCE", ""),
		AccessTokenValiditySec:  accessTokenValidity,
		RefreshTokenValiditySec: refreshTokenValidity,
		Stage:                   getEnv("STAGE", "dev"),
		LogDirectory:            getEnv("LOG_DIRECTORY", "logs"),
		TimeoutSeconds:          getEnv("TIMEOUT_SECONDS", "12"),
		KafkaAddress:            getEnv("KAFKA_ADDRESS", "localhost:9092"),
		KafkaTopic: 			 getEnv("KAFKA_TOPIC", "submission-created"),		
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

	if Cfg.Stage == "" || (Cfg.Stage != "dev" && Cfg.Stage != "prod") {
		return errors.New("STAGE must be either dev or prod")
	}

	if Cfg.JwtSecret == "" {
		return errors.New("JWT_SECRET must be provided")
	}

	if Cfg.KafkaAddress == "" {
		return errors.New("KAFKA_ADDRESS must be provided")
	}

	if Cfg.KafkaTopic == "" {
		return errors.New("KAFKA_TOPIC must be provided")
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

	if !strings.HasPrefix(dbURL, "postgresql://") && !strings.HasPrefix(dbURL, "postgres://") {
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
