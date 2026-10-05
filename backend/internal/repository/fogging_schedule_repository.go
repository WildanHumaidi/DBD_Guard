package repository

import (
	"context"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FoggingScheduleRepository interface {
	CreateSchedule(ctx context.Context, schedule *domain.FoggingSchedule) error
	CreateScheduleWithTx(ctx context.Context, tx *gorm.DB, schedule *domain.FoggingSchedule) error
	GetScheduleByID(ctx context.Context, id uuid.UUID) (*domain.FoggingSchedule, error)
	UpdateExecutionStatus(ctx context.Context, tx *gorm.DB, scheduleID uuid.UUID, status, notes string) error
	GetActiveScheduleByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.FoggingSchedule, error)
	CountSchedulesByTeamAndDate(ctx context.Context, teamID uuid.UUID, date time.Time) (int64, error)
	SaveEvidence(ctx context.Context, evidence *domain.FoggingEvidence) error
	WithTransaction(fn func(tx *gorm.DB) error) error
}

type foggingScheduleRepository struct {
	db *gorm.DB
}

func NewFoggingScheduleRepository(db *gorm.DB) FoggingScheduleRepository {
	return &foggingScheduleRepository{db: db}
}

func (r *foggingScheduleRepository) CreateSchedule(ctx context.Context, schedule *domain.FoggingSchedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *foggingScheduleRepository) CreateScheduleWithTx(ctx context.Context, tx *gorm.DB, schedule *domain.FoggingSchedule) error {
	if tx == nil {
		return r.CreateSchedule(ctx, schedule)
	}
	return tx.WithContext(ctx).Create(schedule).Error
}

func (r *foggingScheduleRepository) GetActiveScheduleByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.FoggingSchedule, error) {
	var schedule domain.FoggingSchedule
	err := r.db.WithContext(ctx).Where("request_id = ? AND execution_status IN (?, ?, ?)", requestID, "SCHEDULED", "ON_THE_WAY", "IN_PROGRESS").First(&schedule).Error
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *foggingScheduleRepository) GetScheduleByID(ctx context.Context, id uuid.UUID) (*domain.FoggingSchedule, error) {
	var schedule domain.FoggingSchedule
	err := r.db.WithContext(ctx).Preload("Request.PE").First(&schedule, "schedule_id = ?", id).Error
	return &schedule, err
}

func (r *foggingScheduleRepository) UpdateExecutionStatus(ctx context.Context, tx *gorm.DB, scheduleID uuid.UUID, status, notes string) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}

	return execDB.WithContext(ctx).
		Model(&domain.FoggingSchedule{}).
		Where("schedule_id = ?", scheduleID).
		Updates(map[string]interface{}{
			"execution_status": status,
			"notes":            notes,
		}).Error
}

func (r *foggingScheduleRepository) CountSchedulesByTeamAndDate(ctx context.Context, teamID uuid.UUID, date time.Time) (int64, error) {
	var count int64
	// date truncation for Postgres: DATE(planned_time)
	// Truncate parameter to day
	dateStr := date.Format("2006-01-02")
	err := r.db.WithContext(ctx).
		Model(&domain.FoggingSchedule{}).
		Where("team_id = ? AND DATE(planned_time) = ? AND execution_status != 'FAILED'", teamID, dateStr).
		Count(&count).Error
	return count, err
}

func (r *foggingScheduleRepository) SaveEvidence(ctx context.Context, evidence *domain.FoggingEvidence) error {
	return r.db.WithContext(ctx).Create(evidence).Error
}

func (r *foggingScheduleRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
