package usecase_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockStockTransferRepo struct {
	repository.InventoryRepository
	mu           sync.Mutex
	items        map[uuid.UUID]*domain.InventoryItem
	transfers    map[uuid.UUID]*domain.StockTransfer
	requests     map[uuid.UUID]*domain.StockRequest
}

func newMockStockTransferRepo() *mockStockTransferRepo {
	return &mockStockTransferRepo{
		items:     make(map[uuid.UUID]*domain.InventoryItem),
		transfers: make(map[uuid.UUID]*domain.StockTransfer),
		requests:  make(map[uuid.UUID]*domain.StockRequest),
	}
}

func (m *mockStockTransferRepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (m *mockStockTransferRepo) GetItemByID(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[id]
	if !ok {
		return nil, errors.New("item not found")
	}
	return item, nil
}

func (m *mockStockTransferRepo) GetItemByTypeAndWarehouse(ctx context.Context, itemType, warehouse string) (*domain.InventoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, item := range m.items {
		if item.ItemType == itemType && item.Warehouse == warehouse {
			return item, nil
		}
	}
	return nil, errors.New("item not found for warehouse")
}

func (m *mockStockTransferRepo) CreateItem(ctx context.Context, item *domain.InventoryItem) error {
	return m.CreateItemWithTx(ctx, nil, item)
}

func (m *mockStockTransferRepo) CreateItemWithTx(ctx context.Context, tx *gorm.DB, item *domain.InventoryItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if item.ItemID == uuid.Nil {
		item.ItemID = uuid.New()
	}
	m.items[item.ItemID] = item
	return nil
}

func (m *mockStockTransferRepo) GetItemByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.InventoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[id]
	if !ok {
		return nil, errors.New("item not found")
	}
	return item, nil
}

func (m *mockStockTransferRepo) UpdateStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantityDelta float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[itemID]
	if !ok {
		return errors.New("item not found")
	}
	item.Quantity += quantityDelta
	return nil
}

func (m *mockStockTransferRepo) CreateStockTransfer(ctx context.Context, tx *gorm.DB, transfer *domain.StockTransfer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if transfer.ID == uuid.Nil {
		transfer.ID = uuid.New()
	}
	m.transfers[transfer.ID] = transfer
	return nil
}

func (m *mockStockTransferRepo) GetStockTransferByID(ctx context.Context, id uuid.UUID) (*domain.StockTransfer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tr, ok := m.transfers[id]
	if !ok {
		return nil, errors.New("transfer not found")
	}
	return tr, nil
}

func (m *mockStockTransferRepo) UpdateStockTransferStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	tr, ok := m.transfers[id]
	if !ok {
		return errors.New("transfer not found")
	}
	if tr.Status != "IN_TRANSIT" {
		return errors.New("transfer tidak dapat diproses: status bukan IN_TRANSIT")
	}
	tr.Status = status
	return nil
}

func (m *mockStockTransferRepo) CreateStockRequest(ctx context.Context, req *domain.StockRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if req.ID == uuid.Nil {
		req.ID = uuid.New()
	}
	m.requests[req.ID] = req
	return nil
}

func (m *mockStockTransferRepo) GetStockRequestByID(ctx context.Context, id uuid.UUID) (*domain.StockRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	req, ok := m.requests[id]
	if !ok {
		return nil, errors.New("request not found")
	}
	return req, nil
}

func (m *mockStockTransferRepo) UpdateStockRequestStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	req, ok := m.requests[id]
	if !ok {
		return errors.New("request not found")
	}
	req.Status = status
	return nil
}

func (m *mockStockTransferRepo) FindDonorWarehouse(ctx context.Context, tx *gorm.DB, itemType string, requiredAmount float64, excludeWarehouse string) (*domain.InventoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, item := range m.items {
		if item.ItemType == itemType && item.Warehouse != excludeWarehouse && (item.Quantity-item.ReservedQuantity-item.MinimumStock) >= requiredAmount {
			return item, nil
		}
	}
	return nil, errors.New("donor warehouse not found")
}

