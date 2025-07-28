package audit

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Auditlog struct {
	id 	 uint `gorm:"primaryKey"`
	Entity string `gorm:"not null" json:"entity"`
	EntityID uint `gorm:"not null" json:"entity_id"`
	Action string `gorm:"not null"  json:"action"`
	ChangedBy uint `gorm:"not null" json:"change_by"`
	Changes  datatypes.JSON  `json:"changes"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}