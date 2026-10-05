package usecase_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	deliveryHttp "github.com/dbd-guard/backend/internal/delivery/http"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/domain/dto"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type mockAnonReportRepo struct {
	repository.ReportRepository
	savedReporter       *domain.Reporter
	savedReport         *domain.CitizenReport
	reportByRef         *domain.CitizenReport
	reportByID          *domain.CitizenReport
	completedReports    []domain.CitizenReport
	feedbackMarkedCount int
}

func (m *mockAnonReportRepo) GetCompletedReportsForFeedback(ctx context.Context, days int) ([]domain.CitizenReport, error) {
	return m.completedReports, nil
}

func (m *mockAnonReportRepo) MarkFeedbackRequested(ctx context.Context, id uuid.UUID) error {
	m.feedbackMarkedCount++
	return nil
}

func (m *mockAnonReportRepo) CheckDuplicateRadius(ctx context.Context, lat, lng float64, radius float64, days int) (bool, error) {
	return false, nil
}

func (m *mockAnonReportRepo) CheckSpamClustering(ctx context.Context, lat, lng float64, radius float64, hours int) (int, error) {
	return 0, nil
}

func (m *mockAnonReportRepo) CreateReportTransaction(ctx context.Context, reporter *domain.Reporter, report *domain.CitizenReport, location *domain.ReportLocation, photos []domain.ReportPhoto) error {
	m.savedReporter = reporter
	m.savedReport = report
	return nil
}

func (m *mockAnonReportRepo) GetReportByReference(ctx context.Context, refNo string) (*domain.CitizenReport, error) {
	return m.reportByRef, nil
}

func (m *mockAnonReportRepo) UpdateReportStatus(ctx context.Context, tx *gorm.DB, reportID uuid.UUID, status string) error {
	if m.reportByID != nil && m.reportByID.ID == reportID {
		m.reportByID.Status = status
	}
	return nil
}

func (m *mockAnonReportRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	if m.reportByID != nil && m.reportByID.ID == id {
		return m.reportByID, nil
	}
	return nil, nil
}

type mockAnonNotifUsecase struct {
	usecase.NotificationUsecase
	foggingNotifSent bool
}

func (m *mockAnonNotifUsecase) NotifyFoggingScheduled(ctx context.Context, reportID string, targetDate string) error {
	m.foggingNotifSent = true
	return nil
}

func (m *mockAnonNotifUsecase) NotifyNewReport(ctx context.Context, reportID string, areaRT string, areaRW string, areaKelurahan string) error {
	return nil
}

func TestAnonymousReport_EndToEndPIIProtection(t *testing.T) {
	repo := &mockAnonReportRepo{}
	audit := &mockAudRepo{}
	notif := &mockAnonNotifUsecase{}

	u := usecase.NewReportUsecase(repo, audit, notif, &storage.MockStorageService{})
	ctx := context.Background()

	// 1. Submit Anonymous Report with PII in request (user typed name & phone, but checked anonymous)
	req := dto.SubmitReportRequest{
		Name:          "John Doe Secret",
		Contact:       "081299998888",
		AnonymousFlag: true,
		Category:      "JENTIK_NYAMUK",
		Description:   "Ada genangan air di pot tetangga",
		Latitude:      -6.4,
		Longitude:     106.8,
		Kelurahan:     "Sukatani",
		Kecamatan:     "Tapos",
	}

	res, err := u.SubmitCitizenReport(ctx, req)
	if err != nil {
		t.Fatalf("SubmitCitizenReport failed: %v", err)
	}

	// Verify tracking reference is generated
	if res.ReferenceNo == "" {
		t.Fatal("expected ReferenceNo to be generated for anonymous report")
	}

	// Verify reporter entity has NO PII saved
	if repo.savedReporter == nil {
		t.Fatal("expected reporter to be saved")
	}
	if repo.savedReporter.Name != "Anonymous" {
		t.Errorf("expected reporter name 'Anonymous', got: %s", repo.savedReporter.Name)
	}
	if repo.savedReporter.Contact != "" {
		t.Errorf("expected empty contact for anonymous reporter, got: %s", repo.savedReporter.Contact)
	}
	if !repo.savedReporter.AnonymousFlag {
		t.Errorf("expected AnonymousFlag to be true")
	}

	// 2. Test tracking via Reference Number
	repo.reportByRef = &domain.CitizenReport{
		ID:          uuid.New(),
		ReferenceNo: res.ReferenceNo,
		Status:      "MENUNGGU_PE_KADER",
		Reporter: domain.Reporter{
			Name:          "Anonymous",
			Contact:       "",
			AnonymousFlag: true,
		},
	}

	tracked, err := u.TrackReport(ctx, res.ReferenceNo)
	if err != nil {
		t.Fatalf("TrackReport failed: %v", err)
	}
	if tracked.Reporter.Name != "Anonymous" || tracked.Reporter.Contact != "" {
		t.Errorf("PII leaked in tracked report: Name=%s, Contact=%s", tracked.Reporter.Name, tracked.Reporter.Contact)
	}
}

