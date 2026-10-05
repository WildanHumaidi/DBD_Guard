package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PEUsecase interface {
	SubmitPEResult(ctx context.Context, reportID, investigatorID uuid.UUID, peRadius int, larvaFound bool, feverCount int, psnEdukasiDone bool, notes string, photoUrl string) (*uuid.UUID, error)
	ValidatePEResult(ctx context.Context, peID, validatorID uuid.UUID, isApproved bool, reason string, overrideDecision string, overrideReason string) error
}

type peUsecase struct {
	peRepo        repository.PERepository
	reportRepo    repository.ReportRepository
	foggingRepo   repository.FoggingRepository
	inventoryRepo repository.InventoryRepository
	auditRepo     repository.AuditRepository
	masterDataRepo repository.MasterDataRepository
	ewsKlbUsecase EwsKlbUsecase
}

func NewPEUsecase(pr repository.PERepository, rr repository.ReportRepository, fr repository.FoggingRepository, ir repository.InventoryRepository, ar repository.AuditRepository, md repository.MasterDataRepository, ews EwsKlbUsecase) PEUsecase {
	return &peUsecase{
		peRepo:        pr,
		reportRepo:    rr,
		foggingRepo:   fr,
		inventoryRepo: ir,
		auditRepo:     ar,
		masterDataRepo: md,
		ewsKlbUsecase: ews,
	}
}

func (u *peUsecase) SubmitPEResult(ctx context.Context, reportID, investigatorID uuid.UUID, peRadius int, larvaFound bool, feverCount int, psnEdukasiDone bool, notes string, photoUrl string) (*uuid.UUID, error) {
	if !psnEdukasiDone {
		return nil, errors.New("edukasi 3M Plus dan Pemberantasan Sarang Nyamuk (PSN) wajib dilakukan")
	}

	// 1. Validate Report
	report, err := u.reportRepo.GetReportByID(ctx, reportID)
	if err != nil {
		return nil, err
	}

	// Allow MENUNGGU_PE_KADER, APPEALED, and ESKALASI_PUSKESMAS
	if report.Status != "MENUNGGU_PE_KADER" && report.Status != "APPEALED" && report.Status != "ESKALASI_PUSKESMAS" {
		return nil, errors.New("hanya laporan yang berstatus MENUNGGU_PE_KADER, APPEALED, atau ESKALASI_PUSKESMAS yang dapat disubmit hasil PE-nya")
	}

	// 1.1 Verify Ownership (Phase 0.6 Invariant)
	if report.InvestigatorID == nil || *report.InvestigatorID != investigatorID {
		return nil, errors.New("UNAUTHORIZED_PE_SUBMISSION: laporan ini tidak diklaim oleh Anda")
	}

	// 1.2 Check if PE result already exists
	_, err = u.peRepo.GetPEResultByReportID(ctx, reportID)
	if err == nil {
		return nil, errors.New("laporan ini sudah memiliki data Penyelidikan Epidemiologi")
	}

	// 1.2 Calculate Recommendation Automatically based on PRD Business Rule
	recommendation := "EDUKASI_PSN"
	if larvaFound || feverCount >= 3 {
		var puskesmasName string = report.Location.Kelurahan // Fallback
		if u.masterDataRepo != nil {
			kelurahan, err := u.masterDataRepo.GetKelurahanByName(ctx, report.Location.Kelurahan)
			if err == nil && kelurahan.Puskesmas.Name != "" {
				puskesmasName = kelurahan.Puskesmas.Name
			}
		}
		
		item, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, "INSEKTISIDA", puskesmasName)
		if err == nil && (item.Quantity-item.ReservedQuantity) >= 10.0 {
			recommendation = "FOGGING_FOKUS"
		} else {
			recommendation = "FOGGING_TERTUNDA"
		}
	}

	// 2. Perform Transaction
	var newPeID uuid.UUID
	err = u.peRepo.WithTransaction(func(tx *gorm.DB) error {
		// a. Create PE Result with Status PENDING_VALIDATION
		peResult := &domain.PEResult{
			ReportID:       reportID,
			OfficerID:      investigatorID, // Kader ID
			Status:         "PENDING_VALIDATION",
			InspectionDate: time.Now(),
			RadiusMeters:   peRadius,
			LarvaeFound:    larvaFound,
			FeverCount:     feverCount,
			Recommendation: recommendation,
			PsnEdukasiDone: psnEdukasiDone,
			Notes:          notes,
			EvidencePhotoURL: photoUrl,
		}

		if err := u.peRepo.CreatePEResult(ctx, tx, peResult); err != nil {
			return err
		}
		newPeID = peResult.ID

		// b. Log to Audit
		auditDetails := map[string]interface{}{
			"fever_count":    feverCount,
			"larvae_found":   larvaFound,
			"recommendation": recommendation,
		}
		if err := u.auditRepo.LogAction(ctx, tx, investigatorID, "SUBMIT_PE", "PEResult", peResult.ID, auditDetails); err != nil {
			return err
		}

		// c. Update Report Status
		if err := u.reportRepo.UpdateReportStatus(ctx, tx, reportID, "MENUNGGU_VALIDASI_PUSKESMAS"); err != nil {
			return err
		}

		return nil
	})

	return &newPeID, err
}

