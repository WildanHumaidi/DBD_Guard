package repository

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"gorm.io/gorm"
)

type KlbRepository interface {
	GetConfiguration(ctx context.Context) (*domain.KlbConfiguration, error)
	SaveConfiguration(ctx context.Context, config *domain.KlbConfiguration) error
	CreateAlert(ctx context.Context, alert *domain.KlbAlert) error
	GetActiveAlerts(ctx context.Context) ([]domain.KlbAlert, error)
	ResolveAlert(ctx context.Context, alertID string) error
	GetActiveAlertByArea(ctx context.Context, areaName string) (*domain.KlbAlert, error)
}

type klbRepository struct {
	db *gorm.DB
}

func NewKlbRepository(db *gorm.DB) KlbRepository {
	return &klbRepository{db: db}
}

func (r *klbRepository) GetConfiguration(ctx context.Context) (*domain.KlbConfiguration, error) {
	var config domain.KlbConfiguration
	err := r.db.WithContext(ctx).First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a default configuration if it doesn't exist
			defaultConfig := domain.KlbConfiguration{
				RegionScope:   "KELURAHAN",
				MinCases:      5,
				TimeframeDays: 7,
				IsAutoTrigger: true,
			}
			if err := r.db.WithContext(ctx).Create(&defaultConfig).Error; err != nil {
				return nil, err
			}
			return &defaultConfig, nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *klbRepository) SaveConfiguration(ctx context.Context, config *domain.KlbConfiguration) error {
	return r.db.WithContext(ctx).Save(config).Error
}

func (r *klbRepository) CreateAlert(ctx context.Context, alert *domain.KlbAlert) error {
	return r.db.WithContext(ctx).Create(alert).Error
}

func (r *klbRepository) GetActiveAlerts(ctx context.Context) ([]domain.KlbAlert, error) {
	var alerts []domain.KlbAlert
	err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&alerts).Error
	return alerts, err
}

func (r *klbRepository) ResolveAlert(ctx context.Context, alertID string) error {
	return r.db.WithContext(ctx).Model(&domain.KlbAlert{}).Where("alert_id = ?", alertID).Updates(map[string]interface{}{
		"is_active":   false,
		"resolved_at": gorm.Expr("CURRENT_TIMESTAMP"),
	}).Error
}

func (r *klbRepository) GetActiveAlertByArea(ctx context.Context, areaName string) (*domain.KlbAlert, error) {
	var alert domain.KlbAlert
	err := r.db.WithContext(ctx).Where("area_name = ? AND is_active = ?", areaName, true).First(&alert).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No active alert
		}
		return nil, err
	}
	return &alert, nil
}
