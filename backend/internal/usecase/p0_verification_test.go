package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// =========================================================================
// MOCKS FOR COMPREHENSIVE P0 VERIFICATION
// =========================================================================

type p0MockInventoryRepo struct {
	repository.InventoryRepository
	mu        sync.Mutex
	items     map[uuid.UUID]*domain.InventoryItem
	transfers map[uuid.UUID]*domain.StockTransfer
	requests  map[uuid.UUID]*domain.StockRequest
}

func newP0MockInventoryRepo() *p0MockInventoryRepo {
	return &p0MockInventoryRepo{
		items:     make(map[uuid.UUID]*domain.InventoryItem),
		transfers: make(map[uuid.UUID]*domain.StockTransfer),
		requests:  make(map[uuid.UUID]*domain.StockRequest),
	}
}

func (m *p0MockInventoryRepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (m *p0MockInventoryRepo) CreateItem(ctx context.Context, item *domain.InventoryItem) error {
	return m.CreateItemWithTx(ctx, nil, item)
}

func (m *p0MockInventoryRepo) CreateItemWithTx(ctx context.Context, tx *gorm.DB, item *domain.InventoryItem) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if item.ItemID == uuid.Nil {
		item.ItemID = uuid.New()
	}
	m.items[item.ItemID] = item
	return nil
}

func (m *p0MockInventoryRepo) GetItemByID(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.items[id]
	if !ok {
		return nil, errors.New("item not found")
	}
	return item, nil
}

func (m *p0MockInventoryRepo) GetItemByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.InventoryItem, error) {
	return m.GetItemByID(ctx, id)
}

func (m *p0MockInventoryRepo) GetItemByTypeAndWarehouse(ctx context.Context, itemType, warehouse string) (*domain.InventoryItem, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, it := range m.items {
		if it.ItemType == itemType && it.Warehouse == warehouse {
			return it, nil
		}
	}
	return nil, errors.New("item not found for warehouse")
}

func (m *p0MockInventoryRepo) UpdateStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantityDelta float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	it, ok := m.items[itemID]
	if !ok {
		return errors.New("item not found")
	}
	it.Quantity += quantityDelta
	return nil
}

func (m *p0MockInventoryRepo) CreateStockRequest(ctx context.Context, req *domain.StockRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if req.ID == uuid.Nil {
		req.ID = uuid.New()
	}
	m.requests[req.ID] = req
	return nil
}

func (m *p0MockInventoryRepo) GetStockRequestByID(ctx context.Context, id uuid.UUID) (*domain.StockRequest, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	req, ok := m.requests[id]
	if !ok {
		return nil, errors.New("stock request not found")
	}
	return req, nil
}

func (m *p0MockInventoryRepo) UpdateStockRequestStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	req, ok := m.requests[id]
	if !ok {
		return errors.New("stock request not found")
	}
	if req.Status != "PENDING" {
		return errors.New("stock request tidak dapat diproses: status bukan PENDING atau sudah diproses sebelumnya")
	}
	req.Status = status
	return nil
}

func (m *p0MockInventoryRepo) CreateStockTransfer(ctx context.Context, tx *gorm.DB, transfer *domain.StockTransfer) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if transfer.ID == uuid.Nil {
		transfer.ID = uuid.New()
	}
	m.transfers[transfer.ID] = transfer
	return nil
}

func (m *p0MockInventoryRepo) GetStockTransferByID(ctx context.Context, id uuid.UUID) (*domain.StockTransfer, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	tr, ok := m.transfers[id]
	if !ok {
		return nil, errors.New("stock transfer not found")
	}
	return tr, nil
}

func (m *p0MockInventoryRepo) UpdateStockTransferStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	tr, ok := m.transfers[id]
	if !ok {
		return errors.New("stock transfer not found")
	}
	if tr.Status != "IN_TRANSIT" {
		return errors.New("transfer tidak dapat diproses: status bukan IN_TRANSIT")
	}
	tr.Status = status
	return nil
}

type p0MockAuditRepo struct {
	repository.AuditRepository
	mu   sync.Mutex
	logs []domain.AuditLog
}

