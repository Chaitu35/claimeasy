package dto


type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
type UserDTO struct {
	ID string `json:"id"`
	Username string `json:"username"`	
	Role string `json:"role"`
	Permissions  map[string]PermissionsDTO `json:"permissions"`
}
type PermissionsDTO struct {
	View bool `json:"view"`
	Edit bool `json:"edit"`
	Delete bool `json:"delete"`
	Create bool `json:"create"`
	Print bool `json:"print"`
	Export bool `json:"export"`
}

type ForgetPasswordRequest struct{
	Email string `json:"email" validate:"required,email"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User UserDTO `json:"user"`
}
type ForgetPasswordResponse struct{
	Message string `json:"message"`

}
type ChangePasswordResponse struct {
	Message string `json:"message"`	
}