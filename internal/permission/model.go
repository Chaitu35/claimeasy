package permission

import (
	"time"

	"gorm.io/gorm"
)

type Permission struct {
	ID          uint   `gorm:"primaryKey"`
	Page        string `gorm:"index;not null"`       // e.g., "user", "claims", "dashboard"
	Code        string `gorm:"uniqueIndex;not null"` // e.g., "view", "edit", "submit"
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
