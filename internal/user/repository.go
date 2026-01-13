package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/chaitu35/claimeasy/internal/utils"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, db *gorm.DB, user *User) error
	UpdateUser(ctx context.Context, db *gorm.DB, user *User) error
	GetUserByID(ctx context.Context, db *gorm.DB, id uint, includeDeleted bool) (*User, error)
	SearchUsers(ctx context.Context, db *gorm.DB, req SearchUserRequest) ([]*User, int64, error)

	// Unified lifecycle operations
	UpdateUserStatus(ctx context.Context, db *gorm.DB, userID uint, action UserStatusAction) error

	// Password reset
	SavePasswordResetToken(ctx context.Context, tx *gorm.DB, userID uint, tokenHash string, expiresAt time.Time) error
    VerifyPasswordResetToken(ctx context.Context, tx *gorm.DB, token string) (uint, error)
    DeletePasswordResetToken(ctx context.Context, tx *gorm.DB, token string) error
}


type userRepo struct{}

func NewUserRepository() Repository {
	return &userRepo{}
}

// ---------------- Core CRUD ----------------

func (r *userRepo) CreateUser(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

func (r *userRepo) UpdateUser(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Save(user).Error
}

func (r *userRepo) GetUserByID(ctx context.Context, db *gorm.DB, id uint, includeDeleted bool) (*User, error) {
	var user User
	tx := db.WithContext(ctx)

	if includeDeleted {
		tx = tx.Unscoped()
	}

	if err := tx.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) SearchUsers(ctx context.Context, db *gorm.DB, req SearchUserRequest) ([]*User, int64, error) {
	var (
		users        []*User
		totalCount int64
	)

	tx := db.WithContext(ctx).Model(&User{})

	// include soft-deleted or not
	if req.IncludeDeleted {
		tx = tx.Unscoped()
	}

	// filters
	if req.Name != "" {
		tx = tx.Where("first_name ILIKE ? OR last_name ILIKE ?", "%"+req.Name+"%", "%"+req.Name+"%")
	}
	if req.Email != "" {
		tx = tx.Where("email ILIKE ?", "%"+req.Email+"%")
	}
	if req.Phone != "" {
		tx = tx.Where("phone = ?", req.Phone)
	}
	if req.RoleID != nil {
		tx = tx.Where("role_id = ?", *req.RoleID)
	}
	if req.IsActive != nil {
		tx = tx.Where("is_active = ?", *req.IsActive)
	}

	// count total
	if err := tx.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// pagination defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	offset := (req.Page - 1) * req.PageSize

	// fetch data
	if err := tx.
		Limit(req.PageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, totalCount, nil
}

// ---------------- Lifecycle Ops ----------------

func (r *userRepo) UpdateUserStatus(ctx context.Context, db *gorm.DB, userID uint, action UserStatusAction) error {
	switch action {
	case UserActionActivate:
		return db.Model(&User{}).
			Where("id = ?", userID).
			Update("is_active", true).Error

	case UserActionDeactivate:
		return db.Model(&User{}).
			Where("id = ?", userID).
			Update("is_active", false).Error

	case UserActionDelete:
		return db.WithContext(ctx).Delete(&User{}, userID).Error // soft delete

	case UserActionRestore:
		return db.Model(&User{}).
			Unscoped().
			Where("id = ?", userID).
			Update("deleted_at", nil).Error

	default:
		return fmt.Errorf("unsupported action: %s", action)
	}
}


// ---------------- Password Reset ----------------

var ErrInvalidOrExpiredToken = errors.New("invalid or expired reset token")

// Save
func (r *userRepo) SavePasswordResetToken(ctx context.Context, tx *gorm.DB, userID uint, tokenHash string, expiresAt time.Time) error {
    resetToken := PasswordResetToken{
        UserID:    userID,
        TokenHash: tokenHash,
        ExpiresAt: expiresAt,
    }
    return tx.WithContext(ctx).Create(&resetToken).Error
}

// Verify
func (r *userRepo) VerifyPasswordResetToken(ctx context.Context, tx *gorm.DB, token string) (uint, error) {
    var resetToken PasswordResetToken
    err := tx.WithContext(ctx).
        Where("expires_at > ?", time.Now()).
        First(&resetToken).Error
    if err != nil {
        return 0, ErrInvalidOrExpiredToken
    }

    if !utils.CheckPasswordHash(token, resetToken.TokenHash) {
        return 0, ErrInvalidOrExpiredToken
    }

    return resetToken.UserID, nil
}

// Delete
func (r *userRepo) DeletePasswordResetToken(ctx context.Context, tx *gorm.DB, token string) error {
    var resetToken PasswordResetToken
    err := tx.WithContext(ctx).First(&resetToken).Error
    if err != nil {
        return ErrInvalidOrExpiredToken
    }

    return tx.WithContext(ctx).Delete(&resetToken).Error
}