func TestStockTransfer_AccountingConservation(t *testing.T) {
	repo := newMockStockTransferRepo()
	audit := &mockAudRepo{}
	u := usecase.NewInventoryUsecase(repo, audit)

	ctx := context.Background()
	actorID := uuid.New()
	dinkesItemID := uuid.New()
	puskesmasItemID := uuid.New()

	// Initial State: Dinkes has 100L, Puskesmas has 20L
	initialDinkesStock := 100.0
	initialPuskStock := 20.0
	transferQty := 30.0

	repo.items[dinkesItemID] = &domain.InventoryItem{
		ItemID:           dinkesItemID,
		ItemType:         "INSEKTISIDA",
		Warehouse:        "Gudang Dinkes",
		Quantity:         initialDinkesStock,
		ReservedQuantity: 10.0, // Available = 90.0
	}
	repo.items[puskesmasItemID] = &domain.InventoryItem{
		ItemID:           puskesmasItemID,
		ItemType:         "INSEKTISIDA",
		Warehouse:        "Puskesmas Sukmajaya",
		Quantity:         initialPuskStock,
		ReservedQuantity: 0.0,
	}

	// 1. Cannot transfer more than available unreserved stock (e.g. 95L when available is 90L)
	errExcess := u.CreateStockTransfer(ctx, actorID, dinkesItemID, "Gudang Dinkes", "Puskesmas Sukmajaya", 95.0, "Excess")
	if errExcess == nil || !strings.Contains(errExcess.Error(), "stok pengirim tidak cukup") {
		t.Fatalf("expected error for exceeding unreserved stock, got: %v", errExcess)
	}

	// 2. Transfer 30L from Dinkes to Puskesmas
	err := u.CreateStockTransfer(ctx, actorID, dinkesItemID, "Gudang Dinkes", "Puskesmas Sukmajaya", transferQty, "Routine replenishment")
	if err != nil {
		t.Fatalf("CreateStockTransfer failed: %v", err)
	}

	// Find created transfer
	var transferID uuid.UUID
	for id := range repo.transfers {
		transferID = id
		break
	}

	// In-Transit Accounting Check:
	// Dinkes = 70, Puskesmas = 20, In-Transit = 30 -> Total = 120 (Conserved!)
	dinkesCurrent := repo.items[dinkesItemID].Quantity
	puskCurrent := repo.items[puskesmasItemID].Quantity
	inTransit := repo.transfers[transferID].Quantity
	if dinkesCurrent+puskCurrent+inTransit != initialDinkesStock+initialPuskStock {
		t.Fatalf("Accounting broken in transit: Dinkes=%f, Pusk=%f, InTransit=%f, Expected Sum=%f",
			dinkesCurrent, puskCurrent, inTransit, initialDinkesStock+initialPuskStock)
	}

	// 3. Receive transfer (Accept)
	// We need destItem in the repo to be updated by ReceiveStockTransfer
	// In the real DB it finds by warehouse; in mock, let's update puskesmas item
	err = u.ReceiveStockTransfer(ctx, actorID, transferID, true, "Received safely")
	if err != nil {
		t.Fatalf("ReceiveStockTransfer failed: %v", err)
	}

	// Verify status is COMPLETED
	if repo.transfers[transferID].Status != "COMPLETED" {
		t.Fatalf("expected transfer status COMPLETED, got: %s", repo.transfers[transferID].Status)
	}

	// 4. Verify duplicate receive call fails
	errDup := u.ReceiveStockTransfer(ctx, actorID, transferID, true, "Duplicate receive")
	if errDup == nil {
		t.Fatalf("expected duplicate receive to fail, but succeeded")
	}
}

func TestStockTransfer_RejectionRefund(t *testing.T) {
	repo := newMockStockTransferRepo()
	audit := &mockAudRepo{}
	u := usecase.NewInventoryUsecase(repo, audit)

	ctx := context.Background()
	actorID := uuid.New()
	dinkesItemID := uuid.New()

	initialDinkesStock := 50.0
	transferQty := 20.0

	repo.items[dinkesItemID] = &domain.InventoryItem{
		ItemID:           dinkesItemID,
		ItemType:         "INSEKTISIDA",
		Warehouse:        "Gudang Dinkes",
		Quantity:         initialDinkesStock,
		ReservedQuantity: 0.0,
	}

	err := u.CreateStockTransfer(ctx, actorID, dinkesItemID, "Gudang Dinkes", "Puskesmas Beji", transferQty, "Damaged in transit test")
	if err != nil {
		t.Fatalf("CreateStockTransfer failed: %v", err)
	}

	var transferID uuid.UUID
	for id := range repo.transfers {
		transferID = id
		break
	}

	// Reject transfer (e.g. damaged goods during delivery, refund source)
	err = u.ReceiveStockTransfer(ctx, actorID, transferID, false, "Damaged during shipment")
	if err != nil {
		t.Fatalf("ReceiveStockTransfer reject failed: %v", err)
	}

	// Stock Dinkes should be refunded back to 50.0
	if repo.items[dinkesItemID].Quantity != initialDinkesStock {
		t.Fatalf("expected refunded stock %f, got: %f", initialDinkesStock, repo.items[dinkesItemID].Quantity)
	}
	if repo.transfers[transferID].Status != "REJECTED" {
		t.Fatalf("expected transfer status REJECTED, got: %s", repo.transfers[transferID].Status)
	}
}

