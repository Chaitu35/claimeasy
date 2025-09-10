package clinics


type CreateClinicDTO struct {
	Name           string `json:"name" binding:"required"`
	ControlBody    string `json:"control_body" binding:"required,oneof=DHA HAAD MOH"`
	Phone          string `json:"phone" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	FacilityID     string `json:"facility_id"`
	Licence        string `json:"licence"`
	Representative string `json:"representative" binding:"required"`
}

type UpdateClinicDTO struct {
	Name           *string `json:"name"`
	ControlBody    *string `json:"control_body" binding:"omitempty,oneof=DHA HAAD MOH"`
	Phone          *string `json:"phone"`
	Email          *string `json:"email" binding:"omitempty,email"`
	FacilityID     *string `json:"facility_id"`
	Licence        *string `json:"licence"`
	Representative *string `json:"representative"`
	AssignUserID   *uint   `json:"assign_user_id,omitempty"` // optional
}

type ClinicResponseDTO struct {
	ID              uint   `json:"id"`
	Name            string `json:"name"`
	ControlBody     string `json:"control_body"`
	Phone           string `json:"phone"`
	Email           string `json:"email"`
	FacilityID      string `json:"facility_id"`
	Licence         string `json:"licence"`
	Representative  string `json:"representative"`
	AssignedUserID  *uint  `json:"assigned_user_id,omitempty"`  // active assigned user
	AssignedUserName string `json:"assigned_user_name,omitempty"` // optional
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}


type AssignUserDTO struct {
	UserID uint `json:"user_id" binding:"required"`
}

type ListClinicsQuery struct {
	Page          int    `form:"page"`            // default handled in handler
	PageSize      int    `form:"page_size"`       // default handled in handler
	Search        string `form:"search"`          // optional fuzzy search on name/email/phone
	IncludeDeleted bool  `form:"include_deleted"` // pass "true" to include soft-deleted rows
}


type ListClinicsResponse struct {
	Clinics     []*ClinicResponseDTO `json:"clinics"`
	TotalCount  int64                `json:"total_count"`
	CurrentPage int                  `json:"current_page"`
	PageSize    int                  `json:"page_size"`
}

type UpdateClinicStatusDTO struct {
	
	Active bool `json:"active" binding:"required"`
}
