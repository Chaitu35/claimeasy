package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type AuthConfig struct {
	JWTSecret         []byte
	AccessTokenExpiry time.Duration
	SigningMethod     string
}

func LoadAuthConfig() AuthConfig {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET not set")
	}

	expiryMinutes, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRY_MINUTES"))
	if err != nil || expiryMinutes <= 0 {
		expiryMinutes = 15 // default
	}

	method := os.Getenv("JWT_SIGNING_METHOD")
	if method == "" {
		method = "HS256"
	}

	return AuthConfig{
		JWTSecret:         []byte(secret),
		AccessTokenExpiry: time.Duration(expiryMinutes) * time.Minute,
		SigningMethod:     method,
	}
}
