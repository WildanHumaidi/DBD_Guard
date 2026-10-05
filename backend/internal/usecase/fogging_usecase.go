package usecase

import (
	"context"
	"errors"

	"github.com/dbd-guard/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FoggingUsecase interface {
	ApproveFoggingRequest(ctx context.Context, reqID, approverID uuid.UUID, decision, reason string) error
}

type foggingUsecase struct {
	foggingRepo   repository.FoggingRepository
	auditRepo     repository.AuditRepository
	reportRepo    repository.ReportRepository
	inventoryRepo repository.InventoryRepository
	masterDataRepo repository.MasterDataRepository
}

func NewFoggingUsecase(fr repository.FoggingRepository, ar repository.AuditRepository, rr repository.ReportRepository, ir repository.InventoryRepository, md repository.MasterDataRepository) FoggingUsecase {
	return &foggingUsecase{
		foggingRepo:    fr,
		auditRepo:      ar,
		reportRepo:     rr,
		inventoryRepo:  ir,
		masterDataRepo: md,
	}
}

func (u *foggingUsecase) ApproveFoggingRequest(ctx context.Context, reqID, approverID uuid.UUID, decision, reason string) error {
	// Rule 4: If rejected, reason is mandatory
	if decision == "REJECTED" && reason == "" {
		return errors.New("alasan penolakan wajib diisi")
	}

	// 1. Get the Request
	foggingReq, err := u.foggingRepo.GetFoggingRequestByID(ctx, reqID)
	if err != nil {
		return err
	}

	if foggingReq.ApprovalStatus != "PENDING" {
		return errors.New("hanya pengajuan dengan status PENDING yang dapat di-approve/reject")
	}

	if foggingReq.PEID == uuid.Nil {
		return errors.New("data PE tidak valid untuk pengajuan fogging ini")
	}

	// 2. Perform Transaction
	return u.foggingRepo.WithTransaction(func(tx *gorm.DB) error {
		// a. Update Status
		if err := u.foggingRepo.UpdateFoggingRequestStatus(ctx, tx, reqID, decision, reason); err != nil {
			return err
		}

		// Update CitizenReport status
		newReportStatus := "FOGGING_REJECTED"
		if decision == "APPROVED" {
			newReportStatus = "FOGGING_APPROVED"
			
			// Validate and reserve stock if not yet reserved
			if foggingReq.ReservedStock == 0 {
				pe, err := u.reportRepo.GetReportByID(ctx, foggingReq.PE.ReportID)
				if err != nil {
					return err
				}
				warehouseName := "Puskesmas " + pe.Location.Kecamatan // Fallback
				kelurahan, errMD := u.masterDataRepo.GetKelurahanByName(ctx, pe.Location.Kelurahan)
				if errMD == nil && kelurahan.Puskesmas.Name != "" {
					warehouseName = kelurahan.Puskesmas.Name
				}
				
				item, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, "INSEKTISIDA", warehouseName)
				if err != nil {
					return errors.New("gagal mendapatkan data stok untuk gudang " + warehouseName)
				}
				
				// Atomic lock for stock checking
				lockedItem, err := u.inventoryRepo.GetItemByIDForUpdate(ctx, tx, item.ItemID)
				if err != nil {
					return err
				}
				
				if (lockedItem.Quantity - lockedItem.ReservedQuantity) < 10 {
					return errors.New("stok tersedia (di luar reservasi) tidak mencukupi untuk menyetujui fogging. Harap tunggu stok datang dari Dinkes")
				}
				if err := u.inventoryRepo.ReserveStock(ctx, tx, lockedItem.ItemID, 10); err != nil {
					return err
				}
				if err := u.foggingRepo.UpdateFoggingRequestReservedStock(ctx, tx, reqID, 10); err != nil {
					return err
				}
			}
		}
		if err := u.reportRepo.UpdateReportStatus(ctx, tx, foggingReq.PE.ReportID, newReportStatus); err != nil {
			return err
		}

		// Release reserved stock if rejected
		if decision == "REJECTED" && foggingReq.ReservedStock > 0 {
			pe, err := u.reportRepo.GetReportByID(ctx, foggingReq.PE.ReportID)
			if err == nil {
				warehouseName := "Puskesmas " + pe.Location.Kecamatan // Fallback
				kelurahan, errMD := u.masterDataRepo.GetKelurahanByName(ctx, pe.Location.Kelurahan)
				if errMD == nil && kelurahan.Puskesmas.Name != "" {
					warehouseName = kelurahan.Puskesmas.Name
				}
				
				item, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, "INSEKTISIDA", warehouseName)
				if err == nil {
					if err := u.inventoryRepo.ReleaseReservedStock(ctx, tx, item.ItemID, foggingReq.ReservedStock); err != nil {
						return err
					}
				}
			}
		}

		// b. Log to Audit (Rule 6)
		auditDetails := map[string]interface{}{
			"decision": decision,
			"reason":   reason,
		}
		if err := u.auditRepo.LogAction(ctx, tx, approverID, "APPROVE_FOGGING", "FoggingRequest", reqID, auditDetails); err != nil {
			return err
		}

		return nil
	})
}
