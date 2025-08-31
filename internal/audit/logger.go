package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Logger struct {
	storage Storage
}

func NewLogger(storage Storage) *Logger {
	return &Logger{storage: storage}
}

// LogSnapshot → CREATE / DELETE / RESTORE
func (l *Logger) LogSnapshot(ctx context.Context, action AuditAction, entity string, entityID uint, data interface{}, reason string) error {
	actorID, actorName := GetUser(ctx)
	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("audit marshal error: %w", err)
	}
	entry := AuditEntry{
		Timestamp: time.Now(),
		ActorID:   actorID,
		ActorName: actorName,
		Action:    string(action),
		Entity:    entity,
		EntityID:  entityID,
		Reason:    reason,
		FullData:  raw,
	}
	return l.storage.Save(ctx, entry)
}

// LogDiff → UPDATE with field-level diffs
func (l *Logger) LogDiff(ctx context.Context, action AuditAction, entity string, entityID uint, before, after map[string]interface{}, reason string) error {
	actorID, actorName := GetUser(ctx)
	diffs := CompareSnapshots(before, after)
	raw, err := json.Marshal(after)
	if err != nil {
		return fmt.Errorf("audit marshal error: %w", err)
	}
	entry := AuditEntry{
		Timestamp:  time.Now(),
		ActorID:    actorID,
		ActorName:  actorName,
		Action:     string(action),
		Entity:     entity,
		EntityID:   entityID,
		Reason:     reason,
		FullData:   raw,
		FieldDiffs: diffs,
	}
	return l.storage.Save(ctx, entry)
}

// LogEvent → VIEW / SEARCH / LOGIN / LOGOUT / EXPORT
func (l *Logger) LogEvent(ctx context.Context, action AuditAction, entity string, entityID uint, reason string) error {
	actorID, actorName := GetUser(ctx)
	entry := AuditEntry{
		Timestamp: time.Now(),
		ActorID:   actorID,
		ActorName: actorName,
		Action:    string(action),
		Entity:    entity,
		EntityID:  entityID,
		Reason:    reason,
	}
	return l.storage.Save(ctx, entry)
}
