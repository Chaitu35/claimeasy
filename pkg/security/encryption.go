package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/base64"
    "errors"
    "io"
)

// Encrypt encrypts plaintext using AES-GCM with a random nonce
func Encrypt(secretKey, plaintext string) (string, error) {
    if len(secretKey) != 32 {
        return "", errors.New("secret key must be exactly 32 bytes for AES-256")
    }
    
    block, err := aes.NewCipher([]byte(secretKey))
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    nonce := make([]byte, aesGCM.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM
func Decrypt(secretKey, ciphertext string) (string, error) {
    if len(secretKey) != 32 {
        return "", errors.New("secret key must be exactly 32 bytes for AES-256")
    }
    
    block, err := aes.NewCipher([]byte(secretKey))
    if err != nil {
        return "", err
    }

    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    raw, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    nonceSize := aesGCM.NonceSize()
    if len(raw) < nonceSize {
        return "", errors.New("ciphertext too short")
    }

    nonce, encBytes := raw[:nonceSize], raw[nonceSize:]
    plaintext, err := aesGCM.Open(nil, nonce, encBytes, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}

// Note: HashString, CheckStringHash, GenerateNumericCode, and GenerateRecoveryCodes 
// are already implemented in your password.go file - use those instead!