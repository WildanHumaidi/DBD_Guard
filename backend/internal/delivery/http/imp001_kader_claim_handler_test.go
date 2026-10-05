package http_test

import (
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/dbd-guard/backend/internal/config"
	deliveryHttp "github.com/dbd-guard/backend/internal/delivery/http"
	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const testJWTSecret = "test-secret-for-imp001-verification"

func init() {
	if config.Env == nil {
		config.Env = &config.AppConfig{
			JWTSecret: testJWTSecret,
		}
	} else {
		config.Env.JWTSecret = testJWTSecret
	}
}

// mockHandlerReportRepo provides mock report repository for HTTP handler tests
type mockHandlerReportRepo struct {
	repository.ReportRepository
	mu          sync.Mutex
	report      *domain.CitizenReport
	reportLocID uuid.UUID
}

func (m *mockHandlerReportRepo) GetReportByID(ctx context.Context, id uuid.UUID) (*domain.CitizenReport, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.report == nil || m.report.ID != id {
		return nil, errors.New("record not found")
	}
	cp := *m.report
	return &cp, nil
}

func (m *mockHandlerReportRepo) VerifyReportLocation(ctx context.Context, reportID uuid.UUID, locationID uuid.UUID) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.report == nil || m.report.ID != reportID {
		return false, nil
	}
	return m.reportLocID == locationID, nil
}

func (m *mockHandlerReportRepo) ClaimReportAtomic(ctx context.Context, id uuid.UUID, kaderID uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.report == nil || m.report.ID != id {
		return errors.New("laporan tidak ditemukan")
	}
	if m.report.Status != "MENUNGGU_PE_KADER" || m.report.InvestigatorID != nil {
		return errors.New("laporan sudah diklaim oleh kader lain atau status tidak valid")
	}
	m.report.InvestigatorID = &kaderID
	return nil
}

func generateTestToken(userID uuid.UUID, role string, locationID *uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(),
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 1).Unix(),
	}
	if locationID != nil {
		claims["location_id"] = locationID.String()
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(testJWTSecret))
}

func setupTestApp(repo repository.ReportRepository) *fiber.App {
	app := fiber.New()
	reportUsecase := usecase.NewReportUsecase(repo, nil, nil, &storage.MockStorageService{})
	reportHandler := deliveryHttp.NewReportHandler(reportUsecase)

	api := app.Group("/api/v1")
	protected := api.Group("/", middleware.Protected())
	kader := protected.Group("/kader", middleware.RequireRole(domain.RoleKader))
	kader.Post("/reports/:id/claim", reportHandler.ClaimReport)

	return app
}

// TEST-01: HTTP Same Region Claim -> HTTP 200 OK
func TestIMP001_HTTP_SameRegion_Success(t *testing.T) {
	reportID := uuid.New()
	regionBeji := uuid.New()
	kaderID := uuid.New()

	repo := &mockHandlerReportRepo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionBeji,
	}

	app := setupTestApp(repo)

	token, err := generateTestToken(kaderID, domain.RoleKader, &regionBeji)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/kader/reports/"+reportID.String()+"/claim", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200 OK, got: %d", resp.StatusCode)
	}

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["message"] != "Berhasil mengklaim tugas laporan" {
		t.Fatalf("unexpected message: %v", body["message"])
	}

	// Verify DB state
	if repo.report.InvestigatorID == nil || *repo.report.InvestigatorID != kaderID {
		t.Fatalf("expected investigator_id %s, got %v", kaderID, repo.report.InvestigatorID)
	}
}

// TEST-02: HTTP Cross Region Claim -> HTTP 403 Forbidden
func TestIMP001_HTTP_CrossRegion_Forbidden(t *testing.T) {
	reportID := uuid.New()
	regionPancoranMas := uuid.New()
	regionBeji := uuid.New()
	kaderBejiID := uuid.New()

	repo := &mockHandlerReportRepo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionPancoranMas, // Report in Pancoran Mas
	}

	app := setupTestApp(repo)

	// Kader Beji
	token, err := generateTestToken(kaderBejiID, domain.RoleKader, &regionBeji)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/kader/reports/"+reportID.String()+"/claim", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403 Forbidden, got: %d", resp.StatusCode)
	}

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)

	if body["error"] != "Akses ditolak: Laporan berada di luar yurisdiksi wilayah kerja Anda" {
		t.Fatalf("unexpected error message: %v", body["error"])
	}
	if body["message"] != "akses ditolak: laporan di luar yurisdiksi wilayah kerja" {
		t.Fatalf("unexpected detail message: %v", body["message"])
	}

	// Invariant: no side effect
	if repo.report.InvestigatorID != nil {
		t.Fatalf("side effect detected! investigator_id changed to: %v", repo.report.InvestigatorID)
	}
}

