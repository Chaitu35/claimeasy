package config

import (
	"log"
	"os"
)

type MFAConfig struct {
	SecretKey string
}

func LoadMFAConfig() MFAConfig {
	key := os.Getenv("MFA_SECRET_KEY")
	if key == "" {
		log.Fatal("MFA_SECRET_KEY not set")
	}

	if len(key) != 32 {
		log.Fatal("MFA_SECRET_KEY must be 32 bytes long (for AES-256)")
	}

	return MFAConfig{
		SecretKey: key,
	}
}
