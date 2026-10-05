package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getTestDB(t *testing.T) *gorm.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5433"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "secret"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "dbd_guard"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, pass, name, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping DB test: database not reachable: %v", err)
	}
	db.AutoMigrate(&domain.CitizenReport{})
	return db
}

func TestCountCasesByArea_KelurahanAndKecamatanScope(t *testing.T) {
	db := getTestDB(t)

	tx := db.Begin()
	defer tx.Rollback()

	repo := repository.NewReportRepository(tx)
	ctx := context.Background()

	testKec := "TestKecamatanAggregation"
	testKel1 := "TestKelurahanAgg1"
	testKel2 := "TestKelurahanAgg2"

	// 1. Setup MasterKelurahan entries under testKec
	pusk := domain.MasterPuskesmas{
		ID:        uuid.New(),
		Name:      "Puskesmas Test Agg",
		Kecamatan: testKec,
	}
	if err := tx.Create(&pusk).Error; err != nil {
		t.Fatalf("failed to create puskesmas: %v", err)
	}

	mKel1 := domain.MasterKelurahan{
		ID:          uuid.New(),
		Name:        testKel1,
		Kecamatan:   testKec,
		PuskesmasID: pusk.ID,
	}
	mKel2 := domain.MasterKelurahan{
		ID:          uuid.New(),
		Name:        testKel2,
		Kecamatan:   testKec,
		PuskesmasID: pusk.ID,
	}
	if err := tx.Create(&mKel1).Error; err != nil {
		t.Fatalf("failed to create mKel1: %v", err)
	}
	if err := tx.Create(&mKel2).Error; err != nil {
		t.Fatalf("failed to create mKel2: %v", err)
	}

	// 2. Create Reporter
	reporter := domain.Reporter{
		ID:            uuid.New(),
		Name:          "Agg Tester",
		Contact:       "08123456789",
		AnonymousFlag: false,
	}
	if err := tx.Create(&reporter).Error; err != nil {
		t.Fatalf("failed to create reporter: %v", err)
	}

	// 3. Create 2 reports in Kelurahan 1, and 1 report in Kelurahan 2
	createTestReport := func(kel, kec string) {
		reportID := uuid.New()
		rep := domain.CitizenReport{
			ID:          reportID,
			ReferenceNo: fmt.Sprintf("REF-%s", reportID.String()[:8]),
			Category:    "JENTIK_NYAMUK",
			Description: "Test aggregation",
			Status:      "VERIFIED",
			ReporterID:  reporter.ID,
			CreatedAt:   time.Now(),
		}
		if err := tx.Create(&rep).Error; err != nil {
			t.Fatalf("failed to create report: %v", err)
		}
		loc := domain.ReportLocation{
			ID:        uuid.New(),
			ReportID:  reportID,
			Kelurahan: kel,
			Kecamatan: kec,
			Latitude:  -6.4,
			Longitude: 106.8,
			Geom:      "SRID=4326;POINT(106.8 -6.4)",
		}
		if err := tx.Create(&loc).Error; err != nil {
			t.Fatalf("failed to create location: %v", err)
		}
	}

	createTestReport(testKel1, testKec)
	createTestReport(testKel1, testKec)
	createTestReport(testKel2, testKec)

	// 4. Test Kelurahan 1 count: should be 2
	countKel1, err := repo.CountCasesByArea(ctx, testKel1, "KELURAHAN", 7)
	if err != nil {
		t.Fatalf("CountCasesByArea Kelurahan 1 failed: %v", err)
	}
	if countKel1 != 2 {
		t.Errorf("expected 2 cases in %s, got: %d", testKel1, countKel1)
	}

	// 5. Test Kelurahan 2 count: should be 1
	countKel2, err := repo.CountCasesByArea(ctx, testKel2, "KELURAHAN", 7)
	if err != nil {
		t.Fatalf("CountCasesByArea Kelurahan 2 failed: %v", err)
	}
	if countKel2 != 1 {
		t.Errorf("expected 1 case in %s, got: %d", testKel2, countKel2)
	}

	// 6. Test Kecamatan count: must cover all Kelurahans under that Kecamatan (2 + 1 = 3)
	countKec, err := repo.CountCasesByArea(ctx, testKec, "KECAMATAN", 7)
	if err != nil {
		t.Fatalf("CountCasesByArea Kecamatan failed: %v", err)
	}
	if countKec != 3 {
		t.Errorf("expected 3 total cases in %s (sum of Kel1 and Kel2), got: %d", testKec, countKec)
	}
}

