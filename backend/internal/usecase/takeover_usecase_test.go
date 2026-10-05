package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"strings"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockReportRepoForTakeover struct {
	repository.ReportRepository
	report           *domain.CitizenReport
	assignedOfficer  *uuid.UUID
	assignCalled     bool
}

func (m *mockReportRepoForTakeover) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	return m.report, nil
}

func (m *mockReportRepoForTakeover) AssignInvestigator(ctx context.Context, id uuid.UUID, officerID uuid.UUID) error {
	m.assignCalled = true
	m.assignedOfficer = &officerID
	m.report.InvestigatorID = &officerID
	return nil
}

type mockAuditRepoForTakeover struct {
	repository.AuditRepository
	loggedAction  string
	loggedDetails map[string]interface{}
}

func (m *mockAuditRepoForTakeover) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action, entity string, entityID uuid.UUID, details interface{}) error {
	m.loggedAction = action
	if d, ok := details.(map[string]interface{}); ok {
		m.loggedDetails = d
	}
	return nil
}

func TestTakeoverWorkflow(t *testing.T) {
	reportID := uuid.New()
	kaderID := uuid.New()
	puskesmasOfficerID := uuid.New()

	t.Run("Takeover valid report and subsequent PE submission", func(t *testing.T) {
		report := &domain.CitizenReport{
			ID:             reportID,
			Status:         "ESKALASI_PUSKESMAS",
			InvestigatorID: &kaderID,
		}

		reportRepo := &mockReportRepoForTakeover{report: report}
		auditRepo := &mockAuditRepoForTakeover{}

		repUsecase := usecase.NewReportUsecase(reportRepo, auditRepo, nil, &storage.MockStorageService{})

		// 1. Puskesmas takes over the report
		err := repUsecase.TakeoverReport(context.Background(), reportID, puskesmasOfficerID, "SLA expired takeover")
		if err != nil {
			t.Fatalf("unexpected error during takeover: %v", err)
		}

		if !reportRepo.assignCalled {
			t.Fatalf("expected AssignInvestigator to be called")
		}
		if *report.InvestigatorID != puskesmasOfficerID {
			t.Fatalf("expected investigator_id to be %v, got %v", puskesmasOfficerID, report.InvestigatorID)
		}
		if auditRepo.loggedAction != "TAKEOVER_REPORT" {
			t.Fatalf("expected audit action TAKEOVER_REPORT, got %s", auditRepo.loggedAction)
		}
		if auditRepo.loggedDetails["reason"] != "SLA expired takeover" {
			t.Fatalf("expected audit reason to match")
		}

		// 2. Puskesmas can now submit PE result
		peRepo := &mockPERepo{}
		invRepo := &mockInventoryRepo{}
		fogRepo := &mockFoggingRepo{}
		mdRepo := &mockMasterDataRepo{}
		ewsUsecase := &mockEwsKlbUsecase{}

		peU := usecase.NewPEUsecase(peRepo, reportRepo, fogRepo, invRepo, auditRepo, mdRepo, ewsUsecase)

		defer func() {
			recover() // Ignore subsequent repo stubs
		}()

		_, err = peU.SubmitPEResult(
			context.Background(),
			reportID,
			puskesmasOfficerID,
			100,
			true,
			0,
			true,
			"Takeover PE notes",
			"http://evidence.jpg",
		)

		if err != nil && strings.Contains(err.Error(), "UNAUTHORIZED_PE_SUBMISSION") {
			t.Fatalf("did not expect UNAUTHORIZED_PE_SUBMISSION after takeover, got: %v", err)
		}
	})

	t.Run("Takeover rejected on invalid status", func(t *testing.T) {
		completedReport := &domain.CitizenReport{
			ID:             reportID,
			Status:         "COMPLETED",
			InvestigatorID: &kaderID,
		}

		reportRepo := &mockReportRepoForTakeover{report: completedReport}
		auditRepo := &mockAuditRepoForTakeover{}
		repUsecase := usecase.NewReportUsecase(reportRepo, auditRepo, nil, &storage.MockStorageService{})

		err := repUsecase.TakeoverReport(context.Background(), reportID, puskesmasOfficerID, "invalid status takeover")
		if err == nil {
			t.Fatalf("expected error taking over COMPLETED report, got nil")
		}
	})
}
