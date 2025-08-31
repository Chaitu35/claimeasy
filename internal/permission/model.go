package permission

import (
	"time"

	"gorm.io/gorm"
)

type Action string

const (
	ActionView   Action = "view"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
	ActionExport Action = "export"
	ActionPrint  Action = "print"
)

// ----------------------
// Models
// ----------------------

// Base Permission (e.g. "user.view")
type Permission struct {
	ID          uint           `gorm:"primaryKey"`
	Page        string         `gorm:"index;not null"`
	Action      Action         `gorm:"type:varchar(20);not null"`
	Key         string         `gorm:"uniqueIndex;not null"` // e.g., "user.view"
	   // Added: active, inactive
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RolePermission struct {
	ID           uint           `gorm:"primaryKey"`
	RoleID       uint           `gorm:"not null;index:idx_role_permission,unique"`
	PermissionID uint           `gorm:"not null;index:idx_role_permission,unique"`
	Granted      bool           `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relations (no direct import, avoid cycle)
	Permission Permission `gorm:"foreignKey:PermissionID"`
}

type UserPermission struct {
	ID           uint           `gorm:"primaryKey"`
	UserID       uint           `gorm:"not null;index:idx_user_permission,unique"`
	PermissionID uint           `gorm:"not null;index:idx_user_permission,unique"`
	Granted      bool           `gorm:"default:true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relations
	Permission Permission `gorm:"foreignKey:PermissionID"`
}
