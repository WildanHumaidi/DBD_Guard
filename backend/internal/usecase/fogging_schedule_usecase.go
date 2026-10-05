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

type UsedItemInput struct {
	ItemID   uuid.UUID
	Quantity float64
}

type FoggingScheduleUsecase interface {
	CreateSchedule(ctx context.Context, requestID, teamID, creatorID uuid.UUID, plannedTime time.Time) (*uuid.UUID, error)
	UpdateExecutionStatus(ctx context.Context, scheduleID, actorID uuid.UUID, status, notes string, usedItems []UsedItemInput) error
	UploadEvidence(ctx context.Context, scheduleID, uploaderID uuid.UUID, storageKey string) error
}

type foggingScheduleUsecase struct {
	scheduleRepo     repository.FoggingScheduleRepository
	foggingRepo      repository.FoggingRepository
	inventoryRepo    repository.InventoryRepository
	inventoryUsecase InventoryUsecase
	auditRepo        repository.AuditRepository
	reportRepo       repository.ReportRepository
	masterDataRepo   repository.MasterDataRepository
	notifUsecase     NotificationUsecase
}

func NewFoggingScheduleUsecase(
	scheduleRepo repository.FoggingScheduleRepository,
	foggingRepo repository.FoggingRepository,
	inventoryRepo repository.InventoryRepository,
	inventoryUsecase InventoryUsecase,
	auditRepo repository.AuditRepository,
	reportRepo repository.ReportRepository,
	masterDataRepo repository.MasterDataRepository,
	notifUsecase NotificationUsecase,
) FoggingScheduleUsecase {
	return &foggingScheduleUsecase{
		scheduleRepo:     scheduleRepo,
		foggingRepo:      foggingRepo,
		inventoryRepo:    inventoryRepo,
		inventoryUsecase: inventoryUsecase,
		auditRepo:        auditRepo,
		reportRepo:       reportRepo,
		masterDataRepo:   masterDataRepo,
		notifUsecase:     notifUsecase,
	}
}

func (u *foggingScheduleUsecase) CreateSchedule(ctx context.Context, requestID, teamID, creatorID uuid.UUID, plannedTime time.Time) (*uuid.UUID, error) {
	// 1. Validasi Fogging Request
	req, err := u.foggingRepo.GetFoggingRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}

	if req.ApprovalStatus != "APPROVED" {
		return nil, errors.New("hanya pengajuan fogging yang disetujui yang dapat dijadwalkan")
	}

	// BUG-02 & LOGIC-01: Validasi waktu (tidak boleh di masa lalu)
	if plannedTime.Before(time.Now()) {
		return nil, errors.New("waktu pelaksanaan tidak boleh di masa lalu")
	}

	// BR5: Fogging dilakukan minimal 1 hari (24 jam) setelah pengajuan disetujui
	// Asumsi `req.CreatedAt` adalah patokan waktu jika UpdatedAt tidak ada
	if plannedTime.Sub(req.CreatedAt) < 24*time.Hour {
		return nil, errors.New("jadwal fogging harus minimal 24 jam setelah persetujuan (Business Rule BR5)")
	}

	// LOGIC-02: Cek apakah sudah ada jadwal aktif untuk request ini
	activeSchedule, err := u.scheduleRepo.GetActiveScheduleByRequestID(ctx, requestID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if activeSchedule != nil {
		return nil, errors.New("sudah ada jadwal fogging yang aktif (belum selesai/gagal) untuk pengajuan ini")
	}

	// PRD Gap: Kapasitas Tim & Anti-Overbooking (1 tugas per hari)
	count, err := u.scheduleRepo.CountSchedulesByTeamAndDate(ctx, teamID, plannedTime)
	if err != nil {
		return nil, err
	}
	if count >= 1 {
		return nil, errors.New("tim fogging sudah memiliki jadwal pada hari tersebut (maksimal 1 tugas per hari)")
	}

	// 2. Buat Jadwal
	schedule := &domain.FoggingSchedule{
		RequestID:       requestID,
		TeamID:          teamID,
		PlannedTime:     plannedTime,
		ExecutionStatus: "SCHEDULED",
	}

	err = u.scheduleRepo.WithTransaction(func(tx *gorm.DB) error {
		// a. Insert Schedule
		if err := u.scheduleRepo.CreateScheduleWithTx(ctx, tx, schedule); err != nil {
			return err
		}

		// Update Report Status
		if err := u.reportRepo.UpdateReportStatus(ctx, tx, req.PE.ReportID, "FOGGING_SCHEDULED"); err != nil {
			return err
		}

		// b. Audit Trail
		auditDetails := map[string]interface{}{
			"planned_time": plannedTime,
			"team_id":      teamID,
		}
		if err := u.auditRepo.LogAction(ctx, tx, creatorID, "CREATE_FOGGING_SCHEDULE", "FoggingSchedule", schedule.ID, auditDetails); err != nil {
			return err
		}

		return nil
	})

	if err == nil && u.notifUsecase != nil {
		report, errReport := u.reportRepo.GetReportByID(ctx, req.PE.ReportID)
		if errReport == nil && !report.Reporter.AnonymousFlag {
			_ = u.notifUsecase.NotifyFoggingScheduled(ctx, req.PE.ReportID.String(), plannedTime.Format("02 Jan 2006 15:04"))
		}
	}

	return &schedule.ID, err
}

