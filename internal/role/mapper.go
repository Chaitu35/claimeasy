package role

// MapCreateRoleRequestToModel - maps Create DTO to model
func MapCreateRoleRequestToModel(req CreateRoleRequest) *Role {
	return &Role{
		Name:        req.Name,
		Description: req.Description,
		IsSystem:    req.IsSystem,
		Status:      "active", // default on creation
	}
}

// MapUpdateRoleRequestToModel - updates only provided fields
func MapUpdateRoleRequestToModel(role *Role, req UpdateRoleRequest) *Role {
	if req.Name != nil {
		role.Name = *req.Name
	}
	if req.Description != nil {
		role.Description = *req.Description
	}
	if req.IsSystem != nil {
		role.IsSystem = *req.IsSystem
	}
	// Note: no status here — handled by Activate/Deactivate APIs
	return role
}

// MapRoleToResponse - maps model to response DTO
func MapRoleToResponse(role Role) RoleResponse {
	return RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		Status:      role.Status,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

// MapRolesToListResponse - maps slice + metadata
func MapRolesToListResponse(roles []Role, total int64, page, pageSize int) ListRolesResponse {
	responses := make([]RoleResponse, len(roles))
	for i, r := range roles {
		responses[i] = MapRoleToResponse(r)
	}
	return ListRolesResponse{
		Roles:      responses,
		TotalCount: total,
		Page:       page,
		PageSize:   pageSize,
	}
}
// MapUpdateRoleStatusRequestToModel - updates only the status field
func MapUpdateRoleStatusRequestToModel(role *Role, req UpdateRoleStatusRequest) *Role {
	role.Status = req.Status
	return role
}
