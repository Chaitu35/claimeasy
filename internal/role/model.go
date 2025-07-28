package role

import (
	"time"

	"github.com/chaitu35/claimeasy/internal/permission"
	"gorm.io/gorm"
)

type Role struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	Permissions []permission.Permission `gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}