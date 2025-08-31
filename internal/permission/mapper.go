package permission

// ----------------------
// Permission Mappers
// ----------------------

func ToPermissionResponse(p Permission) PermissionResponse {
	return PermissionResponse{
		ID:          p.ID,
		Page:        p.Page,
		Action:      p.Action,
		Key:         p.Key,
		
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
	}
}

func ToPermissionResponses(list []Permission) []PermissionResponse {
	res := make([]PermissionResponse, len(list))
	for i, p := range list {
		res[i] = ToPermissionResponse(p)
	}
	return res
}

// ----------------------
// RolePermission Mappers
// ----------------------

func ToRolePermissionResponse(rp RolePermission) RolePermissionResponse {
	return RolePermissionResponse{
		ID:           rp.ID,
		RoleID:       rp.RoleID,
		PermissionID: rp.PermissionID,
		Granted:      rp.Granted,
		CreatedAt:    rp.CreatedAt,
		UpdatedAt:    rp.UpdatedAt,
	}
}

func ToRolePermissionResponses(list []RolePermission) []RolePermissionResponse {
	res := make([]RolePermissionResponse, len(list))
	for i, rp := range list {
		res[i] = ToRolePermissionResponse(rp)
	}
	return res
}

// ----------------------
// UserPermission Mappers
// ----------------------

func ToUserPermissionResponse(up UserPermission) UserPermissionResponse {
	return UserPermissionResponse{
		ID:           up.ID,
		UserID:       up.UserID,
		PermissionID: up.PermissionID,
		Granted:      up.Granted,
		CreatedAt:    up.CreatedAt,
		UpdatedAt:    up.UpdatedAt,
	}
}

func ToUserPermissionResponses(list []UserPermission) []UserPermissionResponse {
	res := make([]UserPermissionResponse, len(list))
	for i, up := range list {
		res[i] = ToUserPermissionResponse(up)
	}
	return res
}