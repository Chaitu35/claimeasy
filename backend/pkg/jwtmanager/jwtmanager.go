package jwtmanager

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type PermissionClaim struct {
	View   bool `json:"view"`
	Edit   bool `json:"edit"`
	Delete bool `json:"delete"`
	Create bool `json:"create"`
	Print  bool `json:"print"`
	Export bool `json:"export"`
}

type UserClaim struct {
	jwt.RegisteredClaims
	UserId     int                              `json:"user_id"`
	Username   string                           `json:"username"`
	Role       string                           `json:"role"`
	Permissions map[string]PermissionClaim      `json:"permissions"`
}

type JWTManager struct {
	secretKey           string
	tokenDuration       time.Duration
	refreshTokenDuration time.Duration
	issuer              string
}

// Constructor for JWTManager
func NewJWTManager(secretKey string, tokenDuration, refreshTokenDuration time.Duration, issuer string) *JWTManager {
	return &JWTManager{
		secretKey:            secretKey,
		tokenDuration:        tokenDuration,
		refreshTokenDuration: refreshTokenDuration,
		issuer:               issuer,
	}
}

// Generate Access Token
func (jm *JWTManager) Generate(userId int, username, role string, permissions map[string]PermissionClaim) (string, error) {
	now := time.Now()
	claims := UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    jm.issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(jm.tokenDuration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		UserId:     userId,
		Username:   username,
		Role:       role,
		Permissions: permissions,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return token.SignedString([]byte(jm.secretKey))
}

// Generate Refresh Token
func (jm *JWTManager) GenerateRefreshToken(userId int) (string, time.Time, error) {
	exp := time.Now().Add(jm.refreshTokenDuration)
	claims := jwt.RegisteredClaims{
		Subject:   string(rune(userId)),
		ExpiresAt: jwt.NewNumericDate(exp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jm.secretKey))
	return signed, exp, err
}

// Verify and parse the JWT token
func (jm *JWTManager) Verify(tokenStr string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jm.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaim)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}


// VerifyRefreshToken checks if the refresh token is valid and returns the subject (user ID).
func (jm *JWTManager) VerifyRefreshToken(tokenStr string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jm.secretKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, jwt.ErrTokenInvalidClaims
	}

	// Extract and return user ID from the subject claim
	userID := claims.Subject
	// Convert string userID to int
	var uid int
	_, err = fmt.Sscanf(userID, "%d", &uid)
	if err != nil {
		return 0, errors.New("invalid subject in refresh token")
	}

	return uid, nil
}
