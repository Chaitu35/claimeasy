package role

import (
	"context"
	"fmt"

	"github.com/chaitu35/claimeasy/internal/audit"
	"gorm.io/gorm"
)

type Service interface {
	CreateRole(ctx context.Context, db *gorm.DB, req CreateRoleRequest) (RoleResponse, error)
	UpdateRole(ctx context.Context, db *gorm.DB, id uint, req UpdateRoleRequest) (RoleResponse, error)
	ListRoles(ctx context.Context, db *gorm.DB, req ListRolesRequest) (ListRolesResponse, error)
	UpdateRoleStatus(ctx context.Context, db *gorm.DB, id uint, req UpdateRoleStatusRequest) (RoleResponse, error)
	GetRoleByID(ctx context.Context, db *gorm.DB, id uint) (RoleResponse, error)
}

type service struct {
	repo   Repository
	logger *audit.Logger
}

func NewService(repo Repository, logger *audit.Logger) Service {
	return &service{repo: repo, logger: logger}
}

// ----------------------------
// CRUD & Status Methods
// ----------------------------

func (s *service) CreateRole(ctx context.Context, db *gorm.DB, req CreateRoleRequest) (RoleResponse, error) {
	role := MapCreateRoleRequestToModel(req)
	if err := s.repo.CreateRole(ctx, db, role); err != nil {
		return RoleResponse{}, err
	}

	// audit
	_ = s.logger.LogSnapshot(ctx, audit.ActionCreate, "role", role.ID, role, "Role created")

	return MapRoleToResponse(*role), nil
}

func (s *service) UpdateRole(ctx context.Context, db *gorm.DB, id uint, req UpdateRoleRequest) (RoleResponse, error) {
	role, err := s.repo.GetRoleByID(ctx, db, id)
	if err != nil {
		return RoleResponse{}, err
	}

	before := *role
	role = MapUpdateRoleRequestToModel(role, req)

	if err := s.repo.UpdateRole(ctx, db, role); err != nil {
		return RoleResponse{}, err
	}

	// audit diff
	beforeMap := map[string]interface{}{
		"name":        before.Name,
		"description": before.Description,
		"is_system":   before.IsSystem,
		"status":      before.Status,
	}
	afterMap := map[string]interface{}{
		"name":        role.Name,
		"description": role.Description,
		"is_system":   role.IsSystem,
		"status":      role.Status,
	}
	_ = s.logger.LogDiff(ctx, audit.ActionUpdate, "role", role.ID, beforeMap, afterMap, "Role updated")

	return MapRoleToResponse(*role), nil
}

func (s *service) ListRoles(ctx context.Context, db *gorm.DB, req ListRolesRequest) (ListRolesResponse, error) {
	roles, total, err := s.repo.ListRoles(ctx, db, req)
	if err != nil {
		return ListRolesResponse{}, err
	}
	page := req.Page
	if page == 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize == 0 {
		pageSize = int(total)
	}
	return MapRolesToListResponse(roles, total, page, pageSize), nil
}

func (s *service) UpdateRoleStatus(ctx context.Context, db *gorm.DB, id uint, req UpdateRoleStatusRequest) (RoleResponse, error) {
	role, err := s.repo.GetRoleByID(ctx, db, id)
	if err != nil {
		return RoleResponse{}, err
	}

	// prevent deleting roles with active users
	if req.Status == "deleted" {
		hasUsers, err := s.repo.HasActiveUsers(ctx, db, id)
		if err != nil {
			return RoleResponse{}, err
		}
		if hasUsers {
			return RoleResponse{}, fmt.Errorf("cannot delete role with active users")
		}
	}

	beforeStatus := role.Status
	role = MapUpdateRoleStatusRequestToModel(role, req)

	if err := s.repo.UpdateRoleStatus(ctx, db, id, role.Status); err != nil {
		return RoleResponse{}, err
	}

	// audit
	_ = s.logger.LogDiff(ctx, audit.ActionUpdate, "role", role.ID,
		map[string]interface{}{"status": beforeStatus},
		map[string]interface{}{"status": role.Status},
		req.Reason)

	return MapRoleToResponse(*role), nil
}

func (s *service) GetRoleByID(ctx context.Context, db *gorm.DB, id uint) (RoleResponse, error) {
	role, err := s.repo.GetRoleByID(ctx, db, id)
	if err != nil {
		return RoleResponse{}, err
	}
	return MapRoleToResponse(*role), nil
}
