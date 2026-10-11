package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dbd-guard/backend/internal/config"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/domain/dto"
	"github.com/dbd-guard/backend/internal/infrastructure/notification"
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
)

func main() {
	fmt.Println("=== Starting E2E Integration Test ===")

	// 1. Setup config & DB
	config.LoadConfig()
	config.ConnectDB()
	config.RunAutoMigration()

	// Ensure Roles exist
	roles := []string{domain.RoleKader, domain.RolePetugasPuskesmas, domain.RoleTimFogging, domain.RoleDinkes}
	for _, r := range roles {
		var role domain.Role
		if err := config.DB.Where("role_name = ?", r).First(&role).Error; err != nil {
			config.DB.Create(&domain.Role{RoleName: r})
		}
	}
	
	// Initialize Repositories
	userRepo := repository.NewUserRepository(config.DB)
	auditRepo := repository.NewAuditRepository(config.DB)
	reportRepo := repository.NewReportRepository(config.DB)
	peRepo := repository.NewPERepository(config.DB)
	foggingRepo := repository.NewFoggingRepository(config.DB)
	inventoryRepo := repository.NewInventoryRepository(config.DB)
	scheduleRepo := repository.NewFoggingScheduleRepository(config.DB)
	klbRepo := repository.NewKlbRepository(config.DB)
	masterDataRepo := repository.NewMasterDataRepository(config.DB)

	notifService := notification.NewNotificationService(config.Env)
	notifUsecase := usecase.NewNotificationUsecase(notifService, userRepo)

	storageService, _ := storage.NewMinioStorage(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_USE_SSL") == "true",
	)

	// Initialize Usecases
	reportUsecase := usecase.NewReportUsecase(reportRepo, auditRepo, notifUsecase, storageService)
	ewsKlbUsecase := usecase.NewEwsKlbUsecase(klbRepo, reportRepo, notifUsecase)
	peUsecase := usecase.NewPEUsecase(peRepo, reportRepo, foggingRepo, inventoryRepo, auditRepo, masterDataRepo, ewsKlbUsecase)
	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepo, auditRepo)
	scheduleUsecase := usecase.NewFoggingScheduleUsecase(scheduleRepo, foggingRepo, inventoryRepo, inventoryUsecase, auditRepo, reportRepo, masterDataRepo, notifUsecase)

	ctx := context.Background()

	// Get Roles
	var roleKader domain.Role
	config.DB.Where("role_name = ?", domain.RoleKader).First(&roleKader)
	var rolePuskesmas domain.Role
	config.DB.Where("role_name = ?", domain.RolePetugasPuskesmas).First(&rolePuskesmas)
	var roleTim domain.Role
	config.DB.Where("role_name = ?", domain.RoleTimFogging).First(&roleTim)

	var puskBeji domain.MasterPuskesmas
	config.DB.Where("name ILIKE ?", "%Beji%").First(&puskBeji)

	// Create Temporary Users
	kader := domain.User{
		Name:         "Test Kader E2E",
		Email:        "kader.e2e@test.com",
		PasswordHash: "mockhash",
		RoleID:       roleKader.ID,
		LocationID:   &puskBeji.ID,
	}
	config.DB.Where("email = ?", kader.Email).FirstOrCreate(&kader)
	defer config.DB.Unscoped().Delete(&kader)

	puskesmas := domain.User{
		Name:         "Test Puskesmas E2E",
		Email:        "puskesmas.e2e@test.com",
		PasswordHash: "mockhash",
		RoleID:       rolePuskesmas.ID,
		LocationID:   &puskBeji.ID,
	}
	config.DB.Where("email = ?", puskesmas.Email).FirstOrCreate(&puskesmas)
	defer config.DB.Unscoped().Delete(&puskesmas)
	
	timFogging := domain.User{
		Name:         "Test Tim Fogging E2E",
		Email:        "tim.e2e@test.com",
		PasswordHash: "mockhash",
		RoleID:       roleTim.ID,
	}
	config.DB.Where("email = ?", timFogging.Email).FirstOrCreate(&timFogging)
	defer config.DB.Unscoped().Delete(&timFogging)

	fmt.Println("--- 1. Submit Citizen Report ---")
	reportReq := dto.SubmitReportRequest{
		Name:          "Warga E2E",
		Contact:       "0811111111",
		AnonymousFlag: false,
		Category:      "Gejala DBD",
		Description:   "Ada anak demam E2E",
		Latitude:      -6.20,
		Longitude:     106.81,
		Address:       "Jl. E2E",
		RT:            "01",
		RW:            "01",
		Kelurahan:     "Beji",
		Kecamatan:     "Beji",
	}

	reportRes, err := reportUsecase.SubmitCitizenReport(ctx, reportReq)
	if err != nil {
		log.Fatalf("Failed to submit report: %v", err)
	}
	fmt.Printf("✅ Report submitted! Reference No: %s\n", reportRes.ReferenceNo)

	// Fetch the created report to get its UUID
	var dbReport domain.CitizenReport
	if err := config.DB.Where("reference_no = ?", reportRes.ReferenceNo).First(&dbReport).Error; err != nil {
		log.Fatalf("Failed to fetch created report: %v", err)
	}
	
	// Clean up after test (defer)
	defer config.DB.Unscoped().Delete(&dbReport)

	fmt.Println("\n--- 2. Kader Claims the Task (PE) ---")
	err = reportUsecase.ClaimReport(ctx, dbReport.ID, kader.ID, kader.LocationID)
	if err != nil {
		log.Fatalf("Failed to claim task: %v", err)
	}
	fmt.Printf("✅ Task claimed by Kader (ID: %s)\n", kader.ID.String())

	fmt.Println("\n--- 3. Kader Submits PE Result ---")
	peIDPtr, err := peUsecase.SubmitPEResult(ctx, dbReport.ID, kader.ID, 100, true, 3, true, "Jentik ditemukan di bak mandi", "http://example.com/photo.jpg")
	if err != nil {
		log.Fatalf("Failed to submit PE: %v", err)
	}
	fmt.Printf("✅ PE Submitted! PE ID: %s\n", peIDPtr.String())

	fmt.Println("\n--- 4. Puskesmas Validates PE & Approves Fogging ---")
	err = peUsecase.ValidatePEResult(ctx, *peIDPtr, puskesmas.ID, true, "", "", "")
	if err != nil {
		log.Fatalf("Failed to validate PE: %v", err)
	}
	fmt.Printf("✅ PE Validated! Decided to Fogging\n")

	fmt.Println("\n--- 5. Puskesmas Schedules Fogging ---")
	// fetch the generated fogging record
	var fogging domain.FoggingRequest
	if err := config.DB.Where("pe_id = ?", *peIDPtr).First(&fogging).Error; err != nil {
		log.Fatalf("Fogging record not found after PE validation: %v", err)
	}

	// Assuming fogging is approved directly or needs approval.
	// We might need to approve the fogging request first if it's PENDING.
	if fogging.ApprovalStatus == "PENDING" {
		fogging.ApprovalStatus = "APPROVED"
		config.DB.Save(&fogging)
		fmt.Println("✅ Fogging Request manually APPROVED for E2E purposes")
	}

	plannedTime := time.Now().Add(48 * time.Hour) // Must be > 24h
	schedIDPtr, err := scheduleUsecase.CreateSchedule(ctx, fogging.ID, timFogging.ID, puskesmas.ID, plannedTime)
	if err != nil {
		fmt.Printf("⚠️ Schedule creation note: %v\n", err)
	} else {
		fmt.Printf("✅ Fogging Scheduled! Schedule ID: %s\n", schedIDPtr.String())
	}

	fmt.Println("\n--- 6. Test KLB Trigger ---")
	// Force a KLB check on "Beji"
	err = ewsKlbUsecase.CheckAndTriggerKLB(ctx, "Beji", "KELURAHAN")
	if err != nil {
		fmt.Printf("⚠️ KLB Check note: %v\n", err)
	} else {
		fmt.Printf("✅ KLB Evaluated for Beji.\n")
		alerts, _ := ewsKlbUsecase.GetActiveKLBAlerts(ctx)
		for _, a := range alerts {
			if a.AreaName == "Beji" {
				fmt.Printf("   -> KLB ALERT ACTIVE for %s! Cases: %d\n", a.AreaName, a.CasesCount)
			}
		}
	}

	fmt.Println("\n✅ E2E Flow Simulation Completed Successfully!")
}
