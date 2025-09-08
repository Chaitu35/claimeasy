package user

import "time"

// ------------------- Requests -------------------

// CreateUserRequest is used when registering a new user.
type CreateUserRequest struct {
	Username       string        `json:"username" binding:"required,min=3,max=50"`
	Email          string        `json:"email" binding:"required,email"`
	Password       string        `json:"password" binding:"required,min=8"`
	FirstName      string        `json:"first_name" binding:"required"`
	LastName       string        `json:"last_name" binding:"required"`
	Phone          string        `json:"phone" binding:"required"`
	Address        string        `json:"address,omitempty"`
	ProfilePicture string        `json:"profile_picture,omitempty"`
	NationalID     string        `json:"national_id" binding:"required"`
	ExpiryDate     time.Time     `json:"expiry_date" binding:"required"`
	TwoFAMethod    TwoFactorType `json:"two_fa_method,omitempty"`
	Need2FAuth     bool          `json:"need_2f_auth,omitempty"`
	RoleID         uint          `json:"role_id" binding:"required"`
}

// UpdateUserRequest allows updating user details (partial update).
type UpdateUserRequest struct {
	FirstName      *string        `json:"first_name,omitempty"`
	LastName       *string        `json:"last_name,omitempty"`
	Phone          *string        `json:"phone,omitempty"`
	Address        *string        `json:"address,omitempty"`
	ProfilePicture *string        `json:"profile_picture,omitempty"`
	ExpiryDate     *time.Time     `json:"expiry_date,omitempty"`
	Need2FAuth     *bool          `json:"need_2f_auth,omitempty"`
	TwoFAMethod    *TwoFactorType `json:"two_fa_method,omitempty"`
	RoleID         *uint          `json:"role_id,omitempty"`
	IsActive       *bool          `json:"is_active,omitempty"`
	IsValidated    *bool          `json:"is_validated,omitempty"`
}

// SearchUserRequest allows searching/filtering users.
type SearchUserRequest struct {
	Name     string `json:"name,omitempty"`   // fuzzy search on first_name/last_name
	Email    string `json:"email,omitempty"`  // fuzzy search on email
	Phone    string `json:"phone,omitempty"`
	RoleID   *uint  `json:"role_id,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	IncludeDeleted bool   `json:"include_deleted,omitempty"` // if true, include soft-deleted users
	Page     int `form:"page,default=1" json:"page,omitempty" binding:"gte=1"`
    PageSize int `form:"page_size,default=20" json:"page_size,omitempty" binding:"gte=1,lte=100"`

}

// ------------------- Responses -------------------

// UserResponse represents a safe version of User returned in APIs.
type UserResponse struct {
	ID             uint          `json:"id"`
	Username       string        `json:"username"`
	Email          string        `json:"email"`
	FirstName      string        `json:"first_name"`
	LastName       string        `json:"last_name"`
	Phone          string        `json:"phone"`
	Address        string        `json:"address,omitempty"`
	ProfilePicture string        `json:"profile_picture,omitempty"`
	ExpiryDate     time.Time     `json:"expiry_date"`
	Need2FAuth     bool          `json:"need_2f_auth"`
	TwoFAMethod    TwoFactorType `json:"two_fa_method,omitempty"`
	IsActive       bool          `json:"is_active"`
	IsValidated    bool          `json:"is_validated"`
	RoleID         uint          `json:"role_id"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// PagedUserResponse wraps a list of users with pagination info.
type PagedUserResponse struct {
    Users      []*UserResponse `json:"users"`
    TotalCount int64           `json:"total_count"`
    Page       int             `json:"page"`
    PageSize   int             `json:"page_size"`
}
