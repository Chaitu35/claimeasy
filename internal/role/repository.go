package role

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	CreateRole(ctx context.Context, db *gorm.DB, role *Role) error
	UpdateRole(ctx context.Context, db *gorm.DB, role *Role) error
	GetRoleByID(ctx context.Context, db *gorm.DB, id uint) (*Role, error)
	ListRoles(ctx context.Context, db *gorm.DB, req ListRolesRequest) ([]Role, int64, error)
	UpdateRoleStatus(ctx context.Context, db *gorm.DB, id uint, status string) error
	HasActiveUsers(ctx context.Context, db *gorm.DB, roleID uint) (bool, error)
}

type roleRepo struct{}

func NewRepository() Repository {
	return &roleRepo{}
}

func (r *roleRepo) CreateRole(ctx context.Context, db *gorm.DB, role *Role) error {
	return db.WithContext(ctx).Create(role).Error
}

func (r *roleRepo) UpdateRole(ctx context.Context, db *gorm.DB, role *Role) error {
	return db.WithContext(ctx).Save(role).Error
}

func (r *roleRepo) GetRoleByID(ctx context.Context, db *gorm.DB, id uint) (*Role, error) {
	var role Role
	if err := db.WithContext(ctx).First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepo) ListRoles(ctx context.Context, db *gorm.DB, req ListRolesRequest) ([]Role, int64, error) {
	var roles []Role
	var total int64

	query := db.WithContext(ctx).Model(&Role{})

	// search filter
	if req.Search != "" {
		like := "%" + req.Search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", like, like)
	}

	// system filter
	if req.IsSystem != nil {
		query = query.Where("is_system = ?", *req.IsSystem)
	}

	// include deleted if requested
	if req.IncludeDeleted == nil || !*req.IncludeDeleted {
		query = query.Where("status != ?", "deleted")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// pagination
	if req.Page > 0 && req.PageSize > 0 {
		offset := (req.Page - 1) * req.PageSize
		query = query.Offset(offset).Limit(req.PageSize)
	}

	if err := query.Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

func (r *roleRepo) UpdateRoleStatus(ctx context.Context, db *gorm.DB, id uint, status string) error {
	return db.WithContext(ctx).Model(&Role{}).Where("id = ?", id).Update("status", status).Error
}

func (r *roleRepo) HasActiveUsers(ctx context.Context, db *gorm.DB, roleID uint) (bool, error) {
	var count int64
	err := db.WithContext(ctx).Table("users").
		Where("role_id = ? AND status = ?", roleID, "active").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
