package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockScheduleRepo struct {
	repository.FoggingScheduleRepository
	schedule               *domain.FoggingSchedule
	updatedExecutionStatus string
	updatedNotes           string
}

func (m *mockScheduleRepo) GetScheduleByID(ctx context.Context, id uuid.UUID) (*domain.FoggingSchedule, error) {
	return m.schedule, nil
}

func (m *mockScheduleRepo) UpdateExecutionStatus(ctx context.Context, tx *gorm.DB, scheduleID uuid.UUID, status, notes string) error {
	m.updatedExecutionStatus = status
	m.updatedNotes = notes
	return nil
}

func (m *mockScheduleRepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

type mockInvRepoForExecution struct {
	repository.InventoryRepository
	item               *domain.InventoryItem
	releasedItemID     uuid.UUID
	releasedQty        float64
	createdUsageLog    *domain.InventoryUsageLog
	updatedUsageStatus string
	usageLog           *domain.InventoryUsageLog
	stockUpdatedQty    float64
}

func (m *mockInvRepoForExecution) GetItemByTypeAndWarehouse(ctx context.Context, itemType, warehouse string) (*domain.InventoryItem, error) {
	return m.item, nil
}

func (m *mockInvRepoForExecution) ReleaseReservedStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error {
	m.releasedItemID = itemID
	m.releasedQty = quantity
	return nil
}

func (m *mockInvRepoForExecution) CreateInventoryUsageLog(ctx context.Context, tx *gorm.DB, log *domain.InventoryUsageLog) error {
	m.createdUsageLog = log
	return nil
}

func (m *mockInvRepoForExecution) GetInventoryUsageLogByID(ctx context.Context, id uuid.UUID) (*domain.InventoryUsageLog, error) {
	if m.usageLog != nil {
		return m.usageLog, nil
	}
	return nil, errors.New("not found")
}

func (m *mockInvRepoForExecution) UpdateInventoryUsageLogStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	m.updatedUsageStatus = status
	return nil
}

func (m *mockInvRepoForExecution) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (m *mockInvRepoForExecution) UpdateStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantityDelta float64) error {
	m.stockUpdatedQty += quantityDelta
	return nil
}

type mockInvUsecase struct {
	usecase.InventoryUsecase
	deductedItems []struct {
		ItemID uuid.UUID
		Qty    float64
	}
	failUpdateStock bool
}

func (m *mockInvUsecase) UpdateStockTx(ctx context.Context, tx *gorm.DB, itemID, actorID uuid.UUID, quantityDelta float64, reason string) error {
	if m.failUpdateStock {
		return errors.New("insufficient stock")
	}
	m.deductedItems = append(m.deductedItems, struct {
		ItemID uuid.UUID
		Qty    float64
	}{ItemID: itemID, Qty: quantityDelta})
	return nil
}

type mockNotifUsecase struct {
	usecase.NotificationUsecase
}

