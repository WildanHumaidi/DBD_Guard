package usecase

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
)

type AuditUsecase interface {
	GetRecentLogs(ctx context.Context, limit int) ([]domain.AuditLog, error)
}

type auditUsecase struct {
	auditRepo repository.AuditRepository
}

func NewAuditUsecase(auditRepo repository.AuditRepository) AuditUsecase {
	return &auditUsecase{auditRepo}
}

func (u *auditUsecase) GetRecentLogs(ctx context.Context, limit int) ([]domain.AuditLog, error) {
	if limit <= 0 {
		limit = 10 // default limit
	}
	if limit > 50 {
		limit = 50 // max limit
	}
	return u.auditRepo.GetRecentLogs(ctx, limit)
}
