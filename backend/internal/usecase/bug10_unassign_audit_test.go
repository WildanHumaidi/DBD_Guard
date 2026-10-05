package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockForceUnassignReportRepo struct {
	repository.ReportRepository
	report   *domain.CitizenReport
	unassigned bool
}

func (m *mockForceUnassignReportRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	return m.report, nil
}

func (m *mockForceUnassignReportRepo) UnassignInvestigator(ctx context.Context, id uuid.UUID) error {
	m.unassigned = true
	m.report.InvestigatorID = nil
	return nil
}

type mockForceUnassignAuditRepo struct {
	repository.AuditRepository
	action  string
	details map[string]interface{}
}

func (m *mockForceUnassignAuditRepo) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action, entity string, entityID uuid.UUID, details interface{}) error {
	m.action = action
	if d, ok := details.(map[string]interface{}); ok {
		m.details = d
	}
	return nil
}

func TestBUGNEW010_ForceUnassignReport_AuditTrail(t *testing.T) {
	reportID := uuid.New()
	prevInvestigatorID := uuid.New()

	repo := &mockForceUnassignReportRepo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: &prevInvestigatorID,
		},
	}
	audit := &mockForceUnassignAuditRepo{}

	u := usecase.NewReportUsecase(repo, audit, nil, &storage.MockStorageService{})

	err := u.ForceUnassignReport(context.Background(), reportID)
	if err != nil {
		t.Fatalf("ForceUnassignReport failed: %v", err)
	}

	if !repo.unassigned {
		t.Fatal("expected investigator to be unassigned")
	}

	if audit.action != "FORCE_UNASSIGN_REPORT" {
		t.Fatalf("BUG-NEW-010 Detected! Expected audit action FORCE_UNASSIGN_REPORT, got: %s", audit.action)
	}

	if audit.details["previous_investigator_id"] != prevInvestigatorID {
		t.Fatalf("expected previous_investigator_id %v, got: %v", prevInvestigatorID, audit.details["previous_investigator_id"])
	}
}
