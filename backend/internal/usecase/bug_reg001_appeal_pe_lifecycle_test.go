package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
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

// In-memory mock for full lifecycle verification
type mockFullLifecycleReportRepo struct {
	repository.ReportRepository
	report *domain.CitizenReport
}

func (m *mockFullLifecycleReportRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	if m.report == nil || m.report.ID != id {
		return nil, errors.New("report not found")
	}
	return m.report, nil
}

func (m *mockFullLifecycleReportRepo) UpdateReportStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	if m.report == nil || m.report.ID != id {
		return errors.New("report not found")
	}
	m.report.Status = status
	return nil
}

func (m *mockFullLifecycleReportRepo) AssignInvestigator(ctx context.Context, id uuid.UUID, investigatorID uuid.UUID) error {
	if m.report == nil || m.report.ID != id {
		return errors.New("report not found")
	}
	m.report.InvestigatorID = &investigatorID
	return nil
}

func (m *mockFullLifecycleReportRepo) SubmitAppealAtomic(ctx context.Context, id uuid.UUID) error {
	if m.report == nil || m.report.ID != id {
		return errors.New("report not found")
	}
	if (m.report.Status != "REJECTED" && m.report.Status != "REJECTED_HOAX") || m.report.AppealCount >= 1 {
		return errors.New("laporan tidak dapat dibanding")
	}
	m.report.Status = "APPEALED"
	m.report.AppealCount++
	return nil
}

func (m *mockFullLifecycleReportRepo) VerifyReportLocation(ctx context.Context, reportID uuid.UUID, locationID uuid.UUID) (bool, error) {
	return true, nil
}

type mockVerificationRepo struct {
	repository.VerificationRepository
	verifications []*domain.Verification
}

func (m *mockVerificationRepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (m *mockVerificationRepo) Create(ctx context.Context, tx *gorm.DB, v *domain.Verification) error {
	m.verifications = append(m.verifications, v)
	return nil
}

type mockAuditRepoLifecycle struct {
	repository.AuditRepository
	logs []string
}

func (m *mockAuditRepoLifecycle) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action, entity string, entityID uuid.UUID, details interface{}) error {
	m.logs = append(m.logs, action)
	return nil
}

type mockPELifecycleRepo struct {
	repository.PERepository
	peResult *domain.PEResult
}

func (m *mockPELifecycleRepo) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

func (m *mockPELifecycleRepo) GetPEResultByReportID(ctx context.Context, reportID uuid.UUID) (*domain.PEResult, error) {
	if m.peResult != nil && m.peResult.ReportID == reportID {
		return m.peResult, nil
	}
	return nil, errors.New("pe not found")
}

func (m *mockPELifecycleRepo) CreatePEResult(ctx context.Context, tx *gorm.DB, pe *domain.PEResult) error {
	pe.ID = uuid.New()
	m.peResult = pe
	return nil
}

type mockNotifLifecycleUsecase struct {
	usecase.NotificationUsecase
}

func (m *mockNotifLifecycleUsecase) NotifyEWSKader(ctx context.Context, kelurahan, officerID string) error {
	return nil
}

type mockEwsLifecycleUsecase struct {
	usecase.EwsKlbUsecase
}

func (m *mockEwsLifecycleUsecase) CheckAndTriggerKLB(ctx context.Context, areaName, scope string) error {
	return nil
}

type mockMasterDataLifecycleRepo struct {
	repository.MasterDataRepository
}

func (m *mockMasterDataLifecycleRepo) GetKelurahanByName(ctx context.Context, name string) (*domain.MasterKelurahan, error) {
	return &domain.MasterKelurahan{
		Name: name,
		Puskesmas: domain.MasterPuskesmas{
			Name: "Puskesmas Pancoran Mas",
		},
	}, nil
}

type mockInventoryLifecycleRepo struct {
	repository.InventoryRepository
}

func (m *mockInventoryLifecycleRepo) GetItemByTypeAndWarehouse(ctx context.Context, itemType, warehouse string) (*domain.InventoryItem, error) {
	return &domain.InventoryItem{
		ItemID:           uuid.New(),
		ItemType:         itemType,
		Warehouse:        warehouse,
		Quantity:         50.0,
		ReservedQuantity: 0.0,
	}, nil
}

type mockFoggingLifecycleRepo struct {
	repository.FoggingRepository
}

func (m *mockFoggingLifecycleRepo) CreateFoggingRequest(ctx context.Context, tx *gorm.DB, req *domain.FoggingRequest) error {
	return nil
}

