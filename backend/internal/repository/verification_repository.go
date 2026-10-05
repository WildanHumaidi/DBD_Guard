package repository

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"gorm.io/gorm"
)

type VerificationRepository interface {
	Create(ctx context.Context, tx *gorm.DB, verification *domain.Verification) error
	WithTransaction(fn func(tx *gorm.DB) error) error
}

type verificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) VerificationRepository {
	return &verificationRepository{db: db}
}

func (r *verificationRepository) Create(ctx context.Context, tx *gorm.DB, verification *domain.Verification) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Create(verification).Error
}

func (r *verificationRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
