package user

import (
	"context"
	"fmt"

	"github.com/chaitu35/claimeasy/internal/audit"
	"github.com/chaitu35/claimeasy/internal/permission"
	"github.com/chaitu35/claimeasy/internal/utils"
	"gorm.io/gorm"
)

type Service struct {
	repo   Repository
	db     *gorm.DB
	logger *audit.Logger
    permRepo permission.Repository
    clinicSvc ClinicService 
}

type ClinicService interface {
    RevokeClinicsForUser(ctx context.Context, tx *gorm.DB, userID uint) error
}

func NewService(repo Repository, db *gorm.DB, logger *audit.Logger, permRepo permission.Repository, clinicSvr ClinicService) *Service {
	return &Service{
		repo:   repo,
		db:     db,
		logger: logger,
        permRepo: permRepo,
        clinicSvc : clinicSvr,
	}
}
// withTransaction is a helper to wrap operations in a DB transaction.
func (s *Service) withTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
    return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        return fn(tx)
    })
}

// ---------------- Core CRUD ----------------

func (s *Service) CreateUser(ctx context.Context, req CreateUserRequest) (*UserResponse, error) {
    hashedPwd, err := utils.HashPassword(req.Password)
    if err != nil {
        return nil, err
    }

    user := &User{
        Username:       req.Username,
        Email:          req.Email,
        PasswordHash:   hashedPwd,
        FirstName:      req.FirstName,
        LastName:       req.LastName,
        Phone:          req.Phone,
        Address:        req.Address,
        ProfilePicture: req.ProfilePicture,
        NationalID:     req.NationalID,
        ExpiryDate:     req.ExpiryDate,
        TwoFAMethod:    req.TwoFAMethod,
        Need2FAuth:     req.Need2FAuth,
        RoleID:         req.RoleID,
        IsActive:       true,
        IsValidated:    false,
    }

    var result *UserResponse
    err = s.withTransaction(ctx, func(tx *gorm.DB) error {
        if err := s.repo.CreateUser(ctx, tx, user); err != nil {
            return err
        }
       rolePerms, err := s.permRepo.ListRolePermissions(ctx, tx, req.RoleID)
		if err != nil {
			return err
		}

		// assign each permission to user
		for _, rp := range rolePerms {
			up := &permission.UserPermission{
				UserID:       user.ID,
				PermissionID: rp.PermissionID,
				Granted:      true,
			}
			if err := s.permRepo.AssignPermissionToUser(ctx, tx, up); err != nil {
				return err
			}
		}


        if err := s.logger.LogSnapshot(ctx, "CREATE", "user", user.ID, user, "User created"); err != nil {
            return err
        }
        result = MapUserToResponse(user)
        return nil
    })

    return result, err
}

func (s *Service) UpdateUser(ctx context.Context, id uint, req UpdateUserRequest) (*UserResponse, error) {
    var result *UserResponse
    err := s.withTransaction(ctx, func(tx *gorm.DB) error {
        user, err := s.repo.GetUserByID(ctx, tx, id, false)
        if err != nil {
            return err
        }

        before := MapUserToMap(user)

        // Apply updates
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

        // Save
        if err := s.repo.UpdateUser(ctx, tx, user); err != nil {
            return err
        }

        after := MapUserToMap(user)
        if err := s.logger.LogDiff(ctx, "UPDATE", "user", user.ID, before, after, "User updated"); err != nil {
            return err
        }

        result = MapUserToResponse(user)
        return nil
    })

    return result, err
}

// ---------------- Lifecycle / Status ----------------

func (s *Service) UpdateUserStatus(ctx context.Context, userID uint, action UserStatusAction) error {
    return s.withTransaction(ctx, func(tx *gorm.DB) error {
        user, err := s.repo.GetUserByID(ctx, tx, userID, true)
        if err != nil {
            return err
        }

        before := MapUserToMap(user)

        if err := s.repo.UpdateUserStatus(ctx, tx, userID, action); err != nil {
            return err
        }

        // Revoke assigned clinics if user is deactivated or deleted
        if action == UserActionDeactivate || action == UserActionDelete {
            if err := s.clinicSvc.RevokeClinicsForUser(ctx, tx,userID); err != nil {
                return err
            }
        }

        // Refresh after update
        user,err= s.repo.GetUserByID(ctx, tx, userID, true)
        if err != nil {
            return err
        }
        return s.logger.LogDiff(ctx, "STATUS_CHANGE", "user", user.ID, before, MapUserToMap(user), fmt.Sprintf("User %s", action))
    })
}


// ---------------- Fetch / Search ----------------

func (s *Service) GetUserByID(ctx context.Context, id uint, includeDeleted bool) (*UserResponse, error) {
	user, err := s.repo.GetUserByID(ctx, s.db, id, includeDeleted)
	if err != nil {
		return nil, err
	}

	// Optional: audit view
	_ = s.logger.LogSnapshot(ctx, "VIEW", "user", user.ID, user, "User viewed")

	return MapUserToResponse(user), nil
}

func (s *Service) SearchUsers(ctx context.Context, req SearchUserRequest) (*PagedUserResponse, error) {
	users, total, err := s.repo.SearchUsers(ctx, s.db, req)
	if err != nil {
		return nil, err
	}

	responses := MapUsersToResponses(users)

	return &PagedUserResponse{
		Users:      responses,
		TotalCount: total,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}, nil
}
