package permission

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	// Permission CRUD
	CreatePermission(ctx context.Context, db *gorm.DB, p *Permission) error
	UpdatePermission(ctx context.Context, db *gorm.DB, p *Permission ) error
	GetPermissionByID(ctx context.Context, db *gorm.DB, id uint) (*Permission, error)
	GetPermissionByKey(ctx context.Context, db *gorm.DB, key string) (*Permission, error)
	ListPermissions(ctx context.Context, db *gorm.DB) ([]Permission, error)

	// Role ↔ Permission
	AssignPermissionToRole(ctx context.Context, db *gorm.DB, rp *RolePermission) error
	RevokePermissionFromRole(ctx context.Context, db *gorm.DB, roleID, permissionID uint) error
	ListRolePermissions(ctx context.Context, db *gorm.DB, roleID uint) ([]RolePermission, error)
	GetPermissionsByRoleID(ctx context.Context, db *gorm.DB, roleID uint) ([]Permission, error)

	// User ↔ Permission
	AssignPermissionToUser(ctx context.Context, db *gorm.DB, up *UserPermission) error
	RevokePermissionFromUser(ctx context.Context, db *gorm.DB, userID, permissionID uint) error
	ListUserPermissions(ctx context.Context, db *gorm.DB, userID uint) ([]UserPermission, error)

	//GetPermissionsByUserID(ctx context.Context, db *gorm.DB, userID uint) ([]Permission, error)
}

type permissionRepo struct{}

func NewRepository() Repository {
	return &permissionRepo{}
}

// ----------------------
// Permission CRUD
// ----------------------

func (r *permissionRepo) CreatePermission(ctx context.Context, db *gorm.DB, p *Permission) error {
	return db.WithContext(ctx).Create(p).Error
}

func (r *permissionRepo) UpdatePermission(ctx context.Context, db *gorm.DB, p *Permission) error {
	return db.WithContext(ctx).Save(p).Error
}
func (r *permissionRepo) GetPermissionByID(ctx context.Context, db *gorm.DB, id uint) (*Permission, error) {
	var p Permission
	err := db.WithContext(ctx).First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *permissionRepo) GetPermissionByKey(ctx context.Context, db *gorm.DB, key string) (*Permission, error) {
	var p Permission
	err := db.WithContext(ctx).Where("key = ?", key).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *permissionRepo) ListPermissions(ctx context.Context, db *gorm.DB) ([]Permission, error) {
	var list []Permission
	err := db.WithContext(ctx).Find(&list).Error
	return list, err
}

// ----------------------
// Role ↔ Permission
// ----------------------

func (r *permissionRepo) AssignPermissionToRole(ctx context.Context, db *gorm.DB, rp *RolePermission) error {
	return db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", rp.RoleID, rp.PermissionID).
		Assign(RolePermission{Granted: true, UpdatedAt: time.Now()}).
		FirstOrCreate(rp).Error
}

func (r *permissionRepo) RevokePermissionFromRole(ctx context.Context, db *gorm.DB, roleID, permissionID uint) error {
	return db.WithContext(ctx).
		Model(&RolePermission{}).
		Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		UpdateColumns(map[string]interface{}{
			"granted":    false,
			"updated_at": time.Now(),
		}).Error
}

func (r *permissionRepo) ListRolePermissions(ctx context.Context, db *gorm.DB, roleID uint) ([]RolePermission, error) {
	var list []RolePermission
	err := db.WithContext(ctx).
		Where("role_id = ? AND granted = ?", roleID, true).
		Preload("Permission").
		Find(&list).Error
	return list, err
}

func (r *permissionRepo) GetPermissionsByUserID(ctx context.Context, db *gorm.DB, userID uint) ([]Permission, error) {
    var perms []Permission
    err := db.WithContext(ctx).Model(&Permission{}).
        Joins("JOIN user_permissions up ON up.permission_id = permissions.id").
        Where("up.user_id = ?", userID).Find(&perms).Error
    return perms, err
}

// ----------------------
// User ↔ Permission
// ----------------------

func (r *permissionRepo) AssignPermissionToUser(ctx context.Context, db *gorm.DB, up *UserPermission) error {
	return db.WithContext(ctx).
		Where("user_id = ? AND permission_id = ?", up.UserID, up.PermissionID).
		Assign(UserPermission{Granted: true, UpdatedAt: time.Now()}).
		FirstOrCreate(up).Error
}

func (r *permissionRepo) RevokePermissionFromUser(ctx context.Context, db *gorm.DB, userID, permissionID uint) error {
	return db.WithContext(ctx).
		Model(&UserPermission{}).
		Where("user_id = ? AND permission_id = ?", userID, permissionID).
		UpdateColumns(map[string]interface{}{
			"granted":    false,
			"updated_at": time.Now(),
		}).Error
}

func (r *permissionRepo) ListUserPermissions(ctx context.Context, db *gorm.DB, userID uint) ([]UserPermission, error) {
	var list []UserPermission
	err := db.WithContext(ctx).
		Where("user_id = ? AND granted = ?", userID, true).
		Preload("Permission").
		Find(&list).Error
	return list, err
}
func (r *permissionRepo) GetPermissionsByRoleID(ctx context.Context, db *gorm.DB, roleID uint) ([]Permission, error) {
    var perms []Permission
    err := db.WithContext(ctx).
        Model(&RolePermission{}).
        Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
        Where("role_permissions.role_id = ? AND role_permissions.granted = ?", roleID, true).
        Select("permissions.*").
        Scan(&perms).Error
    return perms, err
}
