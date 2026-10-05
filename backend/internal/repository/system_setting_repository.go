package repository

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"gorm.io/gorm"
)

type SystemSettingRepository interface {
	GetSettings(ctx context.Context) (*domain.SystemSetting, error)
	UpdateSettings(ctx context.Context, settings *domain.SystemSetting) error
}

type systemSettingRepository struct {
	db *gorm.DB
}

func NewSystemSettingRepository(db *gorm.DB) SystemSettingRepository {
	return &systemSettingRepository{db}
}

func (r *systemSettingRepository) GetSettings(ctx context.Context) (*domain.SystemSetting, error) {
	var settings domain.SystemSetting
	err := r.db.WithContext(ctx).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *systemSettingRepository) UpdateSettings(ctx context.Context, settings *domain.SystemSetting) error {
	return r.db.WithContext(ctx).Save(settings).Error
}