type mockFogSchedRepoForAnon struct {
	repository.FoggingScheduleRepository
}

func (m *mockFogSchedRepoForAnon) CreateSchedule(ctx context.Context, schedule *domain.FoggingSchedule) error {
	schedule.ID = uuid.New()
	return nil
}

func (m *mockFogSchedRepoForAnon) CreateScheduleWithTx(ctx context.Context, tx *gorm.DB, schedule *domain.FoggingSchedule) error {
	schedule.ID = uuid.New()
	return nil
}

func (m *mockFogSchedRepoForAnon) CountSchedulesByTeamAndDate(ctx context.Context, teamID uuid.UUID, date time.Time) (int64, error) {
	return 0, nil
}

func (m *mockFogSchedRepoForAnon) GetActiveScheduleByRequestID(ctx context.Context, requestID uuid.UUID) (*domain.FoggingSchedule, error) {
	return nil, nil
}

func (m *mockFogSchedRepoForAnon) WithTransaction(fn func(tx *gorm.DB) error) error {
	return fn(nil)
}

type mockApprovedFogRepo struct {
	repository.FoggingRepository
	reportID uuid.UUID
}

func (m *mockApprovedFogRepo) GetFoggingRequestByID(ctx context.Context, id uuid.UUID) (*domain.FoggingRequest, error) {
	peID := uuid.New()
	return &domain.FoggingRequest{
		ID:             id,
		ApprovalStatus: "APPROVED",
		PEID:           peID,
		PE: domain.PEResult{
			ID:       peID,
			ReportID: m.reportID,
		},
	}, nil
}

func TestAnonymousReport_NoNotificationSentOnFoggingSchedule(t *testing.T) {
	anonReportID := uuid.New()
	reqID := uuid.New()

	schedRepo := &mockFogSchedRepoForAnon{}
	fogRepo := &mockApprovedFogRepo{reportID: anonReportID}
	invRepo := &mockInvRepoForFogging{}
	invUsecase := &mockInvUsecase{}
	audRepo := &mockAudRepo{}
	mdRepo := &mockMDRepo{}
	notif := &mockAnonNotifUsecase{}

	// Report repo returns anonymous reporter
	repRepo := &mockAnonReportRepo{
		reportByID: &domain.CitizenReport{
			ID: anonReportID,
			Reporter: domain.Reporter{
				Name:          "Anonymous",
				Contact:       "",
				AnonymousFlag: true,
			},
		},
	}

	u := usecase.NewFoggingScheduleUsecase(schedRepo, fogRepo, invRepo, invUsecase, audRepo, repRepo, mdRepo, notif)

	// Create schedule for anonymous report
	teamID := uuid.New()
	creatorID := uuid.New()
	_, err := u.CreateSchedule(context.Background(), reqID, teamID, creatorID, time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("CreateSchedule failed: %v", err)
	}

	// Verify citizen push notification was NOT triggered for anonymous reporter
	if notif.foggingNotifSent {
		t.Fatalf("BUG-10 violation: push notification sent to anonymous reporter")
	}
}

