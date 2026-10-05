package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockReportRepoForDuplicateMerge struct {
	reports      map[uuid.UUID]*domain.CitizenReport
	mergeCallCnt int
	auditLogged  bool
}

func (m *mockReportRepoForDuplicateMerge) CheckDuplicateRadius(ctx context.Context, lat, lng float64, radius float64, days int) (bool, error) {
	return false, nil
}
func (m *mockReportRepoForDuplicateMerge) CheckSpamClustering(ctx context.Context, lat, lng float64, radius float64, hours int) (int, error) {
	return 0, nil
}
func (m *mockReportRepoForDuplicateMerge) CreateReportTransaction(ctx context.Context, reporter *domain.Reporter, report *domain.CitizenReport, location *domain.ReportLocation, photos []domain.ReportPhoto) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) GetReportByReference(ctx context.Context, refNo string) (*domain.CitizenReport, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	if r, ok := m.reports[id]; ok {
		return r, nil
	}
	return nil, errors.New("not found")
}
func (m *mockReportRepoForDuplicateMerge) GetReports(ctx context.Context, locationID *uuid.UUID) ([]domain.CitizenReport, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) UpdateReportStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) SubmitAppealAtomic(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) SubmitFeedback(ctx context.Context, refNo string, rating int, notes string) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) GetOverdueReports(ctx context.Context, cutoffTime time.Time) ([]domain.CitizenReport, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) GetUnverifiedOverdueReports(ctx context.Context, cutoffTime time.Time) ([]domain.CitizenReport, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) EscalateReport(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) GetReportsByPhone(ctx context.Context, phone string) ([]domain.CitizenReport, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) CountCasesByArea(ctx context.Context, areaName string, areaType string, days int) (int, error) {
	return 0, nil
}
func (m *mockReportRepoForDuplicateMerge) GetRecentCasesInSameKelurahanDifferentRW(ctx context.Context, kelurahan string, excludeRW string, days int) ([]domain.CitizenReport, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) AssignInvestigator(ctx context.Context, id uuid.UUID, investigatorID uuid.UUID) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) ClaimReportAtomic(ctx context.Context, id uuid.UUID, kaderID uuid.UUID) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) UnassignInvestigator(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) GetExpiredInvestigatingReports(ctx context.Context, duration time.Duration) ([]domain.CitizenReport, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) GetSLAMetrics(ctx context.Context) (map[string]interface{}, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) GetCompletedReportsForFeedback(ctx context.Context, days int) ([]domain.CitizenReport, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) MarkFeedbackRequested(ctx context.Context, id uuid.UUID) error {
	return nil
}
func (m *mockReportRepoForDuplicateMerge) GetPhotoByStorageKey(ctx context.Context, storageKey string) (*domain.ReportPhoto, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) GetPEResultByPhotoURL(ctx context.Context, photoURL string) (*domain.PEResult, error) {
	return nil, nil
}
func (m *mockReportRepoForDuplicateMerge) VerifyReportLocation(ctx context.Context, reportID uuid.UUID, locationID uuid.UUID) (bool, error) {
	return true, nil
}
func (m *mockReportRepoForDuplicateMerge) MergeDuplicateReportsAtomic(ctx context.Context, childID uuid.UUID, primaryID uuid.UUID, actorID uuid.UUID) error {
	if m.reports[childID].Status != "MENUNGGU_PE_KADER" {
		return errors.New("status tidak valid")
	}
	m.reports[childID].Status = "MERGED_DUPLICATE"
	m.reports[childID].MergedWithID = &primaryID
	m.reports[childID].InvestigatorID = &actorID
	m.mergeCallCnt++
	return nil
}

type mockAuditRepoForDuplicate struct {
	logged bool
}

func (m *mockAuditRepoForDuplicate) LogAction(ctx context.Context, tx *gorm.DB, actorID uuid.UUID, action string, entity string, entityID uuid.UUID, details interface{}) error {
	if action == "MERGE_DUPLICATE" {
		m.logged = true
	}
	return nil
}
func (m *mockAuditRepoForDuplicate) GetAuditLogs(ctx context.Context, limit int, offset int) ([]domain.AuditLog, error) {
	return nil, nil
}
func (m *mockAuditRepoForDuplicate) GetRecentLogs(ctx context.Context, limit int) ([]domain.AuditLog, error) {
	return nil, nil
}
func (m *mockAuditRepoForDuplicate) LogStatusChange(ctx context.Context, tx *gorm.DB, reportID uuid.UUID, from, to, reason string, user uuid.UUID) error {
	return nil
}


func TestMergeDuplicateLifecycle(t *testing.T) {
	primaryID := uuid.New()
	childID := uuid.New()
	kaderID := uuid.New()

	repo := &mockReportRepoForDuplicateMerge{
		reports: map[uuid.UUID]*domain.CitizenReport{
			primaryID: {ID: primaryID, Status: "MENUNGGU_PE_KADER"},
			childID:   {ID: childID, Status: "MENUNGGU_PE_KADER"},
		},
	}
	auditRepo := &mockAuditRepoForDuplicate{}

	uc := usecase.NewReportUsecase(repo, auditRepo, nil, &storage.MockStorageService{})

	// Valid merge
	err := uc.MergeDuplicateReport(context.Background(), childID, primaryID, kaderID, "Sama-sama genangan air di lokasi yang sama")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if repo.reports[childID].Status != "MERGED_DUPLICATE" {
		t.Errorf("Expected status MERGED_DUPLICATE, got: %s", repo.reports[childID].Status)
	}

	if repo.reports[childID].MergedWithID == nil || *repo.reports[childID].MergedWithID != primaryID {
		t.Errorf("Expected MergedWithID to be set to %s", primaryID)
	}

	if !auditRepo.logged {
		t.Errorf("Expected audit trail for MERGE_DUPLICATE to be logged")
	}

	// Try to merge already merged (Invalid status)
	err = uc.MergeDuplicateReport(context.Background(), childID, primaryID, kaderID, "Another reason")
	if err == nil {
		t.Errorf("Expected error when trying to merge report not in MENUNGGU_PE_KADER status")
	}
}