func (m *p0MockAuditRepo) LogAction(ctx context.Context, tx *gorm.DB, userID uuid.UUID, action, entityType string, entityID uuid.UUID, changes interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.logs = append(m.logs, domain.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
	})
	return nil
}

type p0MockReportRepo struct {
	repository.ReportRepository
	mu        sync.Mutex
	reports   map[uuid.UUID]*domain.CitizenReport
	photos    map[string]*domain.ReportPhoto
	locations map[uuid.UUID]uuid.UUID // reportID -> authorizedLocationID
}

func newP0MockReportRepo() *p0MockReportRepo {
	return &p0MockReportRepo{
		reports:   make(map[uuid.UUID]*domain.CitizenReport),
		photos:    make(map[string]*domain.ReportPhoto),
		locations: make(map[uuid.UUID]uuid.UUID),
	}
}

func (m *p0MockReportRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	rep, ok := m.reports[id]
	if !ok {
		return nil, errors.New("laporan tidak ditemukan")
	}
	return rep, nil
}

func (m *p0MockReportRepo) UpdateReportStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, newStatus string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	rep, ok := m.reports[id]
	if !ok {
		return errors.New("laporan tidak ditemukan")
	}
	rep.Status = newStatus
	return nil
}

func (m *p0MockReportRepo) GetPhotoByStorageKey(ctx context.Context, key string) (*domain.ReportPhoto, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, p := range m.photos {
		if strings.Contains(k, key) {
			return p, nil
		}
	}
	return nil, errors.New("photo not found")
}

func (m *p0MockReportRepo) GetPEResultByPhotoURL(ctx context.Context, url string) (*domain.PEResult, error) {
	return nil, errors.New("pe photo not found")
}

func (m *p0MockReportRepo) VerifyReportLocation(ctx context.Context, reportID uuid.UUID, locationID uuid.UUID) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	authLoc, ok := m.locations[reportID]
	if !ok {
		return false, nil
	}
	return authLoc == locationID, nil
}

// =========================================================================
// P0-01: INVENTORY ACCOUNTING INTEGRITY TESTS
// =========================================================================

