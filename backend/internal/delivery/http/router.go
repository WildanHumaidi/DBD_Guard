package http

import (
	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes registers all API routes for the application
func SetupRoutes(
	app *fiber.App,
	authHandler *AuthHandler,
	reportHandler *ReportHandler,
	verificationHandler *VerificationHandler,
	peHandler *PEHandler,
	foggingHandler *FoggingHandler,
	inventoryHandler *InventoryHandler,
	scheduleHandler *FoggingScheduleHandler,
	gisHandler *GisHandler,
	publicHandler *PublicHandler,
	userMgmtHandler *UserManagementHandler,
	ewsKlbHandler *EwsKlbHandler,
	analyticsHandler *AnalyticsHandler,
	systemSettingHandler *SystemSettingHandler,
	auditHandler *AuditHandler,
) {
	api := app.Group("/api/v1")

	// Serve Public Static Files ONLY for non-sensitive assets (banners/infographics)
	app.Static("/storage/banners", "./storage/banners")

	// P0-02: Protect sensitive patient, citizen, and PE investigation photos
	app.Get("/storage/photos/:filename", middleware.Protected(), reportHandler.GetProtectedPhoto)
	
	// Create Rate Limiter for public routes
	rateLimiter := middleware.PublicRateLimiter()

	// PUBLIC ROUTES
	api.Post("/auth/login", rateLimiter, authHandler.Login)
	api.Post("/reports", rateLimiter, reportHandler.SubmitCitizenReport)
	api.Get("/reports/track/:refNo", rateLimiter, reportHandler.TrackReport)
	api.Get("/reports/phone/:phone", rateLimiter, reportHandler.TrackReportsByPhone)
	api.Post("/reports/:refNo/cancel", rateLimiter, reportHandler.CancelReport)
	api.Post("/reports/:id/appeal", rateLimiter, reportHandler.SubmitAppeal)
	api.Post("/reports/:refNo/feedback", rateLimiter, reportHandler.SubmitFeedback)
	api.Get("/public/stats", rateLimiter, publicHandler.GetPublicStats)
	api.Get("/public/settings", rateLimiter, systemSettingHandler.GetSettings)

	// PROTECTED ROUTES
	protected := api.Group("/", middleware.Protected())

	// Generic Uploads & Protected File Access
	protected.Post("/upload/photo", reportHandler.UploadPhoto)
	protected.Get("/storage/photos/:filename", reportHandler.GetProtectedPhoto)

	// Admin Developer Routes (Super Admin)
	developer := protected.Group("/developer", middleware.RequireRole(domain.RoleAdminDeveloper))
	developer.Get("/dinkes-accounts", userMgmtHandler.GetDinkesAccounts)
	developer.Post("/dinkes-accounts", userMgmtHandler.CreateDinkesAccount)
	developer.Patch("/dinkes-accounts/:id", userMgmtHandler.UpdateDinkesAccount)
	developer.Delete("/dinkes-accounts/:id", userMgmtHandler.DeleteDinkesAccount)
	developer.Get("/settings", systemSettingHandler.GetSettings)
	developer.Put("/settings", systemSettingHandler.UpdateSettings)
	developer.Get("/status", systemSettingHandler.GetStatus)
	developer.Get("/audit-logs", auditHandler.GetRecentLogs)

	// Admin Dinkes Routes
	dinkes := protected.Group("/admin", middleware.RequireRole(domain.RoleDinkes, domain.RoleAdminDeveloper))
	dinkes.Get("/puskesmas-accounts", userMgmtHandler.GetPuskesmasAccounts)
	dinkes.Post("/puskesmas-accounts", userMgmtHandler.CreatePuskesmasAccount)
	dinkes.Patch("/puskesmas-accounts/:id", userMgmtHandler.UpdatePuskesmasAccount)
	dinkes.Delete("/puskesmas-accounts/:id", userMgmtHandler.DeletePuskesmasAccount)
	dinkes.Get("/klb/config", ewsKlbHandler.GetConfig)
	dinkes.Patch("/klb/config", ewsKlbHandler.UpdateConfig)
	dinkes.Get("/klb/alerts", ewsKlbHandler.GetActiveAlerts)
	dinkes.Get("/analytics/sla", analyticsHandler.GetSLAMetrics)
	
	// Reports (Puskesmas can manually verify/reject a report from their dashboard before PE)
	protected.Get("/reports", reportHandler.GetReports)
	protected.Get("/reports/:id", reportHandler.GetReportByID)
	protected.Post("/reports/:id/verify", middleware.RequireRole(domain.RolePetugasPuskesmas, domain.RoleKepalaPuskesmas, domain.RoleKader), verificationHandler.VerifyReport)
	protected.Post("/reports/:id/merge", middleware.RequireRole(domain.RolePetugasPuskesmas, domain.RoleKepalaPuskesmas, domain.RoleKader, domain.RoleDinkes), reportHandler.MergeDuplicateReport)

	// GIS Heatmap (Dinkes)
	dinkes.Get("/gis/heatmap", gisHandler.GetHeatmap)
	dinkes.Get("/gis/risk-score", gisHandler.GetRiskScore)

	// PE Results (Kader submits, Puskesmas validates)
	kader := protected.Group("/kader", middleware.RequireRole(domain.RoleKader))
	kader.Post("/reports/:id/claim", reportHandler.ClaimReport)
	kader.Post("/reports/:id/return", reportHandler.ReturnReportTask)
	kader.Post("/reports/:id/pe-results", peHandler.SubmitPEResult)

	puskesmas := protected.Group("/puskesmas", middleware.RequireRole(domain.RolePetugasPuskesmas, domain.RoleKepalaPuskesmas))
	puskesmas.Get("/tim-fogging", userMgmtHandler.GetTimFoggingAccounts)
	puskesmas.Post("/pe-results/:id/validate", peHandler.ValidatePEResult)
	puskesmas.Post("/reports/:id/pe-results", peHandler.SubmitPEResult)
	puskesmas.Post("/reports/:id/unassign", reportHandler.ForceUnassignReport)
	puskesmas.Post("/reports/:id/takeover", reportHandler.TakeoverReport)

	// Fogging Requests (Kepala Puskesmas)
	kepalaPuskesmas := protected.Group("/kepala-puskesmas", middleware.RequireRole(domain.RoleKepalaPuskesmas))
	kepalaPuskesmas.Post("/fogging-requests/:id/approve", foggingHandler.ApproveFogging)

	// Fogging Schedules & Execution (Kepala Puskesmas & Dinkes BKO creates, Tim Fogging updates)
	protected.Post("/fogging-schedules", middleware.RequireRole(domain.RoleKepalaPuskesmas, domain.RoleDinkes), scheduleHandler.CreateSchedule)
	timFogging := protected.Group("/fogging", middleware.RequireRole(domain.RoleTimFogging))
	timFogging.Patch("/schedules/:id/status", scheduleHandler.UpdateExecutionStatus)
	timFogging.Post("/schedules/:id/inventory-usage", inventoryHandler.LogFoggingUsage)
	timFogging.Post("/schedules/:id/evidence", scheduleHandler.UploadEvidence)

	// Inventory
	protected.Post("/inventory", middleware.RequireRole(domain.RoleDinkes), inventoryHandler.CreateItem)
	protected.Get("/inventory", middleware.RequireRole(domain.RoleDinkes, domain.RolePetugasPuskesmas, domain.RoleKepalaPuskesmas, domain.RoleTimFogging), inventoryHandler.GetItems)
	protected.Patch("/inventory/:id/stock", middleware.RequireRole(domain.RoleDinkes, domain.RolePetugasPuskesmas, domain.RoleKepalaPuskesmas), inventoryHandler.UpdateStock) // For usage
	
	protected.Get("/inventory/requests", middleware.RequireRole(domain.RoleDinkes, domain.RoleKepalaPuskesmas, domain.RolePetugasPuskesmas), inventoryHandler.GetStockRequests)
	protected.Post("/inventory/requests", middleware.RequireRole(domain.RoleKepalaPuskesmas), inventoryHandler.CreateStockRequest)
	protected.Post("/inventory/requests/:id/approve", middleware.RequireRole(domain.RoleDinkes), inventoryHandler.ApproveStockRequest)
	
	protected.Get("/inventory/transfers", middleware.RequireRole(domain.RoleDinkes, domain.RolePetugasPuskesmas, domain.RoleKepalaPuskesmas), inventoryHandler.GetStockTransfers)
	protected.Post("/inventory/transfers", middleware.RequireRole(domain.RoleDinkes), inventoryHandler.CreateStockTransfer)
	protected.Post("/inventory/transfers/:id/receive", middleware.RequireRole(domain.RolePetugasPuskesmas, domain.RoleKepalaPuskesmas), inventoryHandler.ReceiveStockTransfer)
	
	protected.Get("/inventory/usage/history", middleware.RequireRole(domain.RoleKepalaPuskesmas, domain.RolePetugasPuskesmas, domain.RoleDinkes), inventoryHandler.GetPendingUsageLogs)
	protected.Get("/inventory/usage/pending", middleware.RequireRole(domain.RoleKepalaPuskesmas), inventoryHandler.GetPendingUsageLogs)
	protected.Post("/inventory/usage/:id/approve", middleware.RequireRole(domain.RoleKepalaPuskesmas), inventoryHandler.ApproveFoggingUsage)

	// EWS & KLB Warnings (Puskesmas & Kader)
	protected.Get("/ews-warnings", middleware.RequireRole(domain.RoleKader, domain.RolePetugasPuskesmas, domain.RoleKepalaPuskesmas), ewsKlbHandler.GetEwsWarnings)
}
