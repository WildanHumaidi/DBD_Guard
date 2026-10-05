package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"strings"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/domain/dto"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockSpamReportRepo struct {
	repository.ReportRepository
	spamCount   int
	savedReport *domain.CitizenReport
}

func (m *mockSpamReportRepo) CheckDuplicateRadius(ctx context.Context, lat, lng float64, radius float64, days int) (bool, error) {
	return false, nil
}

func (m *mockSpamReportRepo) CheckSpamClustering(ctx context.Context, lat, lng float64, radius float64, hours int) (int, error) {
	return m.spamCount, nil
}

func (m *mockSpamReportRepo) CreateReportTransaction(ctx context.Context, reporter *domain.Reporter, report *domain.CitizenReport, location *domain.ReportLocation, photos []domain.ReportPhoto) error {
	m.savedReport = report
	return nil
}

type mockSpamNotifUsecase struct {
	usecase.NotificationUsecase
	newReportNotifCalled bool
}

func (m *mockSpamNotifUsecase) NotifyNewReport(ctx context.Context, reportID string, areaRT string, areaRW string, areaKelurahan string) error {
	m.newReportNotifCalled = true
	return nil
}

type mockSpamAuditRepo struct {
	repository.AuditRepository
	loggedAction string
}

func (m *mockSpamAuditRepo) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action, entity string, entityID uuid.UUID, details interface{}) error {
	m.loggedAction = action
	return nil
}

func TestBUGNEW012_SpamReport_NoNotificationLeak(t *testing.T) {
	ctx := context.Background()

	// Scenario 1: Spam Report (spamCount >= 3)
	t.Run("Spam report is auto-rejected without notifying Kader", func(t *testing.T) {
		repo := &mockSpamReportRepo{spamCount: 4}
		notif := &mockSpamNotifUsecase{}
		audit := &mockSpamAuditRepo{}

		u := usecase.NewReportUsecase(repo, audit, notif, &storage.MockStorageService{})

		req := dto.SubmitReportRequest{
			Name:        "Spammer",
			Contact:     "08111111111",
			Category:    "KASUS_DBD",
			Description: "Spam report test",
			Latitude:    -6.4025,
			Longitude:   106.7942,
			Kelurahan:   "Depok",
			Kecamatan:   "Pancoran Mas",
		}

		res, err := u.SubmitCitizenReport(ctx, req)
		if err != nil {
			t.Fatalf("SubmitCitizenReport failed: %v", err)
		}

		// 1. Saved status must be REJECTED
		if repo.savedReport == nil || repo.savedReport.Status != "REJECTED" {
			t.Fatalf("expected report status REJECTED, got: %v", repo.savedReport)
		}

		// 2. Notification MUST NOT be sent to Kader
		if notif.newReportNotifCalled {
			t.Fatalf("BUG-NEW-012 Violation! Notification was leaked to Kader for a spam-rejected report")
		}

		// 3. Audit action logged
		if audit.loggedAction != "AUTO_REJECT_SPAM" {
			t.Fatalf("expected audit action AUTO_REJECT_SPAM, got: %s", audit.loggedAction)
		}

		// 4. Response message must inform about spam rejection
		if !strings.Contains(res.Message, "spam") {
			t.Fatalf("expected response message to mention spam rejection, got: %s", res.Message)
		}
	})

	// Scenario 2: Legitimate Report (spamCount < 3)
	t.Run("Legitimate report is processed and notifies Kader", func(t *testing.T) {
		repo := &mockSpamReportRepo{spamCount: 0}
		notif := &mockSpamNotifUsecase{}
		audit := &mockSpamAuditRepo{}

		u := usecase.NewReportUsecase(repo, audit, notif, &storage.MockStorageService{})

		req := dto.SubmitReportRequest{
			Name:        "Legit Citizen",
			Contact:     "08122222222",
			Category:    "KASUS_DBD",
			Description: "Gejala DBD nyata",
			Latitude:    -6.4025,
			Longitude:   106.7942,
			Kelurahan:   "Depok",
			Kecamatan:   "Pancoran Mas",
		}

		res, err := u.SubmitCitizenReport(ctx, req)
		if err != nil {
			t.Fatalf("SubmitCitizenReport failed: %v", err)
		}

		if repo.savedReport.Status != "MENUNGGU_PE_KADER" {
			t.Fatalf("expected status MENUNGGU_PE_KADER, got: %s", repo.savedReport.Status)
		}

		if !notif.newReportNotifCalled {
			t.Fatalf("expected NotifyNewReport to be called for legitimate report")
		}

		if !strings.Contains(res.Message, "penyelidikan (PE) oleh Kader") {
			t.Fatalf("expected normal success message, got: %s", res.Message)
		}
	})
}