func (u *peUsecase) ValidatePEResult(ctx context.Context, peID, validatorID uuid.UUID, isApproved bool, reason string, overrideDecision string, overrideReason string) error {
	var peResult *domain.PEResult
	err := u.peRepo.WithTransaction(func(tx *gorm.DB) error {
		var err error
		peResult, err = u.peRepo.GetPEResultByID(ctx, tx, peID)
		if err != nil || peResult == nil {
			return errors.New("hasil PE tidak ditemukan")
		}

		if peResult.Status != "PENDING_VALIDATION" {
			return errors.New("hasil PE ini tidak dalam status menunggu validasi")
		}

		// P0-04: Segregation of Duties (Four-Eyes Principle)
		// The investigator who conducted/submitted the PE cannot validate their own work
		if peResult.OfficerID == validatorID {
			return errors.New("pelanggaran pemisahan tugas: pemeriksa PE dilarang memvalidasi laporannya sendiri (Four-Eyes Principle)")
		}

		if !isApproved {
			if reason == "" {
				return errors.New("alasan penolakan wajib diisi")
			}
			peResult.Status = "REJECTED"
			peResult.ValidatorID = &validatorID
			if err := u.peRepo.UpdatePEResult(ctx, tx, peResult); err != nil {
				return err
			}
			// Update report status
			if err := u.reportRepo.UpdateReportStatus(ctx, tx, peResult.ReportID, "DITOLAK_PUSKESMAS"); err != nil {
				return err
			}
			return u.auditRepo.LogAction(ctx, tx, validatorID, "REJECT_PE", "PEResult", peID, map[string]interface{}{"reason": reason})
		}

		recommendation := peResult.Recommendation
		var puskesmasName string
		if peResult.Report.Location.Kelurahan != "" {
			puskesmasName = peResult.Report.Location.Kelurahan
			if u.masterDataRepo != nil {
				kelurahan, err := u.masterDataRepo.GetKelurahanByName(ctx, peResult.Report.Location.Kelurahan)
				if err == nil && kelurahan.Puskesmas.Name != "" {
					puskesmasName = kelurahan.Puskesmas.Name
				}
			}
		}

		// Re-evaluate stock at validation time if system recommendation is fogging
		if recommendation == "FOGGING_FOKUS" || recommendation == "FOGGING_TERTUNDA" {
			item, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, "INSEKTISIDA", puskesmasName)
			if err == nil {
				if (item.Quantity - item.ReservedQuantity) >= 10.0 {
					recommendation = "FOGGING_FOKUS"
				} else {
					recommendation = "FOGGING_TERTUNDA"
				}
				peResult.Recommendation = recommendation
			}
		}

		// P0-03: DSS Advisory & Human Override
		// Allow authorized authority to override recommendation with mandatory reason & audit trail
		var overrideAudit map[string]interface{}
		if overrideDecision != "" && overrideDecision != recommendation {
			if overrideReason == "" {
				return errors.New("alasan override keputusan DSS wajib diisi")
			}
			validDecisions := map[string]bool{
				"FOGGING_FOKUS":    true,
				"FOGGING_TERTUNDA": true,
				"EDUKASI_PSN":      true,
			}
			if !validDecisions[overrideDecision] {
				return errors.New("keputusan override tidak valid. Pilihan: FOGGING_FOKUS, FOGGING_TERTUNDA, EDUKASI_PSN")
			}

			// If override is to FOGGING_FOKUS, check stock in destination warehouse
			if overrideDecision == "FOGGING_FOKUS" {
				item, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, "INSEKTISIDA", puskesmasName)
				if err != nil || (item.Quantity-item.ReservedQuantity) < 10.0 {
					// Fallback to FOGGING_TERTUNDA if stock is insufficient
					overrideDecision = "FOGGING_TERTUNDA"
				}
			}

			// Audit override details
			overrideAudit = map[string]interface{}{
				"system_recommendation": recommendation,
				"final_human_decision":  overrideDecision,
				"override_reason":       overrideReason,
				"validator_id":          validatorID,
				"timestamp":             time.Now(),
			}

			recommendation = overrideDecision
			peResult.Recommendation = overrideDecision
		}

		// If approved and validation passes, mutate status
		peResult.Status = "VALIDATED"
		peResult.ValidatorID = &validatorID

		if overrideAudit != nil {
			if err := u.auditRepo.LogAction(ctx, tx, validatorID, "OVERRIDE_DSS_RECOMMENDATION", "PEResult", peID, overrideAudit); err != nil {
				return err
			}
		}

		if err := u.peRepo.UpdatePEResult(ctx, tx, peResult); err != nil {
			return err
		}

		// b. Create Fogging Request if recommended
		if recommendation == "FOGGING_FOKUS" || recommendation == "FOGGING_TERTUNDA" {
			foggingReq := &domain.FoggingRequest{
				PEID:           peID,
				ApprovalStatus: "PENDING",
				ReservedStock:  0.0,
			}
			if err := u.foggingRepo.CreateFoggingRequest(ctx, tx, foggingReq); err != nil {
				return err
			}
		}

		// c. Create Stock Request if FOGGING_TERTUNDA
		if recommendation == "FOGGING_TERTUNDA" {
			item, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, "INSEKTISIDA", puskesmasName)
			if err == nil {
				notes := "Auto-draft: Stok tidak cukup untuk fogging pada laporan " + peResult.ReportID.String()
				if peResult.Report.ReferenceNo != "" {
					notes = "Auto-draft: Stok tidak cukup untuk fogging pada laporan " + peResult.Report.ReferenceNo
				}
				stockReq := &domain.StockRequest{
					ItemID:    item.ItemID,
					Warehouse: puskesmasName,
					Quantity:  10.0,
					Status:    "PENDING",
					Notes:     notes,
				}
				if err := tx.WithContext(ctx).Create(stockReq).Error; err != nil {
					return err
				}
			}
		}

		// d. Update Report Status
		newStatus := "CLOSED_NO_FOGGING"
		if recommendation == "FOGGING_FOKUS" || recommendation == "FOGGING_TERTUNDA" {
			newStatus = "WAITING_FOGGING_APPROVAL"
		}
		if err := u.reportRepo.UpdateReportStatus(ctx, tx, peResult.ReportID, newStatus); err != nil {
			return err
		}

		auditDetails := map[string]interface{}{
			"final_recommendation": recommendation,
		}
		return u.auditRepo.LogAction(ctx, tx, validatorID, "APPROVE_PE", "PEResult", peID, auditDetails)
	})

	if err == nil && isApproved && u.ewsKlbUsecase != nil {
		if peResult.Report.Location.Kelurahan != "" {
			_ = u.ewsKlbUsecase.CheckAndTriggerKLB(ctx, peResult.Report.Location.Kelurahan, "KELURAHAN")
		}
		if peResult.Report.Location.Kecamatan != "" {
			_ = u.ewsKlbUsecase.CheckAndTriggerKLB(ctx, peResult.Report.Location.Kecamatan, "KECAMATAN")
		}
	}

	return err
}