func TestAnonymousReport_DirectAPISubmission_EndToEndInvariants(t *testing.T) {
	repo := &mockAnonReportRepo{}
	audit := &mockAudRepo{}
	notif := &mockAnonNotifUsecase{}

	reportUsecase := usecase.NewReportUsecase(repo, audit, notif, &storage.MockStorageService{})
	reportHandler := deliveryHttp.NewReportHandler(reportUsecase)

	app := fiber.New()
	app.Post("/reports", reportHandler.SubmitCitizenReport)

	// Direct API submission with anonymous_flag = true and explicit PII in payload
	payload := `{
		"name": "Jane Doe Secret Citizen",
		"contact": "081987654321",
		"anonymous_flag": true,
		"category": "KASUS_DBD",
		"description": "Gejala demam berdarah dan ditemukan jentik",
		"latitude": -6.4025,
		"longitude": 106.7942,
		"address": "Jl. Margonda Raya No. 100",
		"rt": "02",
		"rw": "05",
		"kelurahan": "Depok",
		"kecamatan": "Pancoran Mas"
	}`

	req := httptest.NewRequest("POST", "/reports", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("API request failed: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Fatalf("expected HTTP 201 Created, got: %d", resp.StatusCode)
	}

	var respBody struct {
		Status string `json:"status"`
		Data   struct {
			ReferenceNo string `json:"reference_no"`
			Message     string `json:"message"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	refNo := respBody.Data.ReferenceNo
	if refNo == "" {
		t.Fatal("expected reference_no in response body")
	}

	// 1. Invariant: PII is NOT stored in database
	if repo.savedReporter == nil {
		t.Fatal("expected reporter to be saved")
	}
	if repo.savedReporter.Name != "Anonymous" {
		t.Errorf("PII Leak! Expected saved name 'Anonymous', got: %s", repo.savedReporter.Name)
	}
	if repo.savedReporter.Contact != "" {
		t.Errorf("PII Leak! Expected empty contact, got: %s", repo.savedReporter.Contact)
	}
	if !repo.savedReporter.AnonymousFlag {
		t.Error("expected AnonymousFlag to be true")
	}

	// 2. Invariant: Reference tracking still works correctly
	repo.reportByRef = &domain.CitizenReport{
		ID:          uuid.New(),
		ReferenceNo: refNo,
		Status:      "MENUNGGU_PE_KADER",
		Reporter: domain.Reporter{
			Name:          "Anonymous",
			Contact:       "",
			AnonymousFlag: true,
		},
	}
	tracked, err := reportUsecase.TrackReport(context.Background(), refNo)
	if err != nil {
		t.Fatalf("TrackReport failed for anonymous report: %v", err)
	}
	if tracked.Reporter.Name != "Anonymous" || tracked.Reporter.Contact != "" {
		t.Errorf("PII Leak in Tracking: %s, %s", tracked.Reporter.Name, tracked.Reporter.Contact)
	}

	// 3. Invariant: Feedback reminder is NOT sent to anonymous report
	repo.completedReports = []domain.CitizenReport{
		{
			ID:          uuid.New(),
			ReferenceNo: refNo,
			Status:      "COMPLETED",
			Reporter: domain.Reporter{
				Name:          "Anonymous",
				Contact:       "",
				AnonymousFlag: true,
			},
			Location: domain.ReportLocation{
				Address: "Jl. Margonda Raya No. 100",
			},
		},
	}
	reminderCount, err := reportUsecase.SendFeedbackReminders(context.Background())
	if err != nil {
		t.Fatalf("SendFeedbackReminders failed: %v", err)
	}
	if reminderCount != 0 {
		t.Errorf("BUG-10 violation! Reminder sent to anonymous report: count=%d", reminderCount)
	}
	if repo.feedbackMarkedCount != 0 {
		t.Errorf("BUG-10 violation! Feedback marked for anonymous report: count=%d", repo.feedbackMarkedCount)
	}
}