func TestVerifyReportLocation(t *testing.T) {
	db := getTestDB(t)
	tx := db.Begin()
	defer tx.Rollback()

	repo := repository.NewReportRepository(tx)
	ctx := context.Background()

	pusk1 := domain.MasterPuskesmas{
		ID:        uuid.New(),
		Name:      "Puskesmas Beji Test",
		Kecamatan: "Beji",
	}
	pusk2 := domain.MasterPuskesmas{
		ID:        uuid.New(),
		Name:      "Puskesmas Pancoran Mas Test",
		Kecamatan: "Pancoran Mas",
	}
	tx.Create(&pusk1)
	tx.Create(&pusk2)

	kel1 := domain.MasterKelurahan{
		ID:          uuid.New(),
		Name:        "Kelurahan Beji Test",
		Kecamatan:   "Beji",
		PuskesmasID: pusk1.ID,
	}
	kel2 := domain.MasterKelurahan{
		ID:          uuid.New(),
		Name:        "Kelurahan Depok Jaya Test",
		Kecamatan:   "Pancoran Mas",
		PuskesmasID: pusk2.ID,
	}
	tx.Create(&kel1)
	tx.Create(&kel2)

	reporter := domain.Reporter{
		ID:            uuid.New(),
		Name:          "Test Reporter",
		Contact:       "0812345678",
		AnonymousFlag: false,
	}
	tx.Create(&reporter)

	repID := uuid.New()
	rep := domain.CitizenReport{
		ID:          repID,
		ReferenceNo: "REF-TEST-001",
		Category:    "JENTIK_NYAMUK",
		Status:      "MENUNGGU_PE_KADER",
		ReporterID:  reporter.ID,
		CreatedAt:   time.Now(),
	}
	tx.Create(&rep)

	loc := domain.ReportLocation{
		ID:        uuid.New(),
		ReportID:  repID,
		Kelurahan: kel1.Name,
		Kecamatan: kel1.Kecamatan,
		Latitude:  -6.4,
		Longitude: 106.8,
		Geom:      "SRID=4326;POINT(106.8 -6.4)",
	}
	tx.Create(&loc)

	// Test same puskesmas
	ok1, err1 := repo.VerifyReportLocation(ctx, repID, pusk1.ID)
	if err1 != nil || !ok1 {
		t.Fatalf("expected true for same puskesmas, got %v, err=%v", ok1, err1)
	}

	// Test different puskesmas
	ok2, err2 := repo.VerifyReportLocation(ctx, repID, pusk2.ID)
	if err2 != nil || ok2 {
		t.Fatalf("expected false for different puskesmas, got %v, err=%v", ok2, err2)
	}

	// Test same kelurahan ID
	okKel1, errKel1 := repo.VerifyReportLocation(ctx, repID, kel1.ID)
	if errKel1 != nil || !okKel1 {
		t.Fatalf("expected true for same kelurahan, got %v, err=%v", okKel1, errKel1)
	}

	// Test different kelurahan ID
	okKel2, errKel2 := repo.VerifyReportLocation(ctx, repID, kel2.ID)
	if errKel2 != nil || okKel2 {
		t.Fatalf("expected false for different kelurahan, got %v, err=%v", okKel2, errKel2)
	}
}


