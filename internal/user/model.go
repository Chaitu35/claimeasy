package user

import (
	"time"

	"gorm.io/gorm"
	//"github.com/chaitu35/claimeasy/internal/permission"
	"github.com/chaitu35/claimeasy/internal/role"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`

	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Phone     string `gorm:"uniqueIndex;not null"`
	Address   string `gorm:"not null"`
    ProfilePicture string `gorm:"not null"`
    NationalID string `gorm:"uniqueIndex;not null"`

	TwoFactorSecret string `gorm:"not null"` // Secret for two-factor authentication
	need2FAuth bool `gorm:"default:false"` // Two-factor authentication enabled

	isActive bool   `gorm:"default:true"`
	isValidated bool `gorm:"default:false"`

	RoleID    uint           `gorm:"index"`            // foreign key
	Role      role.Role      `gorm:"foreignKey:RoleID"`// relation


	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time  `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
