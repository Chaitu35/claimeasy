package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
)
func GenerateRandomString(length int) (string, error) {
    bytes := make([]byte, length)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
// GenerateResetToken generates a random reset token (plaintext + hash)
func GenerateResetToken() (string, string, error) {
	// 32 bytes random token
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", "", err
	}

	// plaintext token (send via email/SMS)
	plainToken := hex.EncodeToString(b)

	// hash token (store in DB)
	hash := sha256.Sum256([]byte(plainToken))
	hashedToken := hex.EncodeToString(hash[:])

	return plainToken, hashedToken, nil
}

// VerifyResetToken verifies plaintext token against stored hash
func VerifyResetToken(plainToken, storedHash string) bool {
	hash := sha256.Sum256([]byte(plainToken))
	return hex.EncodeToString(hash[:]) == storedHash
}