func TestApproveStockRequest_UsesStockTransferWorkflow_NoVirtualStock(t *testing.T) {
	repo := newMockStockTransferRepo()
	audit := &mockAudRepo{}
	u := usecase.NewInventoryUsecase(repo, audit)

	ctx := context.Background()
	actorDinkesID := uuid.New()
	puskOfficerID := uuid.New()
	dinkesItemID := uuid.New()
	puskItemID := uuid.New()

	initialDinkesStock := 100.0
	initialPuskStock := 10.0
	requestQty := 25.0

	// 1. Setup initial inventory
	repo.items[dinkesItemID] = &domain.InventoryItem{
		ItemID:           dinkesItemID,
		ItemType:         "INSEKTISIDA",
		Warehouse:        "Gudang Utama Dinkes",
		Quantity:         initialDinkesStock,
		ReservedQuantity: 0.0,
		Unit:             "Liter",
	}

	repo.items[puskItemID] = &domain.InventoryItem{
		ItemID:           puskItemID,
		ItemType:         "INSEKTISIDA",
		Warehouse:        "Puskesmas Beji",
		Quantity:         initialPuskStock,
		ReservedQuantity: 0.0,
		Unit:             "Liter",
	}

	// 2. Puskesmas Beji creates StockRequest
	reqID := uuid.New()
	repo.requests[reqID] = &domain.StockRequest{
		ID:        reqID,
		ItemID:    puskItemID,
		Warehouse: "Puskesmas Beji",
		Quantity:  requestQty,
		Status:    "PENDING",
		Notes:     "Permintaan stok fogging darurat",
	}

	// 3. Dinkes Approves StockRequest
	err := u.ApproveStockRequest(ctx, actorDinkesID, reqID, "APPROVED")
	if err != nil {
		t.Fatalf("ApproveStockRequest failed: %v", err)
	}

	// Invariant 1: StockRequest status is APPROVED
	if repo.requests[reqID].Status != "APPROVED" {
		t.Fatalf("expected request status APPROVED, got: %s", repo.requests[reqID].Status)
	}

	// Invariant 2: A StockTransfer is created with IN_TRANSIT
	if len(repo.transfers) != 1 {
		t.Fatalf("expected 1 stock transfer created, got: %d", len(repo.transfers))
	}

	var transfer *domain.StockTransfer
	for _, tr := range repo.transfers {
		transfer = tr
		break
	}

	if transfer.Status != "IN_TRANSIT" {
		t.Fatalf("expected transfer status IN_TRANSIT, got: %s", transfer.Status)
	}
	if transfer.Quantity != requestQty {
		t.Fatalf("expected transfer quantity %f, got: %f", requestQty, transfer.Quantity)
	}
	if transfer.FromWarehouse != "Gudang Utama Dinkes" {
		t.Fatalf("expected from warehouse Gudang Utama Dinkes, got: %s", transfer.FromWarehouse)
	}
	if transfer.ToWarehouse != "Puskesmas Beji" {
		t.Fatalf("expected to warehouse Puskesmas Beji, got: %s", transfer.ToWarehouse)
	}

	// Invariant 3: Source stock is immediately deducted
	dinkesCurrent := repo.items[dinkesItemID].Quantity
	if dinkesCurrent != initialDinkesStock-requestQty {
		t.Fatalf("expected Dinkes stock %f, got: %f", initialDinkesStock-requestQty, dinkesCurrent)
	}

	// Invariant 4: Puskesmas stock is NOT yet added (goods are in transit!)
	puskCurrent := repo.items[puskItemID].Quantity
	if puskCurrent != initialPuskStock {
		t.Fatalf("BUG-09 Residual! Virtual stock created immediately! Expected Pusk stock %f, got: %f", initialPuskStock, puskCurrent)
	}

	// Invariant 5: Total stock conservation during IN_TRANSIT
	if dinkesCurrent+puskCurrent+transfer.Quantity != initialDinkesStock+initialPuskStock {
		t.Fatalf("Accounting invariant broken! Sum=%f, Expected=%f", dinkesCurrent+puskCurrent+transfer.Quantity, initialDinkesStock+initialPuskStock)
	}

	// 4. Puskesmas Receives Transfer
	err = u.ReceiveStockTransfer(ctx, puskOfficerID, transfer.ID, true, "Goods accepted at Puskesmas Beji")
	if err != nil {
		t.Fatalf("ReceiveStockTransfer failed: %v", err)
	}

	// Invariant 6: Destination stock is now updated and transfer is COMPLETED
	if transfer.Status != "COMPLETED" {
		t.Fatalf("expected transfer status COMPLETED, got: %s", transfer.Status)
	}
	if repo.items[puskItemID].Quantity != initialPuskStock+requestQty {
		t.Fatalf("expected Pusk stock %f, got: %f", initialPuskStock+requestQty, repo.items[puskItemID].Quantity)
	}

	// Invariant 7: Total stock conserved post-completion
	finalTotal := repo.items[dinkesItemID].Quantity + repo.items[puskItemID].Quantity
	if finalTotal != initialDinkesStock+initialPuskStock {
		t.Fatalf("Total stock not conserved! Final Total=%f, Expected=%f", finalTotal, initialDinkesStock+initialPuskStock)
	}
}
