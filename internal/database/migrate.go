package database

import (
	"log"

	"github.com/chaitu35/claimeasy/internal/audit"
	"github.com/chaitu35/claimeasy/internal/auth"
	"github.com/chaitu35/claimeasy/internal/clinics"
	"github.com/chaitu35/claimeasy/internal/permission"
	"github.com/chaitu35/claimeasy/internal/role"
	"github.com/chaitu35/claimeasy/internal/user"
	"gorm.io/gorm"
)

// Migrate runs AutoMigrate on all models
// Migrate runs AutoMigrate on all models
func Migrate(db *gorm.DB) {
    models := []interface{}{
        &user.User{},
        &role.Role{},
        &permission.Permission{},
        &permission.RolePermission{},
        &permission.UserPermission{},
        &audit.AuditEntry{},
        // Auth
		&auth.LoginAttempt{},
		&auth.MFAChallenge{},
		&auth.MFAAttempt{},
		&auth.RecoveryCode{},
		&auth.PasswordReset{},
		// Clinics
		&clinics.Clinic{},
		&clinics.UserClinic{},
    }

    if err := db.AutoMigrate(models...); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    log.Println("✅ Database migration completed")
}

