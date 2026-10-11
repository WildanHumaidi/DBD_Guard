package usecase_test

import (
	"context"
	"testing"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/google/uuid"
)

type mockKlbRepo struct {
	repository.KlbRepository
	config      *domain.KlbConfiguration
	activeAlert *domain.KlbAlert
	createdAlert *domain.KlbAlert
}

func (m *mockKlbRepo) GetConfiguration(ctx context.Context) (*domain.KlbConfiguration, error) {
	return m.config, nil
}

func (m *mockKlbRepo) GetActiveAlertByArea(ctx context.Context, areaName string) (*domain.KlbAlert, error) {
	return m.activeAlert, nil
}

func (m *mockKlbRepo) CreateAlert(ctx context.Context, alert *domain.KlbAlert) error {
	m.createdAlert = alert
	return nil
}

func (m *mockKlbRepo) ResolveAlert(ctx context.Context, alertID string) error {
	if m.activeAlert != nil {
		m.activeAlert.IsActive = false
	}
	return nil
}

type mockAreaCountReportRepo struct {
	repository.ReportRepository
	kelurahanCounts map[string]int
	kecamatanCounts map[string]int
}

func (m *mockAreaCountReportRepo) CountCasesByArea(ctx context.Context, areaName string, areaType string, days int) (int, error) {
	if areaType == "KELURAHAN" {
		return m.kelurahanCounts[areaName], nil
	}
	if areaType == "KECAMATAN" {
		return m.kecamatanCounts[areaName], nil
	}
	return 0, nil
}

func (m *mockAreaCountReportRepo) VerifyReportLocation(ctx context.Context, reportID uuid.UUID, locationID uuid.UUID) (bool, error) {
	return true, nil
}

func TestCheckAndTriggerKLB_KelurahanScope(t *testing.T) {
	klbRepo := &mockKlbRepo{
		config: &domain.KlbConfiguration{
			RegionScope:   "KELURAHAN",
			MinCases:      3,
			TimeframeDays: 7,
			IsAutoTrigger: true,
		},
	}
	repRepo := &mockAreaCountReportRepo{
		kelurahanCounts: map[string]int{
			"Sukatani": 4, // Exceeds threshold of 3
			"Sukamaju": 1, // Below threshold
		},
	}

	u := usecase.NewEwsKlbUsecase(klbRepo, repRepo, nil)

	ctx := context.Background()

	// 1. Check Sukatani -> should trigger KLB alert
	err := u.CheckAndTriggerKLB(ctx, "Sukatani", "KELURAHAN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if klbRepo.createdAlert == nil {
		t.Fatal("expected KLB alert to be created for Sukatani")
	}
	if klbRepo.createdAlert.AreaName != "Sukatani" || klbRepo.createdAlert.AreaType != "KELURAHAN" {
		t.Errorf("unexpected alert area: %s (%s)", klbRepo.createdAlert.AreaName, klbRepo.createdAlert.AreaType)
	}

	// 2. Check Sukamaju -> below threshold, should NOT trigger
	klbRepo.createdAlert = nil
	err = u.CheckAndTriggerKLB(ctx, "Sukamaju", "KELURAHAN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if klbRepo.createdAlert != nil {
		t.Errorf("expected no alert for Sukamaju, but got: %v", klbRepo.createdAlert)
	}
}

