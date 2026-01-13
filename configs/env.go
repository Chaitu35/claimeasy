package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all environment variables
type Config struct {
	DbDSN                 string
	JWTSecret             string
	AccessTokenExpiryMin  int
	RefreshTokenExpiryDays int
	AppSecret             string
	GoEnv                 string
	Port                  string
	SendGridAPIKey        string
	SendGridFrom          string
	FrontendURL           string
	GormLogLevel          string
}

// LoadConfig loads .env file based on GO_ENV and returns Config
func LoadConfig() *Config {
	// Determine which env file to load
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	envFile := ".env." + env
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("Warning: could not load %s, using system environment variables", envFile)
	}

	cfg := &Config{
		DbDSN:                 os.Getenv("Db_DSN"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		AppSecret:             os.Getenv("APP_SECRET"),
		GoEnv:                 env,
		Port:                  os.Getenv("PORT"),
		SendGridAPIKey:        os.Getenv("SENDGRID_API_KEY"),
		SendGridFrom:          os.Getenv("SENDGRID_FROM"),
		FrontendURL:           os.Getenv("FRONTEND_URL"),
		GormLogLevel:          os.Getenv("GORM_LOG_LEVEL"),
	}

	// Parse optional integer configs
	if v := os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES"); v != "" {
		if mins, err := strconv.Atoi(v); err == nil {
			cfg.AccessTokenExpiryMin = mins
		} else {
			cfg.AccessTokenExpiryMin = 15
		}
	} else {
		cfg.AccessTokenExpiryMin = 15
	}

	if v := os.Getenv("REFRESH_TOKEN_EXPIRY_DAYS"); v != "" {
		if days, err := strconv.Atoi(v); err == nil {
			cfg.RefreshTokenExpiryDays = days
		} else {
			cfg.RefreshTokenExpiryDays = 7
		}
	} else {
		cfg.RefreshTokenExpiryDays = 7
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}