func TestP0_01_InventoryAccountingIntegrity(t *testing.T) {
	ctx := context.Background()
	invRepo := newP0MockInventoryRepo()
	auditRepo := &p0MockAuditRepo{}
	usecaseInv := usecase.NewInventoryUsecase(invRepo, auditRepo)

	actorID := uuid.New()

	// 1. Initial State: Central Warehouse has 50L Insecticide, Puskesmas has 0L
	dinkesItem := &domain.InventoryItem{
		ItemID:       uuid.New(),
		ItemType:     "INSEKTISIDA",
		Warehouse:    "Gudang Utama Dinkes",
		Quantity:     50.0,
		MinimumStock: 10.0,
	}
	_ = invRepo.CreateItem(ctx, dinkesItem)

	puskesmasItem := &domain.InventoryItem{
		ItemID:       uuid.New(),
		ItemType:     "INSEKTISIDA",
		Warehouse:    "Puskesmas Pancoran Mas",
		Quantity:     0.0,
		MinimumStock: 5.0,
	}
	_ = invRepo.CreateItem(ctx, puskesmasItem)

	// TEST A: Over-deduction / Direct reduction beyond stock MUST FAIL (No virtual stock / No silent cross-allocation)
	t.Run("A_OverDeduction_MustFail", func(t *testing.T) {
		err := usecaseInv.UpdateStock(ctx, puskesmasItem.ItemID, actorID, -10.0, "Penggunaan langsung")
		if err == nil {
			t.Fatalf("EXPECTED ERROR for over-deduction beyond stock, but got nil (Virtual stock created!)")
		}
		if !strings.Contains(err.Error(), "stok tidak mencukupi") {
			t.Fatalf("Expected 'stok tidak mencukupi', got: %v", err)
		}
		// Invariant: Puskesmas stock must remain 0
		item, _ := invRepo.GetItemByID(ctx, puskesmasItem.ItemID)
		if item.Quantity != 0.0 {
			t.Fatalf("Expected stock 0, got %f", item.Quantity)
		}
	})

	// TEST B: StockRequest -> Approve -> Transfer IN_TRANSIT -> Accounting conservation
	var requestID uuid.UUID
	t.Run("B_StockRequest_Approval_InTransit", func(t *testing.T) {
		req := &domain.StockRequest{
			ID:        uuid.New(),
			ItemID:    puskesmasItem.ItemID,
			Warehouse: puskesmasItem.Warehouse,
			Quantity:  20.0,
			Status:    "PENDING",
			Notes:     "Permintaan darurat lonjakan kasus DBD",
		}
		_ = invRepo.CreateStockRequest(ctx, req)
		requestID = req.ID

		err := usecaseInv.ApproveStockRequest(ctx, actorID, requestID, "APPROVED")
		if err != nil {
			t.Fatalf("ApproveStockRequest failed: %v", err)
		}

		// Source stock MUST be deducted by 20.0
		src, _ := invRepo.GetItemByID(ctx, dinkesItem.ItemID)
		if src.Quantity != 30.0 {
			t.Fatalf("Expected source stock to be 30.0, got: %f", src.Quantity)
		}

		// Destination stock MUST NOT be added yet (Wait for physical receipt)
		dst, _ := invRepo.GetItemByID(ctx, puskesmasItem.ItemID)
		if dst.Quantity != 0.0 {
			t.Fatalf("Expected destination stock to still be 0.0 (in-transit), got: %f", dst.Quantity)
		}

		// Verify transfer record created with IN_TRANSIT status
		if len(invRepo.transfers) != 1 {
			t.Fatalf("Expected 1 transfer record created, got %d", len(invRepo.transfers))
		}
		for _, tr := range invRepo.transfers {
			if tr.Status != "IN_TRANSIT" || tr.Quantity != 20.0 {
				t.Fatalf("Invalid transfer state: %+v", tr)
			}
		}

		// Invariant check: Source (30) + Destination (0) + In-Transit (20) = Total Accounted Stock (50)
		totalAccounted := src.Quantity + dst.Quantity + 20.0
		if totalAccounted != 50.0 {
			t.Fatalf("Stock Accounting Invariant Violated! Total: %f, Expected: 50.0", totalAccounted)
		}
	})

	// TEST C: Prevent Duplicate Approval (Atomic Status Enforcement)
	t.Run("C_DuplicateApproval_MustFail", func(t *testing.T) {
		err := usecaseInv.ApproveStockRequest(ctx, actorID, requestID, "APPROVED")
		if err == nil {
			t.Fatalf("Expected duplicate ApproveStockRequest to fail, but succeeded!")
		}
	})

	// TEST D: Puskesmas Receives Transfer -> Destination Credited & Transfer COMPLETED
	t.Run("D_ReceiveTransfer_Completion", func(t *testing.T) {
		var transferID uuid.UUID
		for id := range invRepo.transfers {
			transferID = id
			break
		}

		err := usecaseInv.ReceiveStockTransfer(ctx, actorID, transferID, true, "Barang diterima dalam kondisi baik")
		if err != nil {
			t.Fatalf("ReceiveStockTransfer failed: %v", err)
		}

		// Destination stock MUST now be 20.0
		dst, _ := invRepo.GetItemByID(ctx, puskesmasItem.ItemID)
		if dst.Quantity != 20.0 {
			t.Fatalf("Expected destination stock 20.0, got: %f", dst.Quantity)
		}

		// Transfer status MUST be COMPLETED
		tr, _ := invRepo.GetStockTransferByID(ctx, transferID)
		if tr.Status != "COMPLETED" {
			t.Fatalf("Expected transfer COMPLETED, got %s", tr.Status)
		}

		// Invariant check: Source (30) + Destination (20) = 50.0 (Conservation)
		src, _ := invRepo.GetItemByID(ctx, dinkesItem.ItemID)
		if src.Quantity+dst.Quantity != 50.0 {
			t.Fatalf("Stock Accounting Invariant Violated! Total: %f", src.Quantity+dst.Quantity)
		}
	})

	// TEST E: Reject Transfer -> Source Stock Refunded & Destination Unchanged
	t.Run("E_RejectTransfer_RefundSource", func(t *testing.T) {
		// Create another transfer: Dinkes (30) -> Puskesmas (20) for 10 units
		tr2 := &domain.StockTransfer{
			ID:            uuid.New(),
			ItemID:        dinkesItem.ItemID,
			FromWarehouse: dinkesItem.Warehouse,
			ToWarehouse:   puskesmasItem.Warehouse,
			Quantity:      10.0,
			Status:        "IN_TRANSIT",
		}
		_ = invRepo.CreateStockTransfer(ctx, nil, tr2)
		_ = invRepo.UpdateStock(ctx, nil, dinkesItem.ItemID, -10.0) // Dinkes now 20.0

		// Recipient rejects transfer (damaged goods / wrong batch)
		err := usecaseInv.ReceiveStockTransfer(ctx, actorID, tr2.ID, false, "Kemasan rusak saat pengiriman")
		if err != nil {
			t.Fatalf("ReceiveStockTransfer (reject) failed: %v", err)
		}

		// Source stock MUST be refunded by 10.0 (Back to 30.0)
		src, _ := invRepo.GetItemByID(ctx, dinkesItem.ItemID)
		if src.Quantity != 30.0 {
			t.Fatalf("Expected refunded source stock 30.0, got: %f", src.Quantity)
		}

		// Destination stock MUST remain 20.0
		dst, _ := invRepo.GetItemByID(ctx, puskesmasItem.ItemID)
		if dst.Quantity != 20.0 {
			t.Fatalf("Expected destination stock 20.0, got: %f", dst.Quantity)
		}

		// Invariant: Source (30) + Destination (20) = 50.0
		if src.Quantity+dst.Quantity != 50.0 {
			t.Fatalf("Stock Accounting Invariant Violated post-refund! Total: %f", src.Quantity+dst.Quantity)
		}
	})
}

