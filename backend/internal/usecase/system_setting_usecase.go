package usecase

import (
	"context"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/google/uuid"
)

type SystemSettingUsecase interface {
	GetSettings(ctx context.Context) (*domain.SystemSetting, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, heroTitle, announcementText string, isCitizenPortalActive bool, bannerImage *multipart.FileHeader) error
}

type systemSettingUsecase struct {
	settingRepo repository.SystemSettingRepository
	auditRepo   repository.AuditRepository
}

func NewSystemSettingUsecase(settingRepo repository.SystemSettingRepository, auditRepo repository.AuditRepository) SystemSettingUsecase {
	return &systemSettingUsecase{settingRepo, auditRepo}
}

func (u *systemSettingUsecase) GetSettings(ctx context.Context) (*domain.SystemSetting, error) {
	return u.settingRepo.GetSettings(ctx)
}

func (u *systemSettingUsecase) UpdateSettings(ctx context.Context, userID uuid.UUID, heroTitle, announcementText string, isCitizenPortalActive bool, bannerImage *multipart.FileHeader) error {
	settings, err := u.settingRepo.GetSettings(ctx)
	if err != nil {
		return err
	}

	settings.HeroTitle = heroTitle
	settings.AnnouncementText = announcementText
	settings.IsCitizenPortalActive = isCitizenPortalActive
	settings.UpdatedAt = time.Now()

	if bannerImage != nil {
		path, err := u.saveBannerImage(bannerImage)
		if err != nil {
			return err
		}
		// Optional: delete old banner image to save space
		if settings.BannerImagePath != "" {
			oldPath := strings.TrimPrefix(settings.BannerImagePath, "/")
			_ = os.Remove(oldPath)
		}
		// Assume our static route serves /storage -> ./storage
		settings.BannerImagePath = "/" + path
	}

	err = u.settingRepo.UpdateSettings(ctx, settings)
	if err != nil {
		return err
	}

	// Log audit
	_ = u.auditRepo.LogAction(ctx, nil, userID, "UPDATE_SYSTEM_SETTINGS", "SYSTEM_SETTING", settings.ID, map[string]interface{}{
		"is_citizen_portal_active": isCitizenPortalActive,
	})

	return nil
}

func (u *systemSettingUsecase) saveBannerImage(fileHeader *multipart.FileHeader) (string, error) {
	// Create directory if not exists
	storageDir := "storage/banners"
	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		return "", err
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return "", errors.New("hanya file jpg dan png yang diizinkan")
	}

	// Check size (max 5MB)
	if fileHeader.Size > 5*1024*1024 {
		return "", errors.New("ukuran file maksimal 5MB")
	}

	// Generate unique filename
	filename := uuid.New().String() + ext
	path := filepath.Join(storageDir, filename)

	// Save file
	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	// Convert windows slashes to forward slashes for URL
	return filepath.ToSlash(path), nil
}
