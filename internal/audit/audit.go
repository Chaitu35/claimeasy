package audit

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// ---------------- Context for Actor ----------------
type ctxKey string

const (
	userIDKey   ctxKey = "audit_user_id"
	userNameKey ctxKey = "audit_user_name"
)

func WithUser(ctx context.Context, id uint, name string) context.Context {
	ctx = context.WithValue(ctx, userIDKey, id)
	ctx = context.WithValue(ctx, userNameKey, name)
	return ctx
}

func GetUser(ctx context.Context) (uint, string) {
	id, _ := ctx.Value(userIDKey).(uint)
	name, _ := ctx.Value(userNameKey).(string)
	return id, name
}

// ---------------- Audit Actions ----------------
type AuditAction string

const (
	ActionCreate   AuditAction = "CREATE"
	ActionUpdate   AuditAction = "UPDATE"
	ActionDelete   AuditAction = "DELETE"
	ActionRestore  AuditAction = "RESTORE"
	ActionLogin    AuditAction = "LOGIN"
	ActionLogout   AuditAction = "LOGOUT"
	ActionView     AuditAction = "VIEW"
	ActionSearch   AuditAction = "SEARCH"
	ActionExport   AuditAction = "EXPORT"
	ActionSuspend  AuditAction = "SUSPEND"
	ActionActivate AuditAction = "ACTIVATE"
	ActionArchive  AuditAction = "ARCHIVE"
)

// ---------------- AuditEntry / FieldChange ----------------
type AuditEntry struct {
	ID         uint           `gorm:"primaryKey"`
	Timestamp  time.Time      `json:"timestamp"`
	ActorID    uint           `json:"actor_id"`
	ActorName  string         `json:"actor_name"`
	Action     string         `json:"action"`
	Entity     string         `json:"entity"`
	EntityID   uint           `json:"entity_id"`
	Reason     string         `json:"reason,omitempty"`
	FullData   []byte         `json:"full_data"`
	FieldDiffs JSONFieldDiffs `gorm:"type:jsonb" json:"field_diffs"` // store slice as JSON

	FullDataStr string `json:"full_data_str"`
}

// FieldChange represents a single field update
type FieldChange struct {
	Field    string `json:"field"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

// JSONFieldDiffs allows storing []FieldChange as JSON in Postgres
type JSONFieldDiffs []FieldChange

func (j JSONFieldDiffs) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONFieldDiffs) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONFieldDiffs value: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

// ---------------- Storage Interface ----------------
type Storage interface {
	Save(ctx context.Context, entry AuditEntry) error
	List(ctx context.Context, entity string, entityID uint) ([]AuditEntry, error)
}

// ---------------- Utilities ----------------
func CompareSnapshots(prev, curr map[string]interface{}) JSONFieldDiffs {
	var changes JSONFieldDiffs
	for key, newVal := range curr {
		oldVal := prev[key]
		if fmt.Sprintf("%v", oldVal) != fmt.Sprintf("%v", newVal) {
			changes = append(changes, FieldChange{
				Field:    key,
				OldValue: fmt.Sprintf("%v", oldVal),
				NewValue: fmt.Sprintf("%v", newVal),
			})
		}
	}
	return changes
}
