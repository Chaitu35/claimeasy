package clinics

import (
	"time"

	"github.com/chaitu35/claimeasy/internal/user"
	"gorm.io/gorm"
)

type ControlBody string

const (
	DHA  ControlBody = "DHA"
	HAAD ControlBody = "HAAD"
	MOH  ControlBody = "MOH"
)

// Main Clinic model
type Clinic struct {
	ID             uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
    ControlBody string `gorm:"type:varchar(10);not null;check:control_body IN ('DHA','HAAD','MOH')" json:"control_body"`
	Phone          string         `gorm:"type:varchar(20);not null" json:"phone"`
	Email          string         `gorm:"type:varchar(100);not null" json:"email"`
	FacilityID     string         `gorm:"type:varchar(100);" json:"facility_id"`
	Licence        string         `gorm:"type:varchar(100);" json:"licence"`
	Representative string         `gorm:"type:varchar(100);" json:"representative"`
	CreatedAt      time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
	UserID *uint     `json:"user_id"`
    User   *user.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user,omitempty"`

}

// Join table for many-to-many: User ↔ Clinic
type UserClinic struct {
    ID        uint           `gorm:"primaryKey"`
    UserID    uint           `gorm:"not null;index:idx_user_clinic,unique"`
    ClinicID  uint           `gorm:"not null;index:idx_user_clinic,unique"`
    AssignedAt time.Time     `gorm:"autoCreateTime"`
    RevokedAt gorm.DeletedAt `gorm:"index"`
}
