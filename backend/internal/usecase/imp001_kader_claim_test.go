package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// mockIMP001Repo implements repository.ReportRepository for IMP-001 unit tests
type mockIMP001Repo struct {
	repository.ReportRepository
	mu             sync.Mutex
	report         *domain.CitizenReport
	reportLocID    uuid.UUID
	claimAtomicErr error
}

func (m *mockIMP001Repo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.report == nil || m.report.ID != id {
		return nil, errors.New("record not found")
	}
	// Return a copy
	cp := *m.report
	return &cp, nil
}

func (m *mockIMP001Repo) VerifyReportLocation(ctx context.Context, reportID uuid.UUID, locationID uuid.UUID) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.report == nil || m.report.ID != reportID {
		return false, nil
	}
	return m.reportLocID == locationID, nil
}

func (m *mockIMP001Repo) ClaimReportAtomic(ctx context.Context, id uuid.UUID, kaderID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.claimAtomicErr != nil {
		return m.claimAtomicErr
	}
	if m.report == nil || m.report.ID != id {
		return errors.New("laporan tidak ditemukan")
	}
	if m.report.Status != "MENUNGGU_PE_KADER" || m.report.InvestigatorID != nil {
		return errors.New("laporan sudah diklaim oleh kader lain atau status tidak valid")
	}
	m.report.InvestigatorID = &kaderID
	return nil
}

// mockIMP001AuditRepo implements repository.AuditRepository for IMP-001 unit tests
type mockIMP001AuditRepo struct {
	repository.AuditRepository
	mu      sync.Mutex
	actions []struct {
		actorID uuid.UUID
		action  string
		entity  string
		refID   uuid.UUID
		details interface{}
	}
}

func (m *mockIMP001AuditRepo) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action, entity string, refID uuid.UUID, details interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.actions = append(m.actions, struct {
		actorID uuid.UUID
		action  string
		entity  string
		refID   uuid.UUID
		details interface{}
	}{
		actorID: actorID,
		action:  action,
		entity:  entity,
		refID:   refID,
		details: details,
	})
	return nil
}

func (m *mockIMP001AuditRepo) countAction(action string) int {
	m.mu.Lock()
	defer m.mu.Unlock()
	count := 0
	for _, a := range m.actions {
		if a.action == action {
			count++
		}
	}
	return count
}

// TEST-01: Same Region Claim -> SUCCESS
func TestIMP001_TEST01_SameRegion_ClaimSuccess(t *testing.T) {
	ctx := context.Background()
	reportID := uuid.New()
	regionBeji := uuid.New()
	kaderID := uuid.New()

	repo := &mockIMP001Repo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionBeji,
	}
	auditRepo := &mockIMP001AuditRepo{}
	u := usecase.NewReportUsecase(repo, auditRepo, nil, &storage.MockStorageService{})

	err := u.ClaimReport(ctx, reportID, kaderID, &regionBeji)
	if err != nil {
		t.Fatalf("expected claim success, got error: %v", err)
	}

	// Verify investigator_id is set
	if repo.report.InvestigatorID == nil || *repo.report.InvestigatorID != kaderID {
		t.Fatalf("expected investigator_id to be %s, got: %v", kaderID, repo.report.InvestigatorID)
	}

	// Verify status remains MENUNGGU_PE_KADER
	if repo.report.Status != "MENUNGGU_PE_KADER" {
		t.Fatalf("expected status MENUNGGU_PE_KADER, got: %s", repo.report.Status)
	}

	// Verify audit trail logged
	if auditRepo.countAction("CLAIM_REPORT") != 1 {
		t.Fatalf("expected exactly 1 CLAIM_REPORT audit log, got: %d", auditRepo.countAction("CLAIM_REPORT"))
	}
	lastAction := auditRepo.actions[0]
	if lastAction.actorID != kaderID || lastAction.refID != reportID {
		t.Fatalf("audit log mismatch: actorID=%v, refID=%v", lastAction.actorID, lastAction.refID)
	}
}

// TEST-02: Cross Region Claim -> FORBIDDEN (No side effects)
func TestIMP001_TEST02_CrossRegion_ClaimForbidden(t *testing.T) {
	ctx := context.Background()
	reportID := uuid.New()
	regionPancoranMas := uuid.New()
	regionBeji := uuid.New()
	kaderBejiID := uuid.New()

	repo := &mockIMP001Repo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionPancoranMas, // Report is in Pancoran Mas
	}
	auditRepo := &mockIMP001AuditRepo{}
	u := usecase.NewReportUsecase(repo, auditRepo, nil, &storage.MockStorageService{})

	// Kader Beji attempts to claim Pancoran Mas report
	err := u.ClaimReport(ctx, reportID, kaderBejiID, &regionBeji)
	if err == nil {
		t.Fatalf("expected error domain.ErrUnauthorizedJurisdiction, got nil")
	}
	if !errors.Is(err, domain.ErrUnauthorizedJurisdiction) {
		t.Fatalf("expected ErrUnauthorizedJurisdiction, got: %v", err)
	}

	// Invariant: investigator_id must remain nil
	if repo.report.InvestigatorID != nil {
		t.Fatalf("side effect detected! investigator_id should remain nil, got: %v", repo.report.InvestigatorID)
	}

	// Invariant: status must not change
	if repo.report.Status != "MENUNGGU_PE_KADER" {
		t.Fatalf("side effect detected! status changed to: %s", repo.report.Status)
	}

	// Invariant: no audit trail written on failed authorization
	if len(auditRepo.actions) != 0 {
		t.Fatalf("side effect detected! audit logs created on failed authorization: %d", len(auditRepo.actions))
	}
}

