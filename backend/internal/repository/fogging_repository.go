package repository

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FoggingRepository interface {
	CreateFoggingRequest(ctx context.Context, tx *gorm.DB, request *domain.FoggingRequest) error
	UpdateFoggingRequestStatus(ctx context.Context, tx *gorm.DB, requestID uuid.UUID, status, reason string) error
	WithTransaction(fn func(tx *gorm.DB) error) error
	GetFoggingRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.FoggingRequest, error)
	UpdateFoggingRequestReservedStock(ctx context.Context, tx *gorm.DB, requestID uuid.UUID, reservedStock float64) error
}

type foggingRepository struct {
	db *gorm.DB
}

func NewFoggingRepository(db *gorm.DB) FoggingRepository {
	return &foggingRepository{db: db}
}

func (r *foggingRepository) CreateFoggingRequest(ctx context.Context, tx *gorm.DB, request *domain.FoggingRequest) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Create(request).Error
}

func (r *foggingRepository) UpdateFoggingRequestStatus(ctx context.Context, tx *gorm.DB, requestID uuid.UUID, status, reason string) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Model(&domain.FoggingRequest{}).Where("request_id = ?", requestID).Updates(map[string]interface{}{
		"approval_status": status,
		"reason":          reason,
	}).Error
}

func (r *foggingRepository) UpdateFoggingRequestReservedStock(ctx context.Context, tx *gorm.DB, requestID uuid.UUID, reservedStock float64) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Model(&domain.FoggingRequest{}).Where("request_id = ?", requestID).Update("reserved_stock", reservedStock).Error
}

func (r *foggingRepository) GetFoggingRequestByID(ctx context.Context, requestID uuid.UUID) (*domain.FoggingRequest, error) {
	var req domain.FoggingRequest
	if err := r.db.WithContext(ctx).Preload("PE").Where("request_id = ?", requestID).First(&req).Error; err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *foggingRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
