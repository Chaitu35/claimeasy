package permission

import "time"
type CreatePermissionRequest struct {
	Page        string `json:"page" binding:"required"`
	Action      Action `json:"action" binding:"required"`
	Description string `json:"description"`
	 // active, inactive
}

type UpdatePermissionRequest struct {
	
	
	
	Description *string `json:"description,omitempty"`
}


type PermissionResponse struct {
	ID          uint      `json:"id"`
	Page        string    `json:"page"`
	Action      Action    `json:"action"`
	Key         string    `json:"key"`
	
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ----------------------
// Role Permission DTOs
// ----------------------

type AssignRolePermissionRequest struct {
	RoleID       uint `json:"role_id" binding:"required"`
	PermissionID uint `json:"permission_id" binding:"required"`
	Granted      bool `json:"granted"`
}

type UpdateRolePermissionRequest struct {
	Granted *bool `json:"granted,omitempty"`
}

type UpdateRolePermissionStatusRequest struct {
	Granted bool `json:"granted" binding:"required"`
}

type RolePermissionResponse struct {
	ID           uint      `json:"id"`
	RoleID       uint      `json:"role_id"`
	PermissionID uint      `json:"permission_id"`
	Granted      bool      `json:"granted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ----------------------
// User Permission DTOs
// ----------------------

type AssignUserPermissionRequest struct {
	UserID       uint `json:"user_id" binding:"required"`
	PermissionID uint `json:"permission_id" binding:"required"`
	Granted      bool `json:"granted"`
}

type UpdateUserPermissionRequest struct {
	Granted *bool `json:"granted,omitempty"`
}

type UpdateUserPermissionStatusRequest struct {
	Granted bool `json:"granted" binding:"required"`
}

type UserPermissionResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	PermissionID uint      `json:"permission_id"`
	Granted      bool      `json:"granted"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