// =========================================================================
// P0-02: PROTECT SENSITIVE FILE STORAGE TESTS
// =========================================================================

func TestP0_02_SensitiveFileStorageProtection(t *testing.T) {
	ctx := context.Background()
	reportRepo := newP0MockReportRepo()
	auditRepo := &p0MockAuditRepo{}

	reportUsecase := usecase.NewReportUsecase(reportRepo, auditRepo, nil, &storage.MockStorageService{
		ExistsFunc: func(ctx context.Context, bucketName, objectName string) (bool, error) {
			if objectName == "non_existent_file_123.jpg" {
				return false, nil
			}
			return true, nil
		},
	})

	// Create a dummy test file in backend/storage/photos
	testDir := filepath.Join("storage", "photos")
	_ = os.MkdirAll(testDir, 0755)
	dummyFilename := "test_sensitive_patient_doc.jpg"
	dummyFilePath := filepath.Join(testDir, dummyFilename)
	_ = os.WriteFile(dummyFilePath, []byte("sensitive patient data content"), 0644)
	defer os.Remove(dummyFilePath)

	// Register in Mock DB
	reportID := uuid.New()
	correctLocationID := uuid.New()
	wrongLocationID := uuid.New()

	reportRepo.reports[reportID] = &domain.CitizenReport{
		ID:     reportID,
		Status: "MENUNGGU_PE_KADER",
	}
	reportRepo.photos[dummyFilename] = &domain.ReportPhoto{
		ID:         uuid.New(),
		ReportID:   reportID,
		StorageKey: "storage/photos/" + dummyFilename,
	}
	reportRepo.locations[reportID] = correctLocationID

	// 1. Path Traversal Attempt MUST FAIL (400)
	t.Run("PathTraversal_MustFail", func(t *testing.T) {
		traversalFilenames := []string{
			"../photos/secret.jpg",
			"..\\photos\\secret.jpg",
			"../../etc/passwd",
			"photos/../../config.json",
		}
		for _, name := range traversalFilenames {
			_, err := reportUsecase.GetProtectedPhotoPath(ctx, domain.RoleAdminDeveloper, nil, name)
			if err == nil || !strings.Contains(err.Error(), "path traversal") {
				t.Fatalf("Expected path traversal error for %s, got: %v", name, err)
			}
		}
	})

	// 2. Non-Existent File MUST Return Not Found (404)
	t.Run("NonExistentFile_MustFail", func(t *testing.T) {
		_, err := reportUsecase.GetProtectedPhotoPath(ctx, domain.RoleAdminDeveloper, nil, "non_existent_file_123.jpg")
		if err == nil || !strings.Contains(err.Error(), "tidak ditemukan") {
			t.Fatalf("Expected 'tidak ditemukan' error, got: %v", err)
		}
	})

	// 3. Super Admin & Dinkes Have City-Wide Access (200)
	t.Run("AdminAndDinkes_AccessGranted", func(t *testing.T) {
		path, err := reportUsecase.GetProtectedPhotoPath(ctx, domain.RoleAdminDeveloper, nil, dummyFilename)
		if err != nil || path == "" {
			t.Fatalf("Super Admin access should succeed, got: %v", err)
		}

		pathDinkes, err := reportUsecase.GetProtectedPhotoPath(ctx, domain.RoleDinkes, nil, dummyFilename)
		if err != nil || pathDinkes == "" {
			t.Fatalf("Dinkes access should succeed, got: %v", err)
		}
	})

	// 4. Authorized Officer with Matching Location Has Access (200)
	t.Run("AuthorizedLocation_AccessGranted", func(t *testing.T) {
		path, err := reportUsecase.GetProtectedPhotoPath(ctx, domain.RolePetugasPuskesmas, &correctLocationID, dummyFilename)
		if err != nil || path == "" {
			t.Fatalf("Authorized officer with matching location should succeed, got: %v", err)
		}
	})

	// 5. Authorized Officer with Wrong Location Scope MUST FAIL (403)
	t.Run("WrongLocation_AccessDenied", func(t *testing.T) {
		_, err := reportUsecase.GetProtectedPhotoPath(ctx, domain.RolePetugasPuskesmas, &wrongLocationID, dummyFilename)
		if err == nil || !strings.Contains(err.Error(), "akses ditolak") {
			t.Fatalf("Expected access denied for wrong location officer, got: %v", err)
		}
	})

	// 6. Anonymous / Missing Location User MUST FAIL (401/403)
	t.Run("MissingLocationOfficer_AccessDenied", func(t *testing.T) {
		_, err := reportUsecase.GetProtectedPhotoPath(ctx, domain.RolePetugasPuskesmas, nil, dummyFilename)
		if err == nil || !strings.Contains(err.Error(), "akses ditolak") {
			t.Fatalf("Expected access denied for officer without location, got: %v", err)
		}
	})
}

