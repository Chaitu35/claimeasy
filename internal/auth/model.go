package auth

import (
	"time"

	"gorm.io/gorm"
)
type LoginAttempt struct {
    ID         uint      `gorm:"primaryKey"`
    Identifier string    `gorm:"index;not null"`
    Success    bool      `gorm:"default:false"`
    CreatedAt  time.Time `gorm:"autoCreateTime"`
    isLocked   bool      `gorm:"-"`
}

// -------------------------
// MFA Challenge (per OTP)
// -------------------------
// For SMS/Email OTP (short-lived, per request)
// Each challenge is a row
type MFAChallenge struct {
	ID             uint           `gorm:"primaryKey"`
	UserID         uint           `gorm:"not null;index"`
	Channel        string         `gorm:"type:varchar(20);not null"` // "totp", "sms", "email"
	CodeHash       string         `gorm:"not null"`                  // hashed OTP, never raw
	FailedAttempts int            `gorm:"default:0"`                 // attempts for this OTP
	MaxAttempts    int            `gorm:"default:5"`
	SecretEnc      string         `gorm:""`                          // optional: if channel needs secret for validation
	ExpiresAt      time.Time      `gorm:"index"`
	Consumed       bool           `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

// -------------------------
// MFA Attempt (per user/session)
// -------------------------
// For long-lived rate limiting (mainly TOTP)
// One row per user + channel
type MFAAttempt struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null;index"`
	Channel   string         `gorm:"type:varchar(20);not null"` // "totp", "sms", "email"
	Count     int            `gorm:"default:0"`
	Max       int            `gorm:"default:5"`
	LockedAt  *time.Time     // when locked due to too many attempts
	ResetAt   *time.Time     // when counter resets
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}


// --------------------
// Recovery Codes Model
// --------------------
// Permanent backup codes created during MFA enrollment
// Each can be used only once
type RecoveryCode struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null;index"`
	CodeHash  string         `gorm:"not null"`   // bcrypt or sha256 hash of recovery code
	Used      bool           `gorm:"default:false"`
	UsedAt    *time.Time
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type PasswordReset struct {
    ID        uint      `gorm:"primaryKey"`
    UserID    uint      `gorm:"index;not null"`
    Token     string    `gorm:"uniqueIndex;not null"`
    ExpiresAt time.Time `gorm:"not null"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
}
