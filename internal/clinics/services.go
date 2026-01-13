package clinics

import (
	"context"
	"fmt"
	"time"

	"github.com/chaitu35/claimeasy/internal/audit"
	"github.com/chaitu35/claimeasy/internal/user"
	"gorm.io/gorm"
)

// Repository defines the DB operations needed by ClinicService
type Repository interface {
	WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	GetClinicByID(ctx context.Context, tx *gorm.DB, id uint, includeDeleted bool) (*Clinic, error)
	ListClinics(ctx context.Context, query *ListClinicsQuery) ([]*Clinic, int64, error)
	ListActiveUserClinics(ctx context.Context, userID uint) ([]UserClinic, error)
RevokeUserClinic(ctx context.Context, uc *UserClinic) error
RevokeClinicsForUser(ctx context.Context, tx *gorm.DB, userID uint) error
}

// Service provides clinic operations
type Service struct {
	repo   Repository
	logger *audit.Logger
	db    *gorm.DB
}
func NewService(repo Repository, logger *audit.Logger,db *gorm.DB) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
		db:	 db,
    }

}
// CreateClinic creates a new clinic
func (s *Service) CreateClinic(ctx context.Context, dto *CreateClinicDTO) (*ClinicResponseDTO, error) {
	clinic := MapCreateDTOToClinic(dto)

	err := s.repo.WithTransaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(clinic).Error; err != nil {
			return err
		}
		// Audit snapshot
		return s.logger.LogSnapshot(ctx, audit.ActionCreate, "clinic", clinic.ID, MapClinicToMap(clinic), "Created new clinic")
	})
	if err != nil {
		return nil, err
	}
	return MapClinicToDTO(clinic), nil
}

// UpdateClinic updates clinic fields and optionally assigns a user
func (s *Service) UpdateClinic(ctx context.Context, clinicID uint, dto *UpdateClinicDTO) (*ClinicResponseDTO, error) {
	var updatedClinic *Clinic

	err := s.repo.WithTransaction(ctx, func(tx *gorm.DB) error {
		clinic, err := s.repo.GetClinicByID(ctx, tx, clinicID, true)
		if err != nil {
			return err
		}

		before := MapClinicToMap(clinic)

		ApplyUpdateDTOToClinic(clinic, dto)

		// Handle user assignment if provided
		if dto.AssignUserID != nil {
			if err := s.AssignUserToClinic(ctx,  clinicID, *dto.AssignUserID); err != nil {
				return err
			}
		}

		if err := tx.Save(clinic).Error; err != nil {
			return err
		}

		after := MapClinicToMap(clinic)
		if err := s.logger.LogDiff(ctx, audit.ActionUpdate, "clinic", clinic.ID, before, after, "Updated clinic"); err != nil {
			return fmt.Errorf("audit log failed: %w", err)
		}

		updatedClinic = clinic
		return nil
	})

	if err != nil {
		return nil, err
	}

	return MapClinicToDTO(updatedClinic), nil
}

// assignUserToClinic revokes previous user and assigns a new user
func (s *Service) AssignUserToClinic(ctx context.Context, clinicID, userID uint) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // check if clinic exists
        var clinic Clinic
        if err := tx.First(&clinic, clinicID).Error; err != nil {
            return fmt.Errorf("clinic not found: %w", err)
        }

        // check if user exists
        var u user.User
        if err := tx.First(&u, userID).Error; err != nil {
            return fmt.Errorf("user not found: %w", err)
        }

        // revoke any previous assignment for this clinic
        if err := tx.Model(&UserClinic{}).
            Where("clinic_id = ? AND revoked_at IS NULL", clinicID).
            Update("revoked_at", time.Now()).
            Error; err != nil {
            return err
        }

        // insert new assignment into user_clinic
        if err := tx.Create(&UserClinic{
            UserID:    userID,
            ClinicID:  clinicID,
            AssignedAt: time.Now(),
        }).Error; err != nil {
            return err
        }

        // update clinics.user_id to reflect latest active user
        if err := tx.Model(&Clinic{}).
            Where("id = ?", clinicID).
            Update("user_id", userID).
            Error; err != nil {
            return err
        }

        return nil
    })
}


// RevokeUserFromClinic soft-deletes UserClinic record
func (s *Service) RevokeUserFromClinic(ctx context.Context, tx *gorm.DB, userID, clinicID uint) error {
	return tx.Model(&UserClinic{}).
		Where("user_id = ? AND clinic_id = ? AND revoked_at IS NULL", userID, clinicID).
		Update("revoked_at", time.Now()).Error
}

// UpdateClinicStatus updates the active status of a clinic
func (s *Service) UpdateClinicStatus(ctx context.Context, clinicID uint, active bool) error {
	return s.repo.WithTransaction(ctx, func(tx *gorm.DB) error {
		clinic, err := s.repo.GetClinicByID(ctx, tx, clinicID, true)
		if err != nil {
			return err
		}

		before := MapClinicToMap(clinic)

		if err := tx.Model(clinic).Update("active", active).Error; err != nil {
			return err
		}

		after := MapClinicToMap(clinic)
		return s.logger.LogDiff(ctx, audit.ActionUpdate, "clinic", clinic.ID, before, after, "Updated clinic status")
	})
}

// ListClinics returns paginated clinics with optional search and soft-deleted inclusion
func (s *Service) ListClinics(ctx context.Context, query *ListClinicsQuery) (*ListClinicsResponse, error) {
	clinics, totalCount, err := s.repo.ListClinics(ctx, query)
	if err != nil {
		return nil, err
	}

	return &ListClinicsResponse{
		Clinics:     MapClinicsToDTO(clinics),
		TotalCount:  totalCount,
		CurrentPage: query.Page,
		PageSize:    query.PageSize,
	}, nil
}
func (s *Service) ListActiveUserClinics(ctx context.Context, userID uint) ([]UserClinic, error) {
    return s.repo.ListActiveUserClinics(ctx, userID)
}

func (s *Service) RevokeUserClinic(ctx context.Context, uc *UserClinic) error {
    return s.repo.RevokeUserClinic(ctx, uc)
}

func (s *Service) RevokeClinicsForUser(ctx context.Context, tx *gorm.DB, userID uint) error {
	return tx.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&UserClinic{}).Error
}