func (u *foggingScheduleUsecase) UpdateExecutionStatus(ctx context.Context, scheduleID, actorID uuid.UUID, status, notes string, usedItems []UsedItemInput) error {
	schedule, err := u.scheduleRepo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return err
	}

	if schedule.ExecutionStatus == "COMPLETED" || schedule.ExecutionStatus == "FAILED" {
		return errors.New("jadwal fogging yang sudah selesai atau gagal tidak dapat diupdate lagi")
	}

	if status == "COMPLETED" && len(usedItems) == 0 {
		return errors.New("penggunaan inventory wajib diisi saat menyelesaikan jadwal fogging")
	}

	return u.scheduleRepo.WithTransaction(func(tx *gorm.DB) error {
		// 1. Update Status
		if err := u.scheduleRepo.UpdateExecutionStatus(ctx, tx, scheduleID, status, notes); err != nil {
			return err
		}

		// Update CitizenReport status
		reportStatus := "FOGGING_IN_PROGRESS"
		if status == "COMPLETED" {
			reportStatus = "COMPLETED"
		} else if status == "FAILED" {
			reportStatus = "FOGGING_FAILED"
		} else if status == "ON_THE_WAY" || status == "IN_PROGRESS" {
			reportStatus = "FOGGING_IN_PROGRESS"
		}
		if err := u.reportRepo.UpdateReportStatus(ctx, tx, schedule.Request.PE.ReportID, reportStatus); err != nil {
			return err
		}

		// 2. Handle Reserved Stock and Used Items
		if status == "COMPLETED" {
			// Find the reserved item (Insektisida at Kelurahan) to release the reservation
			pe, err := u.reportRepo.GetReportByID(ctx, schedule.Request.PE.ReportID)
			if err == nil && schedule.Request.ReservedStock > 0 {
				puskesmasName := pe.Location.Kelurahan
				kelurahan, err := u.masterDataRepo.GetKelurahanByName(ctx, pe.Location.Kelurahan)
				if err == nil && kelurahan.Puskesmas.Name != "" {
					puskesmasName = kelurahan.Puskesmas.Name
				}
				resItem, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, "INSEKTISIDA", puskesmasName)
				if err == nil {
					// Release the reservation, but do not deduct quantity yet (deduction happens below per usedItems)
					if err := u.inventoryRepo.ReleaseReservedStock(ctx, tx, resItem.ItemID, schedule.Request.ReservedStock); err != nil {
						return err
					}
				}
			}

			for _, item := range usedItems {
				// Pemotongan stok aktual via Usecase agar mendukung cross-allocation
				if err := u.inventoryUsecase.UpdateStockTx(ctx, tx, item.ItemID, actorID, -item.Quantity, "Pemakaian material fogging"); err != nil {
					return err
				}

				// Catat sebagai Audited Usage Record (status RECORDED) dalam transaksi yang sama
				usage := &domain.InventoryUsageLog{
					ItemID:     item.ItemID,
					TeamID:     schedule.TeamID,
					ScheduleID: schedule.ID,
					Quantity:   item.Quantity,
					Status:     "RECORDED",
					Notes:      notes,
				}
				if err := u.inventoryRepo.CreateInventoryUsageLog(ctx, tx, usage); err != nil {
					return err
				}
			}
		} else if status == "FAILED" {
			// Release reserved stock without deducting actual stock
			if schedule.Request.ReservedStock > 0 {
				pe, err := u.reportRepo.GetReportByID(ctx, schedule.Request.PE.ReportID)
				if err == nil {
					puskesmasName := pe.Location.Kelurahan
					kelurahan, err := u.masterDataRepo.GetKelurahanByName(ctx, pe.Location.Kelurahan)
					if err == nil && kelurahan.Puskesmas.Name != "" {
						puskesmasName = kelurahan.Puskesmas.Name
					}
					resItem, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, "INSEKTISIDA", puskesmasName)
					if err == nil {
						if err := u.inventoryRepo.ReleaseReservedStock(ctx, tx, resItem.ItemID, schedule.Request.ReservedStock); err != nil {
							return err
						}
					}
				}
			}
		}

		// 3. Audit Trail untuk update status Fogging Schedule
		auditDetails := map[string]interface{}{
			"execution_status": status,
			"notes":            notes,
			"used_items":       usedItems,
		}
		if err := u.auditRepo.LogAction(ctx, tx, actorID, "UPDATE_FOGGING_EXECUTION_STATUS", "FoggingSchedule", scheduleID, auditDetails); err != nil {
			return err
		}

		return nil
	})
}

func (u *foggingScheduleUsecase) UploadEvidence(ctx context.Context, scheduleID, uploaderID uuid.UUID, storageKey string) error {
	evidence := &domain.FoggingEvidence{
		ScheduleID: scheduleID,
		StorageKey: storageKey,
		UploadedBy: uploaderID,
	}

	err := u.scheduleRepo.SaveEvidence(ctx, evidence)
	if err != nil {
		return err
	}

	auditDetails := map[string]interface{}{
		"evidence_id": evidence.ID,
		"storage_key": storageKey,
	}
	_ = u.auditRepo.LogAction(ctx, nil, uploaderID, "UPLOAD_FOGGING_EVIDENCE", "FoggingEvidence", evidence.ID, auditDetails)

	return nil
}