// TEST-05: Concurrent Claim -> Exactly 1 Winner, 1 Loser
func TestIMP001_TEST05_ConcurrentClaim_AtomicRaceCondition(t *testing.T) {
	ctx := context.Background()
	reportID := uuid.New()
	regionBeji := uuid.New()
	kader1 := uuid.New()
	kader2 := uuid.New()

	repo := &mockIMP001Repo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionBeji,
	}
	auditRepo := &mockIMP001AuditRepo{}
	u := usecase.NewReportUsecase(repo, auditRepo, nil, &storage.MockStorageService{})

	var wg sync.WaitGroup
	successCount := 0
	failCount := 0
	var mu sync.Mutex

	claim := func(kaderID uuid.UUID) {
		defer wg.Done()
		err := u.ClaimReport(ctx, reportID, kaderID, &regionBeji)
		mu.Lock()
		defer mu.Unlock()
		if err == nil {
			successCount++
		} else {
			failCount++
		}
	}

	wg.Add(2)
	go claim(kader1)
	go claim(kader2)
	wg.Wait()

	if successCount != 1 {
		t.Fatalf("expected exactly 1 winner, got: %d", successCount)
	}
	if failCount != 1 {
		t.Fatalf("expected exactly 1 failure, got: %d", failCount)
	}

	// Check investigator_id is either kader1 or kader2, never corrupted
	if repo.report.InvestigatorID == nil || (*repo.report.InvestigatorID != kader1 && *repo.report.InvestigatorID != kader2) {
		t.Fatalf("unexpected investigator_id: %v", repo.report.InvestigatorID)
	}

	// Check exactly 1 audit log created
	if auditRepo.countAction("CLAIM_REPORT") != 1 {
		t.Fatalf("expected exactly 1 CLAIM_REPORT audit log, got: %d", auditRepo.countAction("CLAIM_REPORT"))
	}
}

// TEST-06: Missing or Nil Location Context -> REJECTED SAFELY
func TestIMP001_TEST06_MissingLocationContext_Rejected(t *testing.T) {
	ctx := context.Background()
	reportID := uuid.New()
	regionBeji := uuid.New()
	kaderID := uuid.New()

	repo := &mockIMP001Repo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionBeji,
	}
	auditRepo := &mockIMP001AuditRepo{}
	u := usecase.NewReportUsecase(repo, auditRepo, nil, &storage.MockStorageService{})

	// Actor passes nil locationID
	err := u.ClaimReport(ctx, reportID, kaderID, nil)
	if err == nil {
		t.Fatalf("expected error domain.ErrMissingLocationContext, got nil")
	}
	if !errors.Is(err, domain.ErrMissingLocationContext) {
		t.Fatalf("expected ErrMissingLocationContext, got: %v", err)
	}

	// Invariant: no side effects
	if repo.report.InvestigatorID != nil {
		t.Fatalf("side effect detected! investigator_id changed: %v", repo.report.InvestigatorID)
	}
	if len(auditRepo.actions) != 0 {
		t.Fatalf("side effect detected! audit log created: %d", len(auditRepo.actions))
	}
}

// Additional guardrails: Invalid report status rejected
func TestIMP001_InvalidStatus_Rejected(t *testing.T) {
	ctx := context.Background()
	reportID := uuid.New()
	regionBeji := uuid.New()
	kaderID := uuid.New()

	repo := &mockIMP001Repo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "COMPLETED",
			InvestigatorID: nil,
		},
		reportLocID: regionBeji,
	}
	u := usecase.NewReportUsecase(repo, nil, nil, &storage.MockStorageService{})

	err := u.ClaimReport(ctx, reportID, kaderID, &regionBeji)
	if err == nil {
		t.Fatalf("expected error when status is not MENUNGGU_PE_KADER, got nil")
	}
	if repo.report.InvestigatorID != nil {
		t.Fatalf("investigator_id should remain nil, got: %v", repo.report.InvestigatorID)
	}
}

// Additional guardrails: Already claimed report rejected
func TestIMP001_AlreadyClaimed_Rejected(t *testing.T) {
	ctx := context.Background()
	reportID := uuid.New()
	regionBeji := uuid.New()
	existingKader := uuid.New()
	newKader := uuid.New()

	repo := &mockIMP001Repo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: &existingKader,
		},
		reportLocID: regionBeji,
	}
	u := usecase.NewReportUsecase(repo, nil, nil, &storage.MockStorageService{})

	err := u.ClaimReport(ctx, reportID, newKader, &regionBeji)
	if err == nil {
		t.Fatalf("expected error when report already claimed, got nil")
	}
	// Invariant: existing investigator must not be overwritten
	if *repo.report.InvestigatorID != existingKader {
		t.Fatalf("investigator overwritten! expected: %s, got: %s", existingKader, *repo.report.InvestigatorID)
	}
}
