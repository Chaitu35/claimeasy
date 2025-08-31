package role

import "time"

// CreateRoleRequest - when creating a new role
type CreateRoleRequest struct {
    Name        string `json:"name" binding:"required"`
    Description string `json:"description,omitempty"`
    IsSystem    bool   `json:"is_system,omitempty"` // optional; usually false
}

// UpdateRoleRequest - when updating an existing role
type UpdateRoleRequest struct {
    Name        *string `json:"name,omitempty"`
    Description *string `json:"description,omitempty"`
    IsSystem    *bool   `json:"is_system,omitempty"`
    // No Status here → status updates are handled by Deactivate/Activate APIs
}

// RoleResponse - what we send back to the client
type RoleResponse struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description,omitempty"`
    IsSystem    bool      `json:"is_system"`
    Status      string    `json:"status"`       // <-- NEW
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// ListRolesRequest - filters + pagination for listing roles
type ListRolesRequest struct {
    Search    string `form:"search"`     // fuzzy search on Name/Description
    IsSystem  *bool  `form:"is_system"`  // optional filter
    Status    *string `form:"status"` 
    IncludeDeleted *bool   `form:"include_deleted"` // new field   // <-- NEW optional filter (active/inactive)
    Page      int    `form:"page" binding:"omitempty,min=1"`
    PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// ListRolesResponse - response with roles + pagination metadata
type ListRolesResponse struct {
    Roles      []RoleResponse `json:"roles"`
    TotalCount int64          `json:"total_count"`
    Page       int            `json:"page"`
    PageSize   int            `json:"page_size"`
}
// UpdateRoleStatusRequest - used when activating/deactivating a role
type UpdateRoleStatusRequest struct {
    Status string `json:"status" binding:"required,oneof=active inactive blocked"`
    Reason string `json:"reason" binding:"required"`
     // optional reason for status change
}
