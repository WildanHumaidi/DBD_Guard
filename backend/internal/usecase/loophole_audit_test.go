package usecase_test

import (
	"context"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockReportRepoForLoopholeAudit struct {
	repository.ReportRepository
	reports   map[uuid.UUID]*domain.CitizenReport
	locations map[uuid.UUID]uuid.UUID
}

func (m *mockReportRepoForLoopholeAudit) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	if r, ok := m.reports[id]; ok {
		return r, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockReportRepoForLoopholeAudit) GetReportByReference(ctx context.Context, refNo string) (*domain.CitizenReport, error) {
	for _, r := range m.reports {
		if r.ReferenceNo == refNo {
			return r, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockReportRepoForLoopholeAudit) UpdateReportStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	if r, ok := m.reports[id]; ok {
		r.Status = status
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *mockReportRepoForLoopholeAudit) SubmitFeedback(ctx context.Context, refNo string, rating int, notes string) error {
	for _, r := range m.reports {
		if r.ReferenceNo == refNo {
			r.FeedbackRating = rating
			r.FeedbackNotes = notes
			return nil
		}
	}
	return gorm.ErrRecordNotFound
}

func (m *mockReportRepoForLoopholeAudit) VerifyReportLocation(ctx context.Context, reportID uuid.UUID, locationID uuid.UUID) (bool, error) {
	authLoc, ok := m.locations[reportID]
	if !ok {
		return false, nil
	}
	return authLoc == locationID, nil
}

func (m *mockReportRepoForLoopholeAudit) MergeDuplicateReportsAtomic(ctx context.Context, childID uuid.UUID, primaryID uuid.UUID, actorID uuid.UUID) error {
	if r, ok := m.reports[childID]; ok {
		r.Status = "MERGED_DUPLICATE"
		r.MergedWithID = &primaryID
		return nil
	}
	return gorm.ErrRecordNotFound
}

type mockAuditRepoForLoopholeAudit struct {
	repository.AuditRepository
	actions []string
}

func (m *mockAuditRepoForLoopholeAudit) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action, entity string, entityID uuid.UUID, details interface{}) error {
	m.actions = append(m.actions, action)
	return nil
}

// =========================================================================
// TEST: Anti-Loophole Regression Suite
// =========================================================================

func TestLoophole_MergeDuplicate_SelfMergePrevented(t *testing.T) {
	reportID := uuid.New()
	kaderID := uuid.New()

	repo := &mockReportRepoForLoopholeAudit{
		reports: map[uuid.UUID]*domain.CitizenReport{
			reportID: {ID: reportID, Status: "MENUNGGU_PE_KADER"},
		},
	}
	audit := &mockAuditRepoForLoopholeAudit{}
	uc := usecase.NewReportUsecase(repo, audit, nil, &storage.MockStorageService{})

	// Trying to merge report into itself MUST fail
	err := uc.MergeDuplicateReport(context.Background(), reportID, reportID, kaderID, "Self merge attempt")
	if err == nil {
		t.Fatalf("Expected error when attempting to merge report with itself, got nil")
	}
}

func TestLoophole_MergeDuplicate_InactivePrimaryPrevented(t *testing.T) {
	childID := uuid.New()
	primaryID := uuid.New()
	kaderID := uuid.New()

	repo := &mockReportRepoForLoopholeAudit{
		reports: map[uuid.UUID]*domain.CitizenReport{
			childID:   {ID: childID, Status: "MENUNGGU_PE_KADER"},
			primaryID: {ID: primaryID, Status: "MERGED_DUPLICATE"}, // Already merged
		},
	}
	audit := &mockAuditRepoForLoopholeAudit{}
	uc := usecase.NewReportUsecase(repo, audit, nil, &storage.MockStorageService{})

	err := uc.MergeDuplicateReport(context.Background(), childID, primaryID, kaderID, "Merge into already merged parent")
	if err == nil {
		t.Fatalf("Expected error when attempting to merge into inactive/merged primary report, got nil")
	}
}

func TestLoophole_SubmitFeedback_DuplicatePrevented(t *testing.T) {
	refNo := "DBD-DPK-FEEDBACK-01"
	reportID := uuid.New()

	repo := &mockReportRepoForLoopholeAudit{
		reports: map[uuid.UUID]*domain.CitizenReport{
			reportID: {ID: reportID, ReferenceNo: refNo, Status: "COMPLETED", FeedbackRating: 4}, // Already rated
		},
	}
	audit := &mockAuditRepoForLoopholeAudit{}
	uc := usecase.NewReportUsecase(repo, audit, nil, &storage.MockStorageService{})

	err := uc.SubmitFeedback(context.Background(), refNo, 5, "Rating again")
	if err == nil {
		t.Fatalf("Expected error when attempting to submit feedback twice, got nil")
	}
}

func TestLoophole_CancelReport_AuditLogged(t *testing.T) {
	refNo := "DBD-DPK-CANCEL-01"
	reportID := uuid.New()

	repo := &mockReportRepoForLoopholeAudit{
		reports: map[uuid.UUID]*domain.CitizenReport{
			reportID: {ID: reportID, ReferenceNo: refNo, Status: "MENUNGGU_PE_KADER"},
		},
	}
	audit := &mockAuditRepoForLoopholeAudit{}
	uc := usecase.NewReportUsecase(repo, audit, nil, &storage.MockStorageService{})

	err := uc.CancelReport(context.Background(), refNo)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	foundCancelAudit := false
	for _, act := range audit.actions {
		if act == "CANCEL_REPORT_BY_USER" {
			foundCancelAudit = true
			break
		}
	}
	if !foundCancelAudit {
		t.Fatalf("Expected CANCEL_REPORT_BY_USER audit trail, but none was logged")
	}
}

func TestLoophole_VerifyReport_CrossJurisdictionBlocked(t *testing.T) {
	reportID := uuid.New()
	kaderID := uuid.New()
	kaderLocID := uuid.New()
	differentLocID := uuid.New()

	repo := &mockReportRepoForLoopholeAudit{
		reports: map[uuid.UUID]*domain.CitizenReport{
			reportID: {ID: reportID, Status: "MENUNGGU_PE_KADER"},
		},
		locations: map[uuid.UUID]uuid.UUID{
			reportID: differentLocID, // Report is in a different jurisdiction
		},
	}
	audit := &mockAuditRepoForLoopholeAudit{}
	verifRepo := &mockVerificationRepo{}
	verifUc := usecase.NewVerificationUsecase(verifRepo, repo, audit, nil, nil)

	// Kader from another area attempting to reject report MUST be blocked
	err := verifUc.VerifyReport(context.Background(), reportID, kaderID, domain.RoleKader, &kaderLocID, "REJECTED", "Hoax")
	if err == nil {
		t.Fatalf("Expected jurisdiction error when Kader rejects report outside their territory, got nil")
	}
}

func TestLoophole_FoggingSchedule_CrossTeamExecutionBlocked(t *testing.T) {
	scheduleID := uuid.New()
	assignedTeamID := uuid.New()
	unauthorizedTeamID := uuid.New()

	schedule := &domain.FoggingSchedule{
		ID:              scheduleID,
		ExecutionStatus: "SCHEDULED",
		TeamID:          assignedTeamID,
	}

	schedRepo := &mockScheduleRepo{schedule: schedule}
	uc := usecase.NewFoggingScheduleUsecase(schedRepo, nil, nil, nil, nil, nil, nil, nil)

	err := uc.UpdateExecutionStatus(context.Background(), scheduleID, unauthorizedTeamID, "COMPLETED", "Unauthorized update", []usecase.UsedItemInput{{ItemID: uuid.New(), Quantity: 5.0}})
	if err == nil {
		t.Fatalf("Expected unauthorized error when non-assigned team tries to update schedule, got nil")
	}
}

func TestLoophole_Inventory_NonPositiveUsageBlocked(t *testing.T) {
	uc := usecase.NewInventoryUsecase(nil, nil)

	err := uc.LogFoggingUsage(context.Background(), uuid.New(), uuid.New(), uuid.New(), 0.0, "Zero quantity")
	if err == nil {
		t.Fatalf("Expected error for zero usage quantity, got nil")
	}

	errNeg := uc.LogFoggingUsage(context.Background(), uuid.New(), uuid.New(), uuid.New(), -5.0, "Negative quantity")
	if errNeg == nil {
		t.Fatalf("Expected error for negative usage quantity, got nil")
	}
}
