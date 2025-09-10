package clinics

import (
	"context"

	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

func (r *GormRepository) GetClinicByID(ctx context.Context, tx *gorm.DB, id uint, includeDeleted bool) (*Clinic, error) {
	var clinic Clinic
	query := tx.WithContext(ctx)
	if !includeDeleted {
		query = query.Where("deleted_at IS NULL")
	}
	if err := query.First(&clinic, id).Error; err != nil {
		return nil, err
	}
	return &clinic, nil
}

func (r *GormRepository) ListClinics(ctx context.Context, query *ListClinicsQuery) ([]*Clinic, int64, error) {
	var clinics []*Clinic
	var total int64

	dbQuery := r.db.WithContext(ctx).Model(&Clinic{})
	if !query.IncludeDeleted {
		dbQuery = dbQuery.Where("deleted_at IS NULL")
	}
	if query.Search != "" {
		search := "%" + query.Search + "%"
		dbQuery = dbQuery.Where("name ILIKE ?", search)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.PageSize > 0 {
		offset := (query.Page - 1) * query.PageSize
		dbQuery = dbQuery.Offset(offset).Limit(query.PageSize)
	}

	if err := dbQuery.Find(&clinics).Error; err != nil {
		return nil, 0, err
	}
	return clinics, total, nil
}

func (r *GormRepository) ListActiveUserClinics(ctx context.Context, userID uint) ([]UserClinic, error) {
	var result []UserClinic
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *GormRepository) RevokeUserClinic(ctx context.Context, uc *UserClinic) error {
	return r.db.WithContext(ctx).Model(&UserClinic{}).
		Where("user_id = ? AND clinic_id = ? AND revoked_at IS NULL", uc.UserID, uc.ClinicID).
		Update("revoked_at", uc.RevokedAt).Error
}

func (r *GormRepository) RevokeClinicsForUser(ctx context.Context, tx *gorm.DB, userID uint) error {
	return tx.WithContext(ctx).
		Model(&UserClinic{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", gorm.Expr("NOW()")).Error
}
