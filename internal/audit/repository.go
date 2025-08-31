package audit

import (
	"context"
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type AuditModel struct {
	ID         uint   `gorm:"primaryKey"`
	Timestamp  int64  `gorm:"autoCreateTime:milli"`
	ActorID    uint
	ActorName  string
	Action     string
	Entity     string
	EntityID   uint
	Reason     string
	FullData   string `gorm:"type:jsonb"`
	FieldDiffs string `gorm:"type:jsonb"`
}

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, entry AuditEntry) error {
	fullData, _ := json.Marshal(entry.FullData)
	diffs, _ := json.Marshal(entry.FieldDiffs)

	model := AuditModel{
		Timestamp:  entry.Timestamp.UnixMilli(),
		ActorID:    entry.ActorID,
		ActorName:  entry.ActorName,
		Action:     entry.Action,
		Entity:     entry.Entity,
		EntityID:   entry.EntityID,
		Reason:     entry.Reason,
		FullData:   string(fullData),
		FieldDiffs: string(diffs),
	}
	return r.db.WithContext(ctx).Create(&model).Error
}
type AuditFilters struct {
	ActorID *uint
	Action  string
	From    time.Time
	To      time.Time
	Limit   int
	Offset  int
}


func (r *Repository) List(ctx context.Context, entity string, entityID uint, filters AuditFilters) ([]AuditEntry, error) {
	var models []AuditModel
	tx := r.db.WithContext(ctx).Where("entity = ? AND entity_id = ?", entity, entityID)

	if filters.ActorID != nil {
		tx = tx.Where("actor_id = ?", *filters.ActorID)
	}
	if filters.Action != "" {
		tx = tx.Where("action = ?", filters.Action)
	}
	if !filters.From.IsZero() {
		tx = tx.Where("timestamp >= ?", filters.From.UnixMilli())
	}
	if !filters.To.IsZero() {
		tx = tx.Where("timestamp <= ?", filters.To.UnixMilli())
	}

	if filters.Limit == 0 {
		filters.Limit = 50 // default
	}
	if filters.Offset == 0 {
		filters.Offset = 0
	}

	if err := tx.Order("timestamp desc").Limit(filters.Limit).Offset(filters.Offset).Find(&models).Error; err != nil {
		return nil, err
	}

	var entries []AuditEntry
	for _, m := range models {
		var diffs []FieldChange
		_ = json.Unmarshal([]byte(m.FieldDiffs), &diffs)
		entries = append(entries, AuditEntry{
			Timestamp:  time.UnixMilli(m.Timestamp),
			ActorID:    m.ActorID,
			ActorName:  m.ActorName,
			Action:     m.Action,
			Entity:     m.Entity,
			EntityID:   m.EntityID,
			Reason:     m.Reason,
			FullData:   []byte(m.FullData),
			FieldDiffs: diffs,
		})
	}
	return entries, nil
}

