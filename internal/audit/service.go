package audit

import (
	"context"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListAuditLogs(ctx context.Context, entity string, entityID uint, filters AuditFilters) ([]AuditEntry, error) {
	return s.repo.List(ctx, entity, entityID, filters)
}

