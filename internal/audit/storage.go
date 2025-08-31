package audit

import (
	"context"

	"gorm.io/gorm"
)

// DBStorage implements Storage interface using GORM
type DBStorage struct {
	db *gorm.DB
}

// NewDBStorage creates a new DBStorage
func NewDBStorage(db *gorm.DB) *DBStorage {
	return &DBStorage{db: db}
}

// Save writes an audit entry to the database
func (s *DBStorage) Save(ctx context.Context, entry AuditEntry) error {
	return s.db.WithContext(ctx).Create(&entry).Error
}

// List fetches audit entries for a specific entity and entityID
func (s *DBStorage) List(ctx context.Context, entity string, entityID uint) ([]AuditEntry, error) {
	var entries []AuditEntry
	if err := s.db.WithContext(ctx).
		Where("entity = ? AND entity_id = ?", entity, entityID).
		Order("timestamp DESC").
		Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}

// Optional: define a GORM model for AuditEntry if needed
func (AuditEntry) TableName() string {
	return "audit_entries"
}
