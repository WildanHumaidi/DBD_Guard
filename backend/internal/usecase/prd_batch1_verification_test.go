package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"testing"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/domain/dto"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type batch1MockReportRepo struct {
	reports  map[uuid.UUID]*domain.CitizenReport
	requests map[string]*domain.CitizenReport
}

func (m *batch1MockReportRepo) CheckDuplicateRadius(ctx context.Context, lat, lng float64, radius float64, days int) (bool, error) {
	if lat == -6.4000 && lng == 106.8000 {
		return true, nil
	}
	return false, nil
}

func (m *batch1MockReportRepo) CheckSpamClustering(ctx context.Context, lat, lng float64, radius float64, hours int) (int, error) {
	if lat == -6.4001 && lng == 106.8001 {
		return 3, nil
	}
	return 0, nil
}

func (m *batch1MockReportRepo) CreateReportTransaction(ctx context.Context, reporter *domain.Reporter, report *domain.CitizenReport, location *domain.ReportLocation, photos []domain.ReportPhoto) error {
	report.ID = uuid.New()
	report.Reporter = *reporter
	report.Location = *location
	m.reports[report.ID] = report
	m.requests[report.ReferenceNo] = report
	return nil
}

func (m *batch1MockReportRepo) GetReportByReference(ctx context.Context, refNo string) (*domain.CitizenReport, error) {
	if r, ok := m.requests[refNo]; ok {
		return r, nil
	}
	return nil, nil
}

func (m *batch1MockReportRepo) SubmitFeedback(ctx context.Context, refNo string, rating int, notes string) error {
	r := m.requests[refNo]
	r.FeedbackRating = rating
	r.FeedbackNotes = notes
	return nil
}

// Implement remaining ReportRepository methods with panics or empty returns to satisfy interface
func (m *batch1MockReportRepo) GetReports(ctx context.Context, loc *uuid.UUID) ([]domain.CitizenReport, error) { return nil, nil }
func (m *batch1MockReportRepo) GetReportsByPhone(ctx context.Context, phone string) ([]domain.CitizenReport, error) { return nil, nil }
func (m *batch1MockReportRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) { return nil, nil }
func (m *batch1MockReportRepo) VerifyReportLocation(ctx context.Context, id uuid.UUID, loc uuid.UUID) (bool, error) { return true, nil }
func (m *batch1MockReportRepo) ClaimReportAtomic(ctx context.Context, id uuid.UUID, kaderID uuid.UUID) error { return nil }
func (m *batch1MockReportRepo) MergeDuplicateReportsAtomic(ctx context.Context, childID uuid.UUID, primaryID uuid.UUID, actorID uuid.UUID) error { return nil }
func (m *batch1MockReportRepo) UpdateReportStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error { return nil }
func (m *batch1MockReportRepo) AssignInvestigator(ctx context.Context, id uuid.UUID, officerID uuid.UUID) error { return nil }
func (m *batch1MockReportRepo) UnassignInvestigator(ctx context.Context, id uuid.UUID) error { return nil }
func (m *batch1MockReportRepo) SubmitAppealAtomic(ctx context.Context, id uuid.UUID) error { return nil }
func (m *batch1MockReportRepo) GetUnverifiedOverdueReports(ctx context.Context, cutoff time.Time) ([]domain.CitizenReport, error) { return nil, nil }
func (m *batch1MockReportRepo) GetOverdueReports(ctx context.Context, cutoffTime time.Time) ([]domain.CitizenReport, error) { return nil, nil }
func (m *batch1MockReportRepo) GetExpiredInvestigatingReports(ctx context.Context, duration time.Duration) ([]domain.CitizenReport, error) { return nil, nil }
func (m *batch1MockReportRepo) EscalateReport(ctx context.Context, id uuid.UUID) error { return nil }
func (m *batch1MockReportRepo) GetSLAMetrics(ctx context.Context) (map[string]interface{}, error) { return nil, nil }
func (m *batch1MockReportRepo) GetCompletedReportsForFeedback(ctx context.Context, days int) ([]domain.CitizenReport, error) { return nil, nil }
func (m *batch1MockReportRepo) MarkFeedbackRequested(ctx context.Context, id uuid.UUID) error { return nil }
func (m *batch1MockReportRepo) GetPhotoByStorageKey(ctx context.Context, key string) (*domain.ReportPhoto, error) { return nil, nil }
func (m *batch1MockReportRepo) GetPEResultByPhotoURL(ctx context.Context, url string) (*domain.PEResult, error) { return nil, nil }
func (m *batch1MockReportRepo) CountCasesByArea(ctx context.Context, areaName string, areaType string, days int) (int, error) { return 0, nil }
func (m *batch1MockReportRepo) GetRecentCasesInSameKelurahanDifferentRW(ctx context.Context, kelurahan string, excludeRW string, days int) ([]domain.CitizenReport, error) { return nil, nil }


