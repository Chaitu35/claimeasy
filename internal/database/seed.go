package database

import (
	"log"

	"github.com/chaitu35/claimeasy/internal/permission"
	"github.com/chaitu35/claimeasy/internal/role"
	"github.com/chaitu35/claimeasy/internal/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// --- Seed Roles ---
	roles := []role.Role{
		{Name: "admin", Status: "active"},
		{Name: "manager", Status: "active"},
		{Name: "coder", Status: "active"},
	}
	for _, r := range roles {
		if err := db.FirstOrCreate(&r, role.Role{Name: r.Name}).Error; err != nil {
			log.Fatalf("failed to seed role %s: %v", r.Name, err)
		}
	}

	// --- Seed Permissions ---
	pages := []string{
		"user", "permission", "role", "clinics", "clinician",
		"treatment", "diagnosis", "claim", "submission",
		"remittance", "report", "resubmit", "insurance_companies",
		"insurance_payers", "patient",
	}
	actions := []permission.Action{
		permission.ActionView,
		permission.ActionCreate,
		permission.ActionUpdate,
		permission.ActionDelete,
		permission.ActionExport,
		permission.ActionPrint,
	}

	var allPerms []permission.Permission
	for _, page := range pages {
		for _, act := range actions {
			key := page + "." + string(act)
			p := permission.Permission{
				Page:   page,
				Action: act,
				Key:    key,
			}
			if err := db.FirstOrCreate(&p, permission.Permission{Key: key}).Error; err != nil {
				log.Fatalf("failed to seed permission %s: %v", key, err)
			}
			allPerms = append(allPerms, p)
		}
	}

	// --- Link Roles and Permissions ---
	var admin role.Role
	db.First(&admin, "name = ?", "admin")

	var manager role.Role
	db.First(&manager, "name = ?", "manager")

	var coder role.Role
	db.First(&coder, "name = ?", "coder")

	for _, perm := range allPerms {
		// Admin = granted true
		db.FirstOrCreate(&permission.RolePermission{}, permission.RolePermission{
			RoleID:       admin.ID,
			PermissionID: perm.ID,
			Granted:      true,
		})

		// Manager = granted false
		db.FirstOrCreate(&permission.RolePermission{}, permission.RolePermission{
			RoleID:       manager.ID,
			PermissionID: perm.ID,
			Granted:      false,
		})

		// Coder = granted false
		db.FirstOrCreate(&permission.RolePermission{}, permission.RolePermission{
			RoleID:       coder.ID,
			PermissionID: perm.ID,
			Granted:      false,
		})
	}

	// --- Seed Admin User ---
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Admin@12345"), bcrypt.DefaultCost)
	adminUser := user.User{
		Username:     "Admin",
		Email:    "admin@example.com",
		PasswordHash: string(hashedPassword), // TODO: hash properly
		RoleID:   admin.ID,
	}
	if err := db.FirstOrCreate(&adminUser, user.User{Email: adminUser.Email}).Error; err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}

	log.Println("✅ Seeding completed")
}
