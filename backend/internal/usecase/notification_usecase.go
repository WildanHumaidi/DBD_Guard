package usecase

import (
	"context"

	"github.com/dbd-guard/backend/internal/infrastructure/notification"
	"github.com/dbd-guard/backend/internal/repository"
)

type NotificationUsecase interface {
	NotifyNewReport(ctx context.Context, reportID string, areaRT string, areaRW string, areaKelurahan string) error
	NotifyKLBTriggered(ctx context.Context, areaName string, areaType string, casesCount int) error
	NotifyFoggingScheduled(ctx context.Context, reportID string, targetDate string) error
	NotifyEWSKader(ctx context.Context, kelurahan string, excludeActorID string) error
}

type notificationUsecase struct {
	notifService notification.Service
	userRepo     repository.UserRepository
}

func NewNotificationUsecase(ns notification.Service, ur repository.UserRepository) NotificationUsecase {
	return &notificationUsecase{
		notifService: ns,
		userRepo:     ur,
	}
}

// Dummy methods for now that uses the mock service
func (u *notificationUsecase) NotifyNewReport(ctx context.Context, reportID string, areaRT string, areaRW string, areaKelurahan string) error {
	// Ideally we find all kader in this area
	return u.notifService.SendPushNotification(ctx, "kader-token-123", "Laporan Baru", "Ada laporan baru di RT "+areaRT+" RW "+areaRW, nil)
}

func (u *notificationUsecase) NotifyKLBTriggered(ctx context.Context, areaName string, areaType string, casesCount int) error {
	return u.notifService.SendPushNotification(ctx, "dinkes-puskesmas-token", "Peringatan KLB", "Status KLB aktif di "+areaType+" "+areaName, nil)
}

func (u *notificationUsecase) NotifyFoggingScheduled(ctx context.Context, reportID string, targetDate string) error {
	return u.notifService.SendPushNotification(ctx, "reporter-token-123", "Fogging Terjadwal", "Fogging dijadwalkan pada "+targetDate, nil)
}

func (u *notificationUsecase) NotifyEWSKader(ctx context.Context, kelurahan string, excludeActorID string) error {
	// Peringatan EWS (Kader Terdekat): Saat sebuah laporan divalidasi oleh Kader (status VERIFIED),
	// sistem mencari Kader lain dalam satu kelurahan.
	// We'll mock the send for now as per MVP strategy.
	msg := "Kasus baru terkonfirmasi di wilayah Anda. Mohon tingkatkan kewaspadaan dan lakukan PSN."
	return u.notifService.SendPushNotification(ctx, "kader-kelurahan-"+kelurahan, "Peringatan EWS Kader", msg, nil)
}