func TestUpdateExecutionStatus_SingleSourceOfTruth_Success(t *testing.T) {
	scheduleID := uuid.New()
	actorID := uuid.New()
	reportID := uuid.New()
	teamID := actorID
	itemID := uuid.New()

	schedule := &domain.FoggingSchedule{
		ID:              scheduleID,
		ExecutionStatus: "SCHEDULED",
		TeamID:          teamID,
		Request: domain.FoggingRequest{
			ReservedStock: 10.0,
			PE: domain.PEResult{
				ReportID: reportID,
			},
		},
	}

	schedRepo := &mockScheduleRepo{schedule: schedule}
	invRepo := &mockInvRepoForExecution{
		item: &domain.InventoryItem{ItemID: itemID, Warehouse: "Puskesmas Test"},
	}
	invUsecase := &mockInvUsecase{}
	repRepo := &mockRepRepo{}
	mdRepo := &mockMDRepo{}
	audRepo := &mockAudRepo{}
	notif := &mockNotifUsecase{}

	u := usecase.NewFoggingScheduleUsecase(schedRepo, nil, invRepo, invUsecase, audRepo, repRepo, mdRepo, notif)

	usedItems := []usecase.UsedItemInput{
		{ItemID: itemID, Quantity: 8.5},
	}

	err := u.UpdateExecutionStatus(context.Background(), scheduleID, actorID, "COMPLETED", "Fogging sukses selesai", usedItems)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 1. Verify execution status was updated to COMPLETED
	if schedRepo.updatedExecutionStatus != "COMPLETED" {
		t.Errorf("expected schedule status COMPLETED, got %s", schedRepo.updatedExecutionStatus)
	}

	// 2. Verify reserved stock was released (10.0)
	if invRepo.releasedQty != 10.0 {
		t.Errorf("expected released reserved stock 10.0, got %f", invRepo.releasedQty)
	}

	// 3. Verify actual deduction occurred via UpdateStockTx
	if len(invUsecase.deductedItems) != 1 || invUsecase.deductedItems[0].Qty != -8.5 {
		t.Errorf("expected stock deduction of -8.5, got %v", invUsecase.deductedItems)
	}

	// 4. Verify audited InventoryUsageLog was created with RECORDED status
	if invRepo.createdUsageLog == nil {
		t.Fatal("expected InventoryUsageLog to be created")
	}
	if invRepo.createdUsageLog.Status != "RECORDED" {
		t.Errorf("expected usage log status RECORDED, got %s", invRepo.createdUsageLog.Status)
	}
	if invRepo.createdUsageLog.Quantity != 8.5 {
		t.Errorf("expected usage log quantity 8.5, got %f", invRepo.createdUsageLog.Quantity)
	}
}

func TestUpdateExecutionStatus_StockDeductionFailure_Rollback(t *testing.T) {
	scheduleID := uuid.New()
	actorID := uuid.New()
	reportID := uuid.New()
	teamID := actorID
	itemID := uuid.New()

	schedule := &domain.FoggingSchedule{
		ID:              scheduleID,
		ExecutionStatus: "SCHEDULED",
		TeamID:          teamID,
		Request: domain.FoggingRequest{
			ReservedStock: 10.0,
			PE: domain.PEResult{
				ReportID: reportID,
			},
		},
	}

	schedRepo := &mockScheduleRepo{schedule: schedule}
	invRepo := &mockInvRepoForExecution{
		item: &domain.InventoryItem{ItemID: itemID, Warehouse: "Puskesmas Test"},
	}
	invUsecase := &mockInvUsecase{failUpdateStock: true}
	repRepo := &mockRepRepo{}
	mdRepo := &mockMDRepo{}
	audRepo := &mockAudRepo{}
	notif := &mockNotifUsecase{}

	u := usecase.NewFoggingScheduleUsecase(schedRepo, nil, invRepo, invUsecase, audRepo, repRepo, mdRepo, notif)

	usedItems := []usecase.UsedItemInput{
		{ItemID: itemID, Quantity: 15.0},
	}

	err := u.UpdateExecutionStatus(context.Background(), scheduleID, actorID, "COMPLETED", "Stock short", usedItems)
	if err == nil {
		t.Fatal("expected error due to insufficient stock, got nil")
	}
	if !errors.Is(err, errors.New("insufficient stock")) && err.Error() != "insufficient stock" {
		t.Errorf("expected insufficient stock error, got %v", err)
	}
}

func TestApproveFoggingUsage_DoesNotDeductStock(t *testing.T) {
	approverID := uuid.New()
	usageID := uuid.New()
	itemID := uuid.New()

	invRepo := &mockInvRepoForExecution{
		usageLog: &domain.InventoryUsageLog{
			ID:       usageID,
			ItemID:   itemID,
			Quantity: 10.0,
			Status:   "PENDING_APPROVAL",
		},
	}
	audRepo := &mockAudRepo{}

	u := usecase.NewInventoryUsecase(invRepo, audRepo)

	err := u.ApproveFoggingUsage(context.Background(), approverID, usageID, true, "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify status is APPROVED
	if invRepo.updatedUsageStatus != "APPROVED" {
		t.Errorf("expected usage status APPROVED, got %s", invRepo.updatedUsageStatus)
	}

	// Verify stock was NOT deducted during approval (stockUpdatedQty must be 0)
	if invRepo.stockUpdatedQty != 0 {
		t.Errorf("expected no stock deduction during approval, but got delta %f", invRepo.stockUpdatedQty)
	}
}

// Suppress unused imports
var _ = time.Now
