package permission

import (
	"context"

	"github.com/chaitu35/claimeasy/internal/audit"
	"gorm.io/gorm"
)

// ----------------------
// Service Interface
// ----------------------

type Service interface {
	// Permission CRUD
	CreatePermission(ctx context.Context, db *gorm.DB, req CreatePermissionRequest) (PermissionResponse, error)
	UpdatePermission(ctx context.Context, db *gorm.DB,id uint, req UpdatePermissionRequest) (PermissionResponse, error)
	GetPermissionByID(ctx context.Context, db *gorm.DB, id uint) (PermissionResponse, error)
	ListPermissions(ctx context.Context, db *gorm.DB) ([]PermissionResponse, error)
   // GetPermissionsByID(ctx context.Context, db *gorm.DB, ids []uint) ([]PermissionResponse, error)
	// Role ↔ Permission
	AssignPermissionToRole(ctx context.Context, db *gorm.DB, req AssignRolePermissionRequest) (RolePermissionResponse, error)
	RevokePermissionFromRole(ctx context.Context, db *gorm.DB, roleID, permissionID uint, reason string) error
	ListRolePermissions(ctx context.Context, db *gorm.DB, roleID uint) ([]RolePermissionResponse, error)

	// User ↔ Permission
	AssignPermissionToUser(ctx context.Context, db *gorm.DB, req AssignUserPermissionRequest) (UserPermissionResponse, error)
	RevokePermissionFromUser(ctx context.Context, db *gorm.DB, userID, permissionID uint, reason string) error
	ListUserPermissions(ctx context.Context, db *gorm.DB, userID uint) ([]UserPermissionResponse, error)
}

type service struct {
	repo   Repository
	auditor *audit.Logger
}

func NewService(repo Repository, auditor *audit.Logger) Service {
	return &service{repo: repo, auditor: auditor}
}

// ----------------------
// Permission CRUD
// ----------------------

func (s *service) CreatePermission(ctx context.Context, db *gorm.DB, req CreatePermissionRequest) (PermissionResponse, error) {
	perm := Permission{
		Page:        req.Page,
		Action:      req.Action,
		Key:         string(req.Action) + ":" + req.Page,
		Description: req.Description,
	}

	if err := s.repo.CreatePermission(ctx, db, &perm); err != nil {
		return PermissionResponse{}, err
	}

	// Audit
	_ = s.auditor.LogSnapshot(ctx, audit.ActionCreate, "Permission", perm.ID, perm, "Permission created")

	return ToPermissionResponse(perm), nil
}
func (s *service) UpdatePermission(ctx context.Context, db *gorm.DB, id uint, req UpdatePermissionRequest) (PermissionResponse, error) {
    // Fetch existing permission by ID
    perm, err := s.repo.GetPermissionByID(ctx, db, id)
    if err != nil {
        return PermissionResponse{}, err
    }

    // Apply updates only for non-nil fields
  
    if req.Description != nil {
        perm.Description = *req.Description
    }
   

    // Recompute Key if Page or Action updated
    perm.Key = string(perm.Action) + ":" + perm.Page

    // Save updated permission to DB
    if err := s.repo.UpdatePermission(ctx, db, perm); err != nil {
        return PermissionResponse{}, err
    }

    // Audit update snapshot
    _ = s.auditor.LogSnapshot(ctx, audit.ActionUpdate, "Permission", perm.ID, perm, "Permission updated")

    // Return updated permission response
    return ToPermissionResponse(*perm), nil
}

func (s *service) GetPermissionByID(ctx context.Context, db *gorm.DB, id uint) (PermissionResponse, error) {



	perm, err := s.repo.GetPermissionByID(ctx, db, id)
	if err != nil {
		return PermissionResponse{}, err
	}

	// Audit (view event)
	_ = s.auditor.LogEvent(ctx, audit.ActionView, "Permission", id, "Permission viewed")

	return ToPermissionResponse(*perm), nil
}

func (s *service) ListPermissions(ctx context.Context, db *gorm.DB) ([]PermissionResponse, error) {
	list, err := s.repo.ListPermissions(ctx, db)
	if err != nil {
		return nil, err
	}

	// Audit (search/list event)
	_ = s.auditor.LogEvent(ctx, audit.ActionSearch, "Permission", 0, "Permissions listed")

	return ToPermissionResponses(list), nil
}

// ----------------------
// Role ↔ Permission
// ----------------------

func (s *service) AssignPermissionToRole(ctx context.Context, db *gorm.DB, req AssignRolePermissionRequest) (RolePermissionResponse, error) {
	rp := RolePermission{
		RoleID:       req.RoleID,
		PermissionID: req.PermissionID,
		Granted:      req.Granted,
	}

	if err := s.repo.AssignPermissionToRole(ctx, db, &rp); err != nil {
		return RolePermissionResponse{}, err
	}

	// Audit
	_ = s.auditor.LogSnapshot(ctx, audit.ActionCreate, "RolePermission", rp.ID, rp, "Assigned permission to role")

	return ToRolePermissionResponse(rp), nil
}

func (s *service) RevokePermissionFromRole(ctx context.Context, db *gorm.DB, roleID, permissionID uint, reason string) error {
	if err := s.repo.RevokePermissionFromRole(ctx, db, roleID, permissionID); err != nil {
		return err
	}

	// Audit
	_ = s.auditor.LogSnapshot(ctx, audit.ActionDelete, "RolePermission", permissionID, map[string]interface{}{
		"roleID": roleID, "permissionID": permissionID,
	}, reason)

	return nil
}

func (s *service) ListRolePermissions(ctx context.Context, db *gorm.DB, roleID uint) ([]RolePermissionResponse, error) {
	list, err := s.repo.ListRolePermissions(ctx, db, roleID)
	if err != nil {
		return nil, err
	}

	// Audit
	_ = s.auditor.LogEvent(ctx, audit.ActionView, "RolePermission", roleID, "Listed role permissions")

	return ToRolePermissionResponses(list), nil
}

// ----------------------
// User ↔ Permission
// ----------------------

func (s *service) AssignPermissionToUser(ctx context.Context, db *gorm.DB, req AssignUserPermissionRequest) (UserPermissionResponse, error) {
	up := UserPermission{
		UserID:       req.UserID,
		PermissionID: req.PermissionID,
		Granted:      req.Granted,
	}

	if err := s.repo.AssignPermissionToUser(ctx, db, &up); err != nil {
		return UserPermissionResponse{}, err
	}

	// Audit
	_ = s.auditor.LogSnapshot(ctx, audit.ActionCreate, "UserPermission", up.ID, up, "Assigned permission to user")

	return ToUserPermissionResponse(up), nil
}

func (s *service) RevokePermissionFromUser(ctx context.Context, db *gorm.DB, userID, permissionID uint, reason string) error {
	if err := s.repo.RevokePermissionFromUser(ctx, db, userID, permissionID); err != nil {
		return err
	}

	// Audit
	_ = s.auditor.LogSnapshot(ctx, audit.ActionDelete, "UserPermission", permissionID, map[string]interface{}{
		"userID": userID, "permissionID": permissionID,
	}, reason)

	return nil
}

func (s *service) ListUserPermissions(ctx context.Context, db *gorm.DB, userID uint) ([]UserPermissionResponse, error) {
	list, err := s.repo.ListUserPermissions(ctx, db, userID)
	if err != nil {
		return nil, err
	}

	// Audit
	_ = s.auditor.LogEvent(ctx, audit.ActionView, "UserPermission", userID, "Listed user permissions")

	return ToUserPermissionResponses(list), nil
}
