package repository

import (
	"context"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PERepository interface {
	CreatePEResult(ctx context.Context, tx *gorm.DB, pe *domain.PEResult) error
	GetPEResultByID(ctx context.Context, tx *gorm.DB, peID uuid.UUID) (*domain.PEResult, error)
	GetPEResultByReportID(ctx context.Context, reportID uuid.UUID) (*domain.PEResult, error)
	UpdatePEResult(ctx context.Context, tx *gorm.DB, pe *domain.PEResult) error
	CountPositivePEByKelurahan(ctx context.Context, tx *gorm.DB, kelurahan string, since time.Time) (int64, error)
	CreateRiskScore(ctx context.Context, tx *gorm.DB, riskScore *domain.RiskScore) error
	WithTransaction(fn func(tx *gorm.DB) error) error
}

type peRepository struct {
	db *gorm.DB
}

func NewPERepository(db *gorm.DB) PERepository {
	return &peRepository{db: db}
}

func (r *peRepository) GetPEResultByReportID(ctx context.Context, reportID uuid.UUID) (*domain.PEResult, error) {
	var pe domain.PEResult
	err := r.db.WithContext(ctx).First(&pe, "report_id = ?", reportID).Error
	if err != nil {
		return nil, err
	}
	return &pe, nil
}

func (r *peRepository) GetPEResultByID(ctx context.Context, tx *gorm.DB, peID uuid.UUID) (*domain.PEResult, error) {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	var pe domain.PEResult
	err := execDB.WithContext(ctx).Preload("Report").Preload("Report.Location").First(&pe, "pe_id = ?", peID).Error
	if err != nil {
		return nil, err
	}
	return &pe, nil
}

func (r *peRepository) CreatePEResult(ctx context.Context, tx *gorm.DB, pe *domain.PEResult) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Create(pe).Error
}

func (r *peRepository) UpdatePEResult(ctx context.Context, tx *gorm.DB, pe *domain.PEResult) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Save(pe).Error
}

func (r *peRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *peRepository) CountPositivePEByKelurahan(ctx context.Context, tx *gorm.DB, kelurahan string, since time.Time) (int64, error) {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	var count int64
	err := execDB.WithContext(ctx).
		Table("pe_results").
		Joins("JOIN citizen_reports ON pe_results.report_id = citizen_reports.report_id").
		Joins("JOIN report_locations ON citizen_reports.report_id = report_locations.report_id").
		Where("pe_results.recommendation = ?", "FOGGING"). // Assuming FOGGING = Positive case requiring action
		Where("pe_results.created_at >= ?", since).
		Where("report_locations.kelurahan = ?", kelurahan).
		Count(&count).Error
	return count, err
}

func (r *peRepository) CreateRiskScore(ctx context.Context, tx *gorm.DB, riskScore *domain.RiskScore) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Create(riskScore).Error
}
