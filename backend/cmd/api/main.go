package main

import (
	"log"
	"os"

	"github.com/dbd-guard/backend/internal/config"
	"github.com/dbd-guard/backend/internal/delivery/cron"
	httpDelivery "github.com/dbd-guard/backend/internal/delivery/http"
	"github.com/dbd-guard/backend/internal/infrastructure/notification"
	"github.com/dbd-guard/backend/internal/infrastructure/storage"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load Configuration
	config.LoadConfig()

	// Initialize Database Connection
	// This will log a warning if the database is not up yet
	config.ConnectDB()

	// Run Database Migrations
	config.RunAutoMigration()

	// Initialize Fiber App
	app := fiber.New(fiber.Config{
		AppName: "DBD Guard API",
	})

	// Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH, OPTIONS",
	}))
	app.Use(logger.New())

	// Initialize Modules
	userRepo := repository.NewUserRepository(config.DB)
	auditRepo := repository.NewAuditRepository(config.DB)
	reportRepo := repository.NewReportRepository(config.DB)
	verificationRepo := repository.NewVerificationRepository(config.DB)
	peRepo := repository.NewPERepository(config.DB)
	foggingRepo := repository.NewFoggingRepository(config.DB)
	inventoryRepo := repository.NewInventoryRepository(config.DB)
	scheduleRepo := repository.NewFoggingScheduleRepository(config.DB)
	gisRepo := repository.NewGisRepository(config.DB)
	klbRepo := repository.NewKlbRepository(config.DB)
	masterDataRepo := repository.NewMasterDataRepository(config.DB)
	systemSettingRepo := repository.NewSystemSettingRepository(config.DB)

	notifService := notification.NewNotificationService(config.Env)
	notifUsecase := usecase.NewNotificationUsecase(notifService, userRepo)

	// Initialize Storage
	var storageService storage.StorageService
	minioStorage, err := storage.NewMinioStorage(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
		os.Getenv("MINIO_USE_SSL") == "true",
	)
	if err != nil {
		log.Printf("Failed to initialize MinIO storage, fallback to mock: %v", err)
		storageService = &storage.MockStorageService{}
	} else {
		storageService = minioStorage
	}

	authUsecase := usecase.NewAuthUsecase(userRepo)
	reportUsecase := usecase.NewReportUsecase(reportRepo, auditRepo, notifUsecase, storageService)
	foggingUsecase := usecase.NewFoggingUsecase(foggingRepo, auditRepo, reportRepo, inventoryRepo, masterDataRepo)
	inventoryUsecase := usecase.NewInventoryUsecase(inventoryRepo, auditRepo)
	scheduleUsecase := usecase.NewFoggingScheduleUsecase(scheduleRepo, foggingRepo, inventoryRepo, inventoryUsecase, auditRepo, reportRepo, masterDataRepo, notifUsecase)
	gisUsecase := usecase.NewGisUsecase(gisRepo)
	userMgmtUsecase := usecase.NewUserManagementUsecase(userRepo, auditRepo)
	auditUsecase := usecase.NewAuditUsecase(auditRepo)
	
	ewsKlbUsecase := usecase.NewEwsKlbUsecase(klbRepo, reportRepo, notifUsecase)
	verificationUsecase := usecase.NewVerificationUsecase(verificationRepo, reportRepo, auditRepo, ewsKlbUsecase, notifUsecase)
	peUsecase := usecase.NewPEUsecase(peRepo, reportRepo, foggingRepo, inventoryRepo, auditRepo, masterDataRepo, ewsKlbUsecase)
	systemSettingUsecase := usecase.NewSystemSettingUsecase(systemSettingRepo, auditRepo)

	// Validate instance
	validate := validator.New()

	authHandler := httpDelivery.NewAuthHandler(authUsecase, validate)
	reportHandler := httpDelivery.NewReportHandler(reportUsecase)
	verificationHandler := httpDelivery.NewVerificationHandler(verificationUsecase, validate)
	peHandler := httpDelivery.NewPEHandler(peUsecase, validate)
	foggingHandler := httpDelivery.NewFoggingHandler(foggingUsecase, validate)
	inventoryHandler := httpDelivery.NewInventoryHandler(inventoryUsecase, validate)
	scheduleHandler := httpDelivery.NewFoggingScheduleHandler(scheduleUsecase, validate)
	gisHandler := httpDelivery.NewGisHandler(gisUsecase)
	publicHandler := httpDelivery.NewPublicHandler(reportRepo)
	userMgmtHandler := httpDelivery.NewUserManagementHandler(userMgmtUsecase)
	auditHandler := httpDelivery.NewAuditHandler(auditUsecase)
	ewsKlbHandler := httpDelivery.NewEwsKlbHandler(ewsKlbUsecase)
	analyticsHandler := httpDelivery.NewAnalyticsHandler(reportUsecase)
	systemSettingHandler := httpDelivery.NewSystemSettingHandler(systemSettingUsecase)

	// Basic route to test server
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "DBD Guard API is running",
		})
	})

	// Start SLA Cron Background Job
	slaCron := cron.NewSLACron(reportUsecase)
	slaCron.Start()

	// Start Feedback Reminder Cron Background Job
	feedbackCron := cron.NewFeedbackCron(reportUsecase)
	feedbackCron.Start()
	
	// Setup API Routes
	httpDelivery.SetupRoutes(app, authHandler, reportHandler, verificationHandler, peHandler, foggingHandler, inventoryHandler, scheduleHandler, gisHandler, publicHandler, userMgmtHandler, ewsKlbHandler, analyticsHandler, systemSettingHandler, auditHandler)

	// Start Server
	log.Printf("Starting server in %s mode on port %s...", config.Env.AppEnv, config.Env.Port)
	if err := app.Listen(":" + config.Env.Port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
