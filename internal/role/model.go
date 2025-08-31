package role

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          uint           `gorm:"primaryKey"`
	Name        string         `gorm:"uniqueIndex;not null"`
	Description string         `gorm:"type:text"`
	IsSystem    bool           `gorm:"default:false"`
	Status      string         `gorm:"type:varchar(20);default:'active';index"` // active, inactive, deleted
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
