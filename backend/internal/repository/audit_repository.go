package repository

import (
	"context"
	"encoding/json"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AuditRepository interface {
	LogAction(ctx context.Context, tx *gorm.DB, userID uuid.UUID, action, entityType string, entityID uuid.UUID, changes interface{}) error
	GetRecentLogs(ctx context.Context, limit int) ([]domain.AuditLog, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db: db}
}

// LogAction inserts an audit log using either an ongoing transaction (tx) or the default DB instance
func (r *auditRepository) LogAction(ctx context.Context, tx *gorm.DB, userID uuid.UUID, action, entityType string, entityID uuid.UUID, changes interface{}) error {
	changesBytes, _ := json.Marshal(changes)

	log := domain.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		Changes:    datatypes.JSON(changesBytes),
	}

	execDB := r.db
	if tx != nil {
		execDB = tx
	}

	return execDB.WithContext(ctx).Create(&log).Error
}

// GetRecentLogs retrieves the most recent audit logs with their associated User, limited by the 'limit' parameter
func (r *auditRepository) GetRecentLogs(ctx context.Context, limit int) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.WithContext(ctx).
		Preload("User").
		Order("created_at desc").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}