// TEST-03: Direct API Bypass Attempt -> HTTP 403 Forbidden
func TestIMP001_HTTP_DirectAPIBypass_Forbidden(t *testing.T) {
	reportID := uuid.New()
	regionSukmajaya := uuid.New()
	regionCinere := uuid.New()
	kaderCinereID := uuid.New()

	repo := &mockHandlerReportRepo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionSukmajaya,
	}

	app := setupTestApp(repo)

	// Simulating direct curl / Postman without UI filtering
	token, err := generateTestToken(kaderCinereID, domain.RoleKader, &regionCinere)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/kader/reports/"+reportID.String()+"/claim", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("direct API bypass succeeded! expected 403 Forbidden, got: %d", resp.StatusCode)
	}

	// Invariant: no side effect
	if repo.report.InvestigatorID != nil {
		t.Fatalf("side effect detected! investigator_id changed: %v", repo.report.InvestigatorID)
	}
}

// TEST-04: Unauthorized Role -> HTTP 403 Forbidden (RBAC Guard)
func TestIMP001_HTTP_UnauthorizedRole_Forbidden(t *testing.T) {
	reportID := uuid.New()
	regionBeji := uuid.New()
	petugasPuskesmasID := uuid.New()

	repo := &mockHandlerReportRepo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionBeji,
	}

	app := setupTestApp(repo)

	// User is Petugas Puskesmas, NOT Kader
	token, err := generateTestToken(petugasPuskesmasID, domain.RolePetugasPuskesmas, &regionBeji)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/kader/reports/"+reportID.String()+"/claim", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403 Forbidden for non-Kader role, got: %d", resp.StatusCode)
	}

	// Invariant: no side effect
	if repo.report.InvestigatorID != nil {
		t.Fatalf("side effect detected! investigator_id set by unauthorized role: %v", repo.report.InvestigatorID)
	}
}

// TEST-06: Missing Location Context in Token -> HTTP 403 Forbidden
func TestIMP001_HTTP_MissingLocationContext_Forbidden(t *testing.T) {
	reportID := uuid.New()
	regionBeji := uuid.New()
	kaderNoLocID := uuid.New()

	repo := &mockHandlerReportRepo{
		report: &domain.CitizenReport{
			ID:             reportID,
			Status:         "MENUNGGU_PE_KADER",
			InvestigatorID: nil,
		},
		reportLocID: regionBeji,
	}

	app := setupTestApp(repo)

	// Token without location_id claim
	token, err := generateTestToken(kaderNoLocID, domain.RoleKader, nil)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/v1/kader/reports/"+reportID.String()+"/claim", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected status 403 Forbidden for missing location context, got: %d", resp.StatusCode)
	}

	var body map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	if body["message"] != "akses ditolak: akun tidak memiliki wilayah kerja yang valid" {
		t.Fatalf("unexpected message: %v", body["message"])
	}

	// Invariant: no side effect
	if repo.report.InvestigatorID != nil {
		t.Fatalf("side effect detected! investigator_id changed: %v", repo.report.InvestigatorID)
	}
}

// Unauthenticated Request -> HTTP 401 Unauthorized
func TestIMP001_HTTP_Unauthenticated_Unauthorized(t *testing.T) {
	reportID := uuid.New()
	repo := &mockHandlerReportRepo{
		report: &domain.CitizenReport{
			ID:     reportID,
			Status: "MENUNGGU_PE_KADER",
		},
	}
	app := setupTestApp(repo)

	req := httptest.NewRequest("POST", "/api/v1/kader/reports/"+reportID.String()+"/claim", nil)
	// No Authorization header

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized, got: %d", resp.StatusCode)
	}
}

// Malformed Report ID -> HTTP 400 Bad Request
func TestIMP001_HTTP_InvalidReportID_BadRequest(t *testing.T) {
	regionBeji := uuid.New()
	kaderID := uuid.New()
	repo := &mockHandlerReportRepo{}
	app := setupTestApp(repo)

	token, _ := generateTestToken(kaderID, domain.RoleKader, &regionBeji)

	req := httptest.NewRequest("POST", "/api/v1/kader/reports/invalid-uuid-string/claim", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request, got: %d", resp.StatusCode)
	}
}