func TestBUGREG001_FullLifecycle_AppealToPEValidation(t *testing.T) {
	ctx := context.Background()

	reportID := uuid.New()
	kaderID := uuid.New()
	puskesmasOfficerID := uuid.New()
	puskesmasLocationID := uuid.New()

	// Initial State: Report is REJECTED_HOAX by Kader
	report := &domain.CitizenReport{
		ID:             reportID,
		Status:         "REJECTED_HOAX",
		ReferenceNo:    "REP-2026-001",
		InvestigatorID: &kaderID,
		AppealCount:    0,
		Location: domain.ReportLocation{
			Kelurahan: "Depok",
			Kecamatan: "Pancoran Mas",
		},
		CreatedAt: time.Now(),
	}

	reportRepo := &mockFullLifecycleReportRepo{report: report}
	auditRepo := &mockAuditRepoLifecycle{}
	verificationRepo := &mockVerificationRepo{}
	peRepo := &mockPELifecycleRepo{}
	notifUsecase := &mockNotifLifecycleUsecase{}
	ewsUsecase := &mockEwsLifecycleUsecase{}
	masterDataRepo := &mockMasterDataLifecycleRepo{}
	invRepo := &mockInventoryLifecycleRepo{}
	foggingRepo := &mockFoggingLifecycleRepo{}

	reportUsecase := usecase.NewReportUsecase(reportRepo, auditRepo, nil, &storage.MockStorageService{})
	verificationUsecase := usecase.NewVerificationUsecase(verificationRepo, reportRepo, auditRepo, ewsUsecase, notifUsecase)
	peUsecase := usecase.NewPEUsecase(peRepo, reportRepo, foggingRepo, invRepo, auditRepo, masterDataRepo, ewsUsecase)

	// Step 1: Citizen appeals against REJECTED_HOAX
	t.Log("Step 1: Citizen appeals against REJECTED_HOAX")
	err := reportUsecase.SubmitAppeal(ctx, reportID, "Mohon ditinjau ulang, anggota keluarga positif DBD di RS")
	if err != nil {
		t.Fatalf("SubmitAppeal failed: %v", err)
	}
	if report.Status != "APPEALED" {
		t.Fatalf("Expected report status APPEALED, got %s", report.Status)
	}

	// Step 2: Puskesmas reviews and APPROVES the appeal
	t.Log("Step 2: Puskesmas reviews and APPROVES the appeal")
	err = verificationUsecase.VerifyReport(
		ctx,
		reportID,
		puskesmasOfficerID,
		domain.RolePetugasPuskesmas,
		&puskesmasLocationID,
		"APPROVED",
		"Hasil lab RS valid, eskalasi untuk PE Puskesmas",
	)
	if err != nil {
		t.Fatalf("VerifyReport by Puskesmas failed: %v", err)
	}

	// Invariant Check: Status MUST be ESKALASI_PUSKESMAS (not VERIFIED or orphan state)
	if report.Status != "ESKALASI_PUSKESMAS" {
		t.Fatalf("BUG-REG-001 Detected! Expected status ESKALASI_PUSKESMAS, got %s", report.Status)
	}

	// Step 3: Puskesmas officer takes over the escalated report
	t.Log("Step 3: Puskesmas officer takes over the escalated report")
	err = reportUsecase.TakeoverReport(ctx, reportID, puskesmasOfficerID, "Puskesmas takes over investigation post-appeal")
	if err != nil {
		t.Fatalf("TakeoverReport failed on ESKALASI_PUSKESMAS: %v", err)
	}
	if report.InvestigatorID == nil || *report.InvestigatorID != puskesmasOfficerID {
		t.Fatalf("Expected investigator to be Puskesmas officer %v, got %v", puskesmasOfficerID, report.InvestigatorID)
	}

	// Step 4: Puskesmas officer submits PE result
	t.Log("Step 4: Puskesmas officer submits PE result")
	peID, err := peUsecase.SubmitPEResult(
		ctx,
		reportID,
		puskesmasOfficerID,
		100,  // radius
		true, // larva found
		3,    // fever count >= 3 -> recommendation FOGGING_FOKUS
		true, // psn edukasi
		"Ditemukan jentik di penampungan air dan 3 penderita demam",
		"http://evidence.jpg",
	)
	if err != nil {
		t.Fatalf("SubmitPEResult failed after takeover: %v", err)
	}
	if peID == nil {
		t.Fatalf("Expected non-nil peID")
	}
	if report.Status != "MENUNGGU_VALIDASI_PUSKESMAS" {
		t.Fatalf("Expected report status MENUNGGU_VALIDASI_PUSKESMAS, got %s", report.Status)
	}

	// Step 5: Puskesmas validates PE (simulating validation approval)
	t.Log("Step 5: Puskesmas validates PE result")
	if peRepo.peResult == nil {
		t.Fatalf("Expected PE result to be created")
	}
	if peRepo.peResult.Recommendation != "FOGGING_FOKUS" {
		t.Fatalf("Expected recommendation FOGGING_FOKUS, got %s", peRepo.peResult.Recommendation)
	}

	t.Log("Full lifecycle completed successfully: REJECTED_HOAX -> APPEALED -> ESKALASI_PUSKESMAS -> TAKEOVER -> SUBMIT_PE -> PENDING_VALIDATION")
}
