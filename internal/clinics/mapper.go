package clinics

import (
	"time"
)

// Map Clinic model → ClinicResponseDTO
func MapClinicToDTO(clinic *Clinic) *ClinicResponseDTO {
	dto := &ClinicResponseDTO{
		ID:             clinic.ID,
		Name:           clinic.Name,
		ControlBody:    string(clinic.ControlBody),
		Phone:          clinic.Phone,
		Email:          clinic.Email,
		FacilityID:     clinic.FacilityID,
		Licence:        clinic.Licence,
		Representative: clinic.Representative,
		CreatedAt:      clinic.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      clinic.UpdatedAt.Format(time.RFC3339),
	}

	// Include assigned user if exists
	if clinic.UserID != nil && clinic.User != nil {
		dto.AssignedUserID = clinic.UserID
		dto.AssignedUserName = clinic.User.FirstName + " " + clinic.User.LastName
	}

	return dto
}

func MapClinicResponseDTOToMap(dto *ClinicResponseDTO) map[string]interface{} {
	return map[string]interface{}{
		"id":                 dto.ID,
		"name":               dto.Name,
		"control_body":       dto.ControlBody,
		"phone":              dto.Phone,
		"email":              dto.Email,
		"facility_id":        dto.FacilityID,
		"licence":            dto.Licence,
		"representative":     dto.Representative,
		"assigned_user_id":   dto.AssignedUserID,
		"assigned_user_name": dto.AssignedUserName,
	}
}

// MapClinicToMap → convert Clinic model to map for audit diff
func MapClinicToMap(clinic *Clinic) map[string]interface{} {
	m := map[string]interface{}{
		"id":             clinic.ID,
		"name":           clinic.Name,
		"control_body":   string(clinic.ControlBody),
		"phone":          clinic.Phone,
		"email":          clinic.Email,
		"facility_id":    clinic.FacilityID,
		"licence":        clinic.Licence,
		"representative": clinic.Representative,
		"user_id":        clinic.UserID,
	}

	if clinic.User != nil {
		m["assigned_user_name"] = clinic.User.FirstName + " " + clinic.User.LastName
	}

	return m
}

// Map slice of Clinics → slice of ClinicResponseDTOs
func MapClinicsToDTO(clinics []*Clinic) []*ClinicResponseDTO {
	dtos := make([]*ClinicResponseDTO, len(clinics))
	for i, c := range clinics {
		dtos[i] = MapClinicToDTO(c)
	}
	return dtos
}

// Map CreateClinicDTO → Clinic model
func MapCreateDTOToClinic(dto *CreateClinicDTO) *Clinic {
	return &Clinic{
		Name:           dto.Name,
		ControlBody:    dto.ControlBody,
		Phone:          dto.Phone,
		Email:          dto.Email,
		FacilityID:     dto.FacilityID,
		Licence:        dto.Licence,
		Representative: dto.Representative,
	}
}

// Apply updates from UpdateClinicDTO → Clinic model
func ApplyUpdateDTOToClinic(clinic *Clinic, dto *UpdateClinicDTO) {
	if dto.Name != nil {
		clinic.Name = *dto.Name
	}
	if dto.ControlBody != nil {
    clinic.ControlBody = *dto.ControlBody // no need to cast
}

	if dto.Phone != nil {
		clinic.Phone = *dto.Phone
	}
	if dto.Email != nil {
		clinic.Email = *dto.Email
	}
	if dto.FacilityID != nil {
		clinic.FacilityID = *dto.FacilityID
	}
	if dto.Licence != nil {
		clinic.Licence = *dto.Licence
	}
	if dto.Representative != nil {
		clinic.Representative = *dto.Representative
	}
	// AssignUserID is handled separately in ClinicService
}
