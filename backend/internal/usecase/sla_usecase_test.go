package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"testing"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockSLAReportRepo struct {
	repository.ReportRepository
	reports []domain.CitizenReport
}

func (m *mockSLAReportRepo) GetUnverifiedOverdueReports(ctx context.Context, cutoffTime time.Time) ([]domain.CitizenReport, error) {
	var overdue []domain.CitizenReport
	for _, r := range m.reports {
		if r.Status == "MENUNGGU_PE_KADER" && r.CreatedAt.Before(cutoffTime) {
			overdue = append(overdue, r)
		}
	}
	return overdue, nil
}

func (m *mockSLAReportRepo) EscalateReport(ctx context.Context, id uuid.UUID) error {
	for i := range m.reports {
		if m.reports[i].ID == id {
			m.reports[i].Status = "ESKALASI_PUSKESMAS"
			m.reports[i].InvestigatorID = nil
			return nil
		}
	}
	return nil
}

type mockSLAAuditRepo struct {
	repository.AuditRepository
	escalationLogged bool
}

func (m *mockSLAAuditRepo) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action string, entity string, entityID uuid.UUID, details interface{}) error {
	if action == "SLA_ESCALATION" {
		m.escalationLogged = true
	}
	return nil
}

func TestAutoReleaseExpiredTasks_72hSLA(t *testing.T) {
	now := time.Now()
	kaderID := uuid.New()

	oldReportID := uuid.New()
	recentReportID := uuid.New()

	repo := &mockSLAReportRepo{
		reports: []domain.CitizenReport{
			{
				ID:             oldReportID,
				Status:         "MENUNGGU_PE_KADER",
				InvestigatorID: &kaderID,
				CreatedAt:      now.Add(-75 * time.Hour), // > 72 hours
			},
			{
				ID:             recentReportID,
				Status:         "MENUNGGU_PE_KADER",
				InvestigatorID: &kaderID,
				CreatedAt:      now.Add(-24 * time.Hour), // < 72 hours
			},
		},
	}
	audit := &mockSLAAuditRepo{}
	u := usecase.NewReportUsecase(repo, audit, nil, &storage.MockStorageService{})

	ctx := context.Background()
	count, err := u.AutoReleaseExpiredTasks(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if count != 1 {
		t.Fatalf("expected 1 escalated report, got: %d", count)
	}

	if !audit.escalationLogged {
		t.Fatalf("expected SLA_ESCALATION audit action to be logged")
	}

	// Verify old report is escalated and unassigned
	if repo.reports[0].Status != "ESKALASI_PUSKESMAS" {
		t.Fatalf("expected status ESKALASI_PUSKESMAS, got: %s", repo.reports[0].Status)
	}
	if repo.reports[0].InvestigatorID != nil {
		t.Fatalf("expected investigatorID to be nil after escalation")
	}

	// Verify recent report is unchanged
	if repo.reports[1].Status != "MENUNGGU_PE_KADER" {
		t.Fatalf("expected recent report to remain MENUNGGU_PE_KADER, got: %s", repo.reports[1].Status)
	}

	// Verify idempotency: running a second time should escalate 0 reports
	count2, err := u.AutoReleaseExpiredTasks(ctx)
	if err != nil || count2 != 0 {
		t.Fatalf("expected 0 reports escalated on second run (idempotent), got: %d, err: %v", count2, err)
	}
}
