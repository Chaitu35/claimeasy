package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DB_DSN     string
	JWT_SECRET string
	GO_ENV     string
}

func LoadConfig() *Config {
	// 1. Detect environment
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	// 2. Load corresponding .env file
	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Printf("Warning: no .env.%s file found, relying on OS env variables", env)
	}

	// 3. Build config from env
	cfg := &Config{
		Port:       os.Getenv("PORT"),
		DB_DSN:     os.Getenv("DB_DSN"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		GO_ENV:     env,
	}

	// 4. Defaults
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DB_DSN == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	return cfg
}
