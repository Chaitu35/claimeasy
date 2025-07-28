package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Port       string `env:"PORT" envDefault:"8080"`
	DB_DSN     string `env:"DB_DSN" envDefault:"postgres://user:password@localhost:5432/claimeasy?sslmode=disable"`
	JWT_SECRET string
	GO_ENV     string `env:"GO_ENV" envDefault:"development"`
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("GO_ENV")
	if env==""{
		env = "development"
	}
		_ = godotenv.Load(".env." + env)

		cfg:=&Config{
		Port:       os.Getenv("PORT"),
		DB_DSN:     os.Getenv("DB_DSN"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
		GO_ENV:     env,
		}
		if cfg.Port == "" {
			cfg.Port = "8080" // Default port
		}
return cfg, nil
}