// =========================================================================
// P0-03 & P0-04: DSS ADVISORY OVERRIDE & SEGREGATION OF DUTIES (SoD)
// =========================================================================

type p0MockPERepo struct {
	repository.PERepository
	mu       sync.Mutex
	peResult *domain.PEResult
	reports  map[uuid.UUID]*domain.CitizenReport
}

func (m *p0MockPERepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (m *p0MockPERepo) GetPEResultByReportID(ctx context.Context, reportID uuid.UUID) (*domain.PEResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.peResult != nil && m.peResult.ReportID == reportID {
		return m.peResult, nil
	}
	return nil, errors.New("pe result not found")
}

func (m *p0MockPERepo) CreatePEResult(ctx context.Context, tx *gorm.DB, pe *domain.PEResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if pe.ID == uuid.Nil {
		pe.ID = uuid.New()
	}
	m.peResult = pe
	return nil
}

func (m *p0MockPERepo) GetPEResultByID(ctx context.Context, tx *gorm.DB, peID uuid.UUID) (*domain.PEResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.peResult != nil && m.peResult.ID == peID {
		if m.reports != nil {
			if rep, ok := m.reports[m.peResult.ReportID]; ok && rep != nil {
				m.peResult.Report = *rep
			}
		}
		return m.peResult, nil
	}
	return nil, errors.New("pe result not found")
}

func (m *p0MockPERepo) UpdatePEResult(ctx context.Context, tx *gorm.DB, pe *domain.PEResult) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.peResult = pe
	return nil
}

type p0MockFoggingRepo struct {
	repository.FoggingRepository
	mu          sync.Mutex
	foggingReqs []*domain.FoggingRequest
}

func (m *p0MockFoggingRepo) CreateFoggingRequest(ctx context.Context, tx *gorm.DB, req *domain.FoggingRequest) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.foggingReqs = append(m.foggingReqs, req)
	return nil
}

func TestP0_03_and_04_DSS_and_SoD(t *testing.T) {
	ctx := context.Background()
	reportRepo := newP0MockReportRepo()
	peRepo := &p0MockPERepo{reports: reportRepo.reports}
	foggingRepo := &p0MockFoggingRepo{}
	invRepo := newP0MockInventoryRepo()
	auditRepo := &p0MockAuditRepo{}

	peUsecase := usecase.NewPEUsecase(peRepo, reportRepo, foggingRepo, invRepo, auditRepo, nil, nil)

	investigatorKaderID := uuid.New()
	validatorKepalaPuskesmasID := uuid.New()
	reportID := uuid.New()

	reportRepo.reports[reportID] = &domain.CitizenReport{
		ID:             reportID,
		Status:         "MENUNGGU_PE_KADER",
		InvestigatorID: &investigatorKaderID,
		Location: domain.ReportLocation{
			Kelurahan: "Pancoran Mas",
			Kecamatan: "Pancoran Mas",
		},
	}

	// Step 1: Investigator submits PE with No Larvae & Fever < 3 -> DSS recommends EDUKASI_PSN
	peID, err := peUsecase.SubmitPEResult(
		ctx,
		reportID,
		investigatorKaderID,
		100,   // radius
		false, // larvae NOT found
		1,     // fever < 3
		true,  // psn done
		"Tidak ditemukan jentik di bak mandi",
		"http://evidence.jpg",
	)
	if err != nil {
		t.Fatalf("SubmitPEResult failed: %v", err)
	}

	if peRepo.peResult.Recommendation != "EDUKASI_PSN" {
		t.Fatalf("Expected DSS recommendation EDUKASI_PSN, got: %s", peRepo.peResult.Recommendation)
	}

	// ---------------------------------------------------------------------
	// P0-04: SEGREGATION OF DUTIES (Four-Eyes Principle)
	// ---------------------------------------------------------------------
	t.Run("P0_04_SoD_SameActor_Validation_MustFail", func(t *testing.T) {
		// Investigator tries to validate their own PE result
		err := peUsecase.ValidatePEResult(ctx, *peID, investigatorKaderID, true, "", "", "")
		if err == nil {
			t.Fatalf("CRITICAL SECURITY FLAW: Investigator was allowed to validate their own PE (SoD Violated)!")
		}
		if !strings.Contains(err.Error(), "pemisahan tugas") && !strings.Contains(err.Error(), "Four-Eyes Principle") {
			t.Fatalf("Expected Four-Eyes Principle error, got: %v", err)
		}
	})

	// ---------------------------------------------------------------------
	// P0-03: DSS ADVISORY & HUMAN OVERRIDE
	// ---------------------------------------------------------------------
	t.Run("P0_03_DSS_Override_WithoutReason_MustFail", func(t *testing.T) {
		// Different authorized validator overrides DSS recommendation but omits reason
		err := peUsecase.ValidatePEResult(ctx, *peID, validatorKepalaPuskesmasID, true, "", "FOGGING_FOKUS", "")
		if err == nil {
			t.Fatalf("Expected override without reason to fail, but succeeded!")
		}
		if !strings.Contains(err.Error(), "alasan override") {
			t.Fatalf("Expected 'alasan override' error, got: %v", err)
		}
	})

	t.Run("P0_03_DSS_Override_WithReason_Success", func(t *testing.T) {
		// Prepare stock in puskesmas warehouse so FOGGING_FOKUS can be approved
		_ = invRepo.CreateItem(ctx, &domain.InventoryItem{
			ItemID:    uuid.New(),
			ItemType:  "INSEKTISIDA",
			Warehouse: "Pancoran Mas",
			Quantity:  25.0,
		})

		overrideReason := "Hasil verifikasi surveilans puskesmas menemukan kasus sekunder di RT sebelah, diputuskan fogging fokus"
		err := peUsecase.ValidatePEResult(ctx, *peID, validatorKepalaPuskesmasID, true, "", "FOGGING_FOKUS", overrideReason)
		if err != nil {
			t.Fatalf("Human override with reason should succeed, got: %v", err)
		}

		// Verify final recommendation updated to FOGGING_FOKUS
		if peRepo.peResult.Recommendation != "FOGGING_FOKUS" {
			t.Fatalf("Expected recommendation updated to FOGGING_FOKUS, got: %s", peRepo.peResult.Recommendation)
		}

		// Verify FoggingRequest was created
		if len(foggingRepo.foggingReqs) != 1 {
			t.Fatalf("Expected FoggingRequest created post-override, got %d", len(foggingRepo.foggingReqs))
		}

		// Verify Report status transitioned to WAITING_FOGGING_APPROVAL
		rep := reportRepo.reports[reportID]
		if rep.Status != "WAITING_FOGGING_APPROVAL" {
			t.Fatalf("Expected status WAITING_FOGGING_APPROVAL, got: %s", rep.Status)
		}

		// Verify Audit Trail has OVERRIDE_DSS_RECOMMENDATION
		foundAudit := false
		for _, log := range auditRepo.logs {
			if log.Action == "OVERRIDE_DSS_RECOMMENDATION" {
				foundAudit = true
				break
			}
		}
		if !foundAudit {
			t.Fatalf("Expected OVERRIDE_DSS_RECOMMENDATION in audit log, but not found!")
		}
	})
}