func TestCheckAndTriggerKLB_KecamatanScope(t *testing.T) {
	// Kecamatan Tapos contains multiple Kelurahans (e.g. Sukatani=2, Sukamaju=2, Cilangkap=1 -> Total=5)
	klbRepo := &mockKlbRepo{
		config: &domain.KlbConfiguration{
			RegionScope:   "KECAMATAN",
			MinCases:      5,
			TimeframeDays: 7,
			IsAutoTrigger: true,
		},
	}
	repRepo := &mockAreaCountReportRepo{
		kecamatanCounts: map[string]int{
			"Tapos": 5, // Sum of cases across all Kelurahans in Tapos
			"Beji":  2,
		},
	}

	u := usecase.NewEwsKlbUsecase(klbRepo, repRepo, nil)

	ctx := context.Background()

	// 1. Check Tapos -> total cases = 5 >= minCases (5) -> triggers KLB
	err := u.CheckAndTriggerKLB(ctx, "Tapos", "KECAMATAN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if klbRepo.createdAlert == nil {
		t.Fatal("expected KLB alert to be created for Kecamatan Tapos")
	}
	if klbRepo.createdAlert.AreaName != "Tapos" || klbRepo.createdAlert.AreaType != "KECAMATAN" {
		t.Errorf("unexpected alert area: %s (%s)", klbRepo.createdAlert.AreaName, klbRepo.createdAlert.AreaType)
	}

	// 2. Check Beji -> total cases = 2 < 5 -> should NOT trigger
	klbRepo.createdAlert = nil
	err = u.CheckAndTriggerKLB(ctx, "Beji", "KECAMATAN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if klbRepo.createdAlert != nil {
		t.Errorf("expected no alert for Beji, but got: %v", klbRepo.createdAlert)
	}
}

type mockKlbTrackingEwsUsecase struct {
	usecase.EwsKlbUsecase
	calledScopes []string
	calledAreas  []string
}

func (m *mockKlbTrackingEwsUsecase) CheckAndTriggerKLB(ctx context.Context, areaName string, areaType string) error {
	m.calledAreas = append(m.calledAreas, areaName)
	m.calledScopes = append(m.calledScopes, areaType)
	return nil
}

type mockNotifForKlbTest struct {
	usecase.NotificationUsecase
}

func (m *mockNotifForKlbTest) NotifyEWSKader(ctx context.Context, kelurahan, officerID string) error {
	return nil
}

func TestCheckAndTriggerKLB_WorkflowCaller_KecamatanIntegration(t *testing.T) {
	ctx := context.Background()

	trackingEws := &mockKlbTrackingEwsUsecase{}
	notif := &mockNotifForKlbTest{}

	reportID := uuid.New()
	officerID := uuid.New()
	locationID := uuid.New()

	report := &domain.CitizenReport{
		ID:     reportID,
		Status: "APPEALED",
		Location: domain.ReportLocation{
			Kelurahan: "Sukatani",
			Kecamatan: "Tapos",
		},
	}

	reportRepo := &mockFullLifecycleReportRepo{report: report}
	auditRepo := &mockAuditRepoLifecycle{}
	verificationRepo := &mockVerificationRepo{}

	vUsecase := usecase.NewVerificationUsecase(verificationRepo, reportRepo, auditRepo, trackingEws, notif)

	// Verify report as Puskesmas
	err := vUsecase.VerifyReport(ctx, reportID, officerID, domain.RolePetugasPuskesmas, &locationID, "APPROVED", "Valid")
	if err != nil {
		t.Fatalf("VerifyReport failed: %v", err)
	}

	// Verify that BOTH KELURAHAN and KECAMATAN were triggered by the caller
	hasKelurahan := false
	hasKecamatan := false
	for i, scope := range trackingEws.calledScopes {
		if scope == "KELURAHAN" && trackingEws.calledAreas[i] == "Sukatani" {
			hasKelurahan = true
		}
		if scope == "KECAMATAN" && trackingEws.calledAreas[i] == "Tapos" {
			hasKecamatan = true
		}
	}

	if !hasKecamatan {
		t.Fatalf("BUG-08 Residual! Kecamatan scope was not triggered by workflow caller. Called scopes: %v, areas: %v", trackingEws.calledScopes, trackingEws.calledAreas)
	}
	if !hasKelurahan {
		t.Fatalf("Expected Kelurahan scope to also be triggered. Called scopes: %v", trackingEws.calledScopes)
	}
}

// Suppress unused imports
var _ = uuid.New
