package security

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService defines contract for issuing and validating tokens
type JWTService interface {
	GenerateAccessToken(userID uint, email string, permissions []string) (string, error)
	GenerateRefreshToken(userID uint, email string) (string, error)
	GenerateResetToken(userID uint, ttl time.Duration) (string, error)

	ValidateAccessToken(token string) (*CustomClaims, error)
	ValidateRefreshToken(token string) (*CustomClaims, error)
	ValidateResetToken(token string) (*CustomClaims, error)
}

// Implementation
type jwtService struct {
	secretKey     string
	issuer        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// Claims structure
type CustomClaims struct {
	UserID      uint     `json:"user_id"`
	Email       string   `json:"email,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

// Constructor
func NewJWTService(secretKey, issuer string, accessExpiry, refreshExpiry time.Duration) JWTService {
	return &jwtService{secretKey, issuer, accessExpiry, refreshExpiry}
}

// ------------------------
// Generate Tokens
// ------------------------
func (j *jwtService) GenerateAccessToken(userID uint, email string, permissions []string) (string, error) {
	claims := CustomClaims{
		UserID:      userID,
		Email:       email,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secretKey))
}

func (j *jwtService) GenerateRefreshToken(userID uint, email string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secretKey))
}

func (j *jwtService) GenerateResetToken(userID uint, ttl time.Duration) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(j.secretKey))
}

// ------------------------
// Validate Tokens
// ------------------------
func (j *jwtService) ValidateAccessToken(tokenStr string) (*CustomClaims, error) {
	return j.parseToken(tokenStr)
}

func (j *jwtService) ValidateRefreshToken(tokenStr string) (*CustomClaims, error) {
	return j.parseToken(tokenStr)
}

func (j *jwtService) ValidateResetToken(tokenStr string) (*CustomClaims, error) {
	return j.parseToken(tokenStr)
}

// ------------------------
// Internal helper
// ------------------------
func (j *jwtService) parseToken(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