// =========================================================================
// P0-06: RESIDUAL BUG REGRESSION: APPEAL -> ESKALASI -> TAKEOVER -> PE
// =========================================================================

type p0MockVerificationRepo struct {
	repository.VerificationRepository
	mu            sync.Mutex
	verifications []*domain.Verification
}

func (m *p0MockVerificationRepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (m *p0MockVerificationRepo) Create(ctx context.Context, tx *gorm.DB, v *domain.Verification) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.verifications = append(m.verifications, v)
	return nil
}

func (m *p0MockReportRepo) AssignInvestigator(ctx context.Context, id uuid.UUID, investigatorID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	rep, ok := m.reports[id]
	if !ok {
		return errors.New("report not found")
	}
	rep.InvestigatorID = &investigatorID
	return nil
}

func (m *p0MockReportRepo) SubmitAppealAtomic(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	rep, ok := m.reports[id]
	if !ok {
		return errors.New("report not found")
	}
	rep.Status = "APPEALED"
	return nil
}

func TestP0_06_AppealToPELifecycle(t *testing.T) {
	ctx := context.Background()
	reportRepo := newP0MockReportRepo()
	auditRepo := &p0MockAuditRepo{}
	verifRepo := &p0MockVerificationRepo{}
	peRepo := &p0MockPERepo{reports: reportRepo.reports}
	foggingRepo := &p0MockFoggingRepo{}
	invRepo := newP0MockInventoryRepo()

	reportUsecase := usecase.NewReportUsecase(reportRepo, auditRepo, nil, &storage.MockStorageService{})
	verifUsecase := usecase.NewVerificationUsecase(verifRepo, reportRepo, auditRepo, nil, nil)
	peUsecase := usecase.NewPEUsecase(peRepo, reportRepo, foggingRepo, invRepo, auditRepo, nil, nil)

	reportID := uuid.New()
	puskesmasOfficerID := uuid.New()
	kepalaPuskesmasID := uuid.New()
	puskesmasLocID := uuid.New()

	// Initial State: Report rejected as HOAX by Kader
	reportRepo.reports[reportID] = &domain.CitizenReport{
		ID:     reportID,
		Status: "REJECTED_HOAX",
		Location: domain.ReportLocation{
			Kelurahan: "Pancoran Mas",
			Kecamatan: "Pancoran Mas",
		},
	}
	reportRepo.locations[reportID] = puskesmasLocID

	// 1. Citizen appeals the rejection
	t.Log("Step 1: Citizen appeals rejection")
	err := reportUsecase.SubmitAppeal(ctx, reportID, "Anggota keluarga dirawat di RS positif DBD, bukti lab terlampir")
	if err != nil {
		t.Fatalf("SubmitAppeal failed: %v", err)
	}
	if reportRepo.reports[reportID].Status != "APPEALED" {
		t.Fatalf("Expected status APPEALED, got: %s", reportRepo.reports[reportID].Status)
	}

	// 2. Puskesmas reviews and APPROVES the appeal
	t.Log("Step 2: Puskesmas verifies and approves appeal")
	err = verifUsecase.VerifyReport(ctx, reportID, puskesmasOfficerID, domain.RolePetugasPuskesmas, &puskesmasLocID, "APPROVED", "Bukti laboratorium valid")
	if err != nil {
		t.Fatalf("VerifyReport failed: %v", err)
	}

	// Status MUST escalate to ESKALASI_PUSKESMAS
	if reportRepo.reports[reportID].Status != "ESKALASI_PUSKESMAS" {
		t.Fatalf("Regression detected! Expected ESKALASI_PUSKESMAS, got: %s", reportRepo.reports[reportID].Status)
	}

	// 3. Puskesmas officer takes over the escalated report
	t.Log("Step 3: Puskesmas officer takes over investigation")
	err = reportUsecase.TakeoverReport(ctx, reportID, puskesmasOfficerID, "Pengambilan alih investigasi pasca banding")
	if err != nil {
		t.Fatalf("TakeoverReport failed: %v", err)
	}
	if reportRepo.reports[reportID].InvestigatorID == nil || *reportRepo.reports[reportID].InvestigatorID != puskesmasOfficerID {
		t.Fatalf("Expected investigator to be Puskesmas officer")
	}

	// 4. Puskesmas officer submits PE result
	t.Log("Step 4: Puskesmas officer submits PE result")
	peID, err := peUsecase.SubmitPEResult(
		ctx,
		reportID,
		puskesmasOfficerID,
		100,
		true, // Larvae found
		3,    // Fever count >= 3 -> FOGGING_FOKUS
		true,
		"Ditemukan banyak jentik aedes aegypti",
		"http://evidence.jpg",
	)
	if err != nil {
		t.Fatalf("SubmitPEResult failed: %v", err)
	}
	if reportRepo.reports[reportID].Status != "MENUNGGU_VALIDASI_PUSKESMAS" {
		t.Fatalf("Expected status MENUNGGU_VALIDASI_PUSKESMAS, got: %s", reportRepo.reports[reportID].Status)
	}

	// 5. SoD verification: Puskesmas officer tries to validate their own PE -> MUST BE BLOCKED
	t.Log("Step 5: SoD check on PE validation")
	err = peUsecase.ValidatePEResult(ctx, *peID, puskesmasOfficerID, true, "", "", "")
	if err == nil {
		t.Fatalf("SoD failed: Officer validated their own PE!")
	}

	// 6. Legitimate Kepala Puskesmas validates PE
	t.Log("Step 6: Kepala Puskesmas validates PE")
	_ = invRepo.CreateItem(ctx, &domain.InventoryItem{
		ItemID:    uuid.New(),
		ItemType:  "INSEKTISIDA",
		Warehouse: "Pancoran Mas",
		Quantity:  30.0,
	})
	err = peUsecase.ValidatePEResult(ctx, *peID, kepalaPuskesmasID, true, "", "", "")
	if err != nil {
		t.Fatalf("Kepala Puskesmas validation failed: %v", err)
	}

	// Final Status MUST be WAITING_FOGGING_APPROVAL
	if reportRepo.reports[reportID].Status != "WAITING_FOGGING_APPROVAL" {
		t.Fatalf("Expected final status WAITING_FOGGING_APPROVAL, got: %s", reportRepo.reports[reportID].Status)
	}
	t.Log("Full lifecycle successfully verified without deadlocks or regressions!")
}
