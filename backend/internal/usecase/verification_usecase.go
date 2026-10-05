package usecase

import (
	"context"
	"errors"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationUsecase interface {
	VerifyReport(ctx context.Context, reportID, actorID uuid.UUID, actorRole string, locationID *uuid.UUID, decision, reason string) error
}

type verificationUsecase struct {
	verificationRepo repository.VerificationRepository
	reportRepo       repository.ReportRepository
	auditRepo        repository.AuditRepository
	ewsKlbUsecase    EwsKlbUsecase
	notifUsecase     NotificationUsecase
}

func NewVerificationUsecase(vr repository.VerificationRepository, rr repository.ReportRepository, ar repository.AuditRepository, ews EwsKlbUsecase, nu NotificationUsecase) VerificationUsecase {
	return &verificationUsecase{
		verificationRepo: vr,
		reportRepo:       rr,
		auditRepo:        ar,
		ewsKlbUsecase:    ews,
		notifUsecase:     nu,
	}
}

func (u *verificationUsecase) VerifyReport(ctx context.Context, reportID, actorID uuid.UUID, actorRole string, locationID *uuid.UUID, decision, reason string) error {
	// Rule 4: If rejected, reason is mandatory
	if decision == "REJECTED" && reason == "" {
		return errors.New("alasan penolakan wajib diisi")
	}

	// 1. Get the report
	report, err := u.reportRepo.GetReportByID(ctx, reportID)
	if err != nil {
		return err
	}

	// 1.5 Role-Based Business Logic Validation
	if actorRole == domain.RolePetugasPuskesmas || actorRole == domain.RoleKepalaPuskesmas {
		if report.Status != "APPEALED" {
			return errors.New("Puskesmas hanya dapat memverifikasi laporan dengan status APPEALED")
		}
		// Optional: Location scope check for Puskesmas
		if locationID == nil {
			return errors.New("Puskesmas tidak memiliki area tugas yang valid")
		}
	} else if actorRole == domain.RoleKader {
		if report.Status != "MENUNGGU_PE_KADER" {
			return errors.New("Kader hanya dapat memverifikasi/menolak laporan dengan status MENUNGGU_PE_KADER")
		}
		if decision != "REJECTED" {
			return errors.New("Kader hanya diizinkan menggunakan endpoint ini untuk menolak (HOAX) laporan. Gunakan endpoint PE untuk menyetujui.")
		}
	} else {
		return errors.New("role anda tidak memiliki izin untuk memverifikasi laporan ini")
	}

	// 2. Perform Transaction
	err = u.verificationRepo.WithTransaction(func(tx *gorm.DB) error {
		// a. Create Verification record
		verification := &domain.Verification{
			ReportID: reportID,
			ActorID:  actorID,
			Decision: decision,
			Reason:   reason,
		}

		if err := u.verificationRepo.Create(ctx, tx, verification); err != nil {
			return err
		}

		// b. Update Report Status
		newStatus := "VERIFIED"
		if actorRole == domain.RolePetugasPuskesmas || actorRole == domain.RoleKepalaPuskesmas {
			if decision == "REJECTED" {
				newStatus = "REJECTED"
			} else {
				newStatus = "ESKALASI_PUSKESMAS"
			}
		} else if actorRole == domain.RoleKader {
			if decision == "REJECTED" {
				newStatus = "REJECTED_HOAX" // Special status for kader hoax rejection
			}
		}
		if err := u.reportRepo.UpdateReportStatus(ctx, tx, reportID, newStatus); err != nil {
			return err
		}

		// c. Log action to Audit Logs (Rule 6)
		auditDetails := map[string]interface{}{
			"decision": decision,
			"reason":   reason,
		}
		if err := u.auditRepo.LogAction(ctx, tx, actorID, "VERIFY_REPORT", "CitizenReport", reportID, auditDetails); err != nil {
			return err
		}

		return nil
	})

	if err == nil && decision == "APPROVED" {
		if report.Location.Kelurahan != "" {
			if u.notifUsecase != nil {
				_ = u.notifUsecase.NotifyEWSKader(ctx, report.Location.Kelurahan, actorID.String())
			}
			if u.ewsKlbUsecase != nil {
				_ = u.ewsKlbUsecase.CheckAndTriggerKLB(ctx, report.Location.Kelurahan, "KELURAHAN")
			}
		}
		if report.Location.Kecamatan != "" && u.ewsKlbUsecase != nil {
			_ = u.ewsKlbUsecase.CheckAndTriggerKLB(ctx, report.Location.Kecamatan, "KECAMATAN")
		}
	}

	return err
}
