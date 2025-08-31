package user

// MapCreateUserRequestToUser converts CreateUserRequest into a User model.
// Note: The password must already be hashed in the service layer.
func MapCreateUserRequestToUser(req CreateUserRequest, hashedPassword string) *User {
	return &User{
		Username:       req.Username,
		Email:          req.Email,
		PasswordHash:   hashedPassword,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		Address:        req.Address,
		ProfilePicture: req.ProfilePicture,
		NationalID:     req.NationalID,
		ExpiryDate:     req.ExpiryDate,
		Need2FAuth:     req.Need2FAuth,
		TwoFAMethod:    req.TwoFAMethod,
		RoleID:         req.RoleID,
	}
}

// MapUpdateUserRequestToUser applies UpdateUserRequest to an existing User.
// Only updates fields that are provided (non-nil).
func MapUpdateUserRequestToUser(user *User, req UpdateUserRequest) {
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Address != nil {
		user.Address = *req.Address
	}
	if req.ProfilePicture != nil {
		user.ProfilePicture = *req.ProfilePicture
	}
	if req.ExpiryDate != nil {
		user.ExpiryDate = *req.ExpiryDate
	}
	if req.Need2FAuth != nil {
		user.Need2FAuth = *req.Need2FAuth
	}
	if req.TwoFAMethod != nil {
		user.TwoFAMethod = *req.TwoFAMethod
	}
	if req.RoleID != nil {
		user.RoleID = *req.RoleID
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.IsValidated != nil {
		user.IsValidated = *req.IsValidated
	}
}

// MapUserToResponse converts a User model to UserResponse.
func MapUserToResponse(u *User) *UserResponse {
	if u == nil {
		return nil
	}
	return &UserResponse{
		ID:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		FirstName:      u.FirstName,
		LastName:       u.LastName,
		Phone:          u.Phone,
		Address:        u.Address,
		ProfilePicture: u.ProfilePicture,
		ExpiryDate:     u.ExpiryDate,
		Need2FAuth:     u.Need2FAuth,
		TwoFAMethod:    u.TwoFAMethod,
		IsActive:       u.IsActive,
		IsValidated:    u.IsValidated,
		RoleID:         u.RoleID,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		
	}
}

// MapUsersToResponses converts a slice of User models into a slice of UserResponse DTOs.
func MapUsersToResponses(users []*User) []*UserResponse {
	responses := make([]*UserResponse, 0, len(users))
	for _, u := range users {
		responses = append(responses, MapUserToResponse(u))
	}
	return responses
}

// MapPagedUsersToResponse builds a PagedUserResponse from users and pagination data.
func MapPagedUsersToResponse(users []*User, totalCount int64, page, pageSize int) *PagedUserResponse {
	return &PagedUserResponse{
		Users:      MapUsersToResponses(users),
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
	}
}

// MapUserToMap converts User model into a flat map (for auditing diffs/snapshots).
func MapUserToMap(u *User) map[string]interface{} {
	if u == nil {
		return nil
	}
	return map[string]interface{}{
		"ID":             u.ID,
		"Username":       u.Username,
		"Email":          u.Email,
		"FirstName":      u.FirstName,
		"LastName":       u.LastName,
		"Phone":          u.Phone,
		"Address":        u.Address,
		"ProfilePicture": u.ProfilePicture,
		"NationalID":     u.NationalID,
		"ExpiryDate":     u.ExpiryDate,
		"TwoFAMethod":    u.TwoFAMethod,
		"Need2FAuth":     u.Need2FAuth,
		"RoleID":         u.RoleID,
		"IsActive":       u.IsActive,
		"IsValidated":    u.IsValidated,
		"CreatedAt":      u.CreatedAt,
		"UpdatedAt":      u.UpdatedAt,
		"DeletedAt":      u.DeletedAt,
	}
}
