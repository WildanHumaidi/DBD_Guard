package usecase

import (
	"context"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
)

type EwsKlbUsecase interface {
	GetConfiguration(ctx context.Context) (*domain.KlbConfiguration, error)
	UpdateConfiguration(ctx context.Context, req UpdateKlbConfigRequest) error
	CheckAndTriggerKLB(ctx context.Context, areaName string, areaType string) error
	GetActiveKLBAlerts(ctx context.Context) ([]domain.KlbAlert, error)
	GetEwsForRw(ctx context.Context, kelurahan, rw string) ([]domain.CitizenReport, error)
}

type UpdateKlbConfigRequest struct {
	RegionScope   string `json:"region_scope"`
	MinCases      int    `json:"min_cases"`
	TimeframeDays int    `json:"timeframe_days"`
	IsAutoTrigger bool   `json:"is_auto_trigger"`
}

type ewsKlbUsecase struct {
	klbRepo    repository.KlbRepository
	reportRepo repository.ReportRepository
	notifUsecase NotificationUsecase
}

func NewEwsKlbUsecase(klbRepo repository.KlbRepository, reportRepo repository.ReportRepository, notifUsecase NotificationUsecase) EwsKlbUsecase {
	return &ewsKlbUsecase{
		klbRepo:      klbRepo,
		reportRepo:   reportRepo,
		notifUsecase: notifUsecase,
	}
}

func (u *ewsKlbUsecase) GetConfiguration(ctx context.Context) (*domain.KlbConfiguration, error) {
	return u.klbRepo.GetConfiguration(ctx)
}

func (u *ewsKlbUsecase) UpdateConfiguration(ctx context.Context, req UpdateKlbConfigRequest) error {
	config, err := u.klbRepo.GetConfiguration(ctx)
	if err != nil {
		return err
	}
	config.RegionScope = req.RegionScope
	config.MinCases = req.MinCases
	config.TimeframeDays = req.TimeframeDays
	config.IsAutoTrigger = req.IsAutoTrigger

	return u.klbRepo.SaveConfiguration(ctx, config)
}

func (u *ewsKlbUsecase) CheckAndTriggerKLB(ctx context.Context, areaName string, areaType string) error {
	config, err := u.klbRepo.GetConfiguration(ctx)
	if err != nil {
		return err
	}

	if !config.IsAutoTrigger {
		return nil
	}

	// Make sure we only check the area that matches the config scope
	if areaType != config.RegionScope {
		return nil
	}

	// 1. Count cases in this area
	casesCount, err := u.reportRepo.CountCasesByArea(ctx, areaName, areaType, config.TimeframeDays)
	if err != nil {
		return err
	}

	// 2. Check threshold
	if casesCount >= config.MinCases {
		// Ensure no active alert already exists for this exact area
		existingAlert, err := u.klbRepo.GetActiveAlertByArea(ctx, areaName)
		if err != nil {
			return err
		}
		
		if existingAlert == nil {
			// Trigger KLB
			newAlert := &domain.KlbAlert{
				AreaName:   areaName,
				AreaType:   areaType,
				CasesCount: casesCount,
				IsActive:   true,
			}
			err := u.klbRepo.CreateAlert(ctx, newAlert)
			if err == nil && u.notifUsecase != nil {
				_ = u.notifUsecase.NotifyKLBTriggered(ctx, areaName, areaType, casesCount)
			}
			return err
		} else {
			// Optionally update the count on the existing alert?
			// But for now, we just let it be active.
		}
	} else {
		// If cases dropped below threshold, maybe resolve the alert automatically?
		// PRD says "KLB Otomatis". If we drop below, it should turn off.
		existingAlert, err := u.klbRepo.GetActiveAlertByArea(ctx, areaName)
		if err == nil && existingAlert != nil {
			return u.klbRepo.ResolveAlert(ctx, existingAlert.ID.String())
		}
	}

	return nil
}

func (u *ewsKlbUsecase) GetActiveKLBAlerts(ctx context.Context) ([]domain.KlbAlert, error) {
	return u.klbRepo.GetActiveAlerts(ctx)
}

func (u *ewsKlbUsecase) GetEwsForRw(ctx context.Context, kelurahan, rw string) ([]domain.CitizenReport, error) {
	// Look back 7 days by default for EWS
	return u.reportRepo.GetRecentCasesInSameKelurahanDifferentRW(ctx, kelurahan, rw, 7)
}