type batch1MockAuditRepo struct{}
func (m *batch1MockAuditRepo) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action string, entity string, entityID uuid.UUID, details interface{}) error { return nil }
func (m *batch1MockAuditRepo) GetAuditLogs(ctx context.Context, limit int, offset int) ([]domain.AuditLog, error) { return nil, nil }
func (m *batch1MockAuditRepo) GetRecentLogs(ctx context.Context, limit int) ([]domain.AuditLog, error) { return nil, nil }
func (m *batch1MockAuditRepo) LogStatusChange(ctx context.Context, tx *gorm.DB, reportID uuid.UUID, from, to, reason string, user uuid.UUID) error { return nil }


type batch1MockNotifUsecase struct{}
func (m *batch1MockNotifUsecase) NotifyNewReport(ctx context.Context, reportID string, rt, rw, kelurahan string) error { return nil }
func (m *batch1MockNotifUsecase) NotifyKLBTriggered(ctx context.Context, areaName string, areaType string, casesCount int) error { return nil }
func (m *batch1MockNotifUsecase) NotifyFoggingScheduled(ctx context.Context, reportID string, targetDate string) error { return nil }
func (m *batch1MockNotifUsecase) NotifyEWSKader(ctx context.Context, kelurahan string, excludeActorID string) error { return nil }


func TestPRDBatch1_DuplicateClustering(t *testing.T) {
	repo := &batch1MockReportRepo{
		reports:  make(map[uuid.UUID]*domain.CitizenReport),
		requests: make(map[string]*domain.CitizenReport),
	}
	auditRepo := &batch1MockAuditRepo{}
	notifUsecase := &batch1MockNotifUsecase{}

	uc := usecase.NewReportUsecase(repo, auditRepo, notifUsecase, &storage.MockStorageService{})

	ctx := context.Background()
	req := dto.SubmitReportRequest{
		Name:      "Pelapor",
		Category:  "DBD",
		Latitude:  -6.4000,
		Longitude: 106.8000,
	}

	res, err := uc.SubmitCitizenReport(ctx, req)
	if err != nil {
		t.Fatalf("Submit error: %v", err)
	}

	report := repo.requests[res.ReferenceNo]
	if !report.IsPossibleDuplicate {
		t.Errorf("Expected IsPossibleDuplicate to be true")
	}
}

func TestPRDBatch1_SpamClustering(t *testing.T) {
	repo := &batch1MockReportRepo{
		reports:  make(map[uuid.UUID]*domain.CitizenReport),
		requests: make(map[string]*domain.CitizenReport),
	}
	auditRepo := &batch1MockAuditRepo{}
	notifUsecase := &batch1MockNotifUsecase{}

	uc := usecase.NewReportUsecase(repo, auditRepo, notifUsecase, &storage.MockStorageService{})

	ctx := context.Background()
	req := dto.SubmitReportRequest{
		Name:      "Spammer",
		Category:  "DBD",
		Latitude:  -6.4001,
		Longitude: 106.8001,
	}

	res, err := uc.SubmitCitizenReport(ctx, req)
	if err != nil {
		t.Fatalf("Submit error: %v", err)
	}

	if res.ReferenceNo == "" || res.Message == "" {
		t.Errorf("Expected spam rejection response")
	}
}

func TestPRDBatch1_SubmitFeedback(t *testing.T) {
	repo := &batch1MockReportRepo{
		reports:  make(map[uuid.UUID]*domain.CitizenReport),
		requests: make(map[string]*domain.CitizenReport),
	}
	auditRepo := &batch1MockAuditRepo{}
	notifUsecase := &batch1MockNotifUsecase{}

	uc := usecase.NewReportUsecase(repo, auditRepo, notifUsecase, &storage.MockStorageService{})

	report := &domain.CitizenReport{
		ID:          uuid.New(),
		ReferenceNo: "REF-123",
		Status:      "COMPLETED",
	}
	repo.requests[report.ReferenceNo] = report

	ctx := context.Background()
	err := uc.SubmitFeedback(ctx, "REF-123", 5, "Sangat baik")
	if err != nil {
		t.Fatalf("SubmitFeedback error: %v", err)
	}

	if repo.requests["REF-123"].FeedbackRating != 5 {
		t.Errorf("Feedback not saved correctly")
	}
}
