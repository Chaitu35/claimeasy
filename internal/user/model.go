package user

import (
	"time"

	"github.com/chaitu35/claimeasy/internal/role"
	"gorm.io/gorm"
)

// User represents the database model for users.
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string         `gorm:"not null;size:255" json:"-"` // hidden from JSON

	// Personal Information
	FirstName      string `gorm:"not null;size:100" json:"first_name"`
	LastName       string `gorm:"not null;size:100" json:"last_name"`
	Phone          string `gorm:"uniqueIndex;not null;size:20" json:"phone"`
	Address        string `gorm:"type:text" json:"address,omitempty"`
	ProfilePicture string `gorm:"size:500" json:"profile_picture,omitempty"`
	NationalID     string `gorm:"uniqueIndex;not null;size:50" json:"-"` // sensitive → never exposed in API
    ExpiryDate     time.Time `json:"expiry_date,omitempty"`
	// Two-Factor Authentication (MFA)
	TwoFactorSecret string   `gorm:"size:255" json:"-"`
	Need2FAuth      bool     `gorm:"default:false;column:need_2f_auth" json:"need_2f_auth"`
	TwoFAMethod TwoFactorType `gorm:"size:50;column:two_fa_method" json:"two_fa_method"`
// store multiple methods (authenticator, sms, email, etc.)

	// Account Status
	IsActive    bool `gorm:"default:true;column:is_active" json:"is_active"`
	IsValidated bool `gorm:"default:false;column:is_validated" json:"is_validated"`

	// Role Relationship
	RoleID uint      `gorm:"index;not null" json:"role_id"`
	Role   role.Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"role,omitempty"`

	// Timestamps
	CreatedAt time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


// Lifecycle action type
type UserStatusAction string

const (
	UserActionActivate   UserStatusAction = "activate"
	UserActionDeactivate UserStatusAction = "deactivate"
	UserActionDelete     UserStatusAction = "delete"
	UserActionRestore    UserStatusAction = "restore"
	// future: suspend, lock, archive, etc.
)

type PasswordResetToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index;not null"` // link to user
	TokenHash string    `gorm:"size:256;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

