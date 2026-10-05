package http

import (
	"runtime"
	"time"

	"github.com/dbd-guard/backend/internal/config"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var serverStartTime = time.Now()

type SystemSettingHandler struct {
	settingUsecase usecase.SystemSettingUsecase
}

func NewSystemSettingHandler(settingUsecase usecase.SystemSettingUsecase) *SystemSettingHandler {
	return &SystemSettingHandler{settingUsecase}
}

func (h *SystemSettingHandler) GetSettings(c *fiber.Ctx) error {
	settings, err := h.settingUsecase.GetSettings(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(settings)
}

func (h *SystemSettingHandler) UpdateSettings(c *fiber.Ctx) error {
	heroTitle := c.FormValue("hero_title")
	announcementText := c.FormValue("announcement_text")
	isCitizenPortalActiveStr := c.FormValue("is_citizen_portal_active")
	isCitizenPortalActive := isCitizenPortalActiveStr == "true"

	if heroTitle == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "hero_title required"})
	}

	bannerImage, _ := c.FormFile("banner_image")

	userIDStr := c.Locals("user_id")
	userID, ok := userIDStr.(uuid.UUID)
	if !ok {
		// fallback just in case
		userID = uuid.Nil
	}

	err := h.settingUsecase.UpdateSettings(c.Context(), userID, heroTitle, announcementText, isCitizenPortalActive, bannerImage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Settings updated successfully"})
}

func (h *SystemSettingHandler) GetStatus(c *fiber.Ctx) error {
	uptime := time.Since(serverStartTime).String()

	// Check DB
	dbStatus := "Operational"
	if sqlDB, err := config.DB.DB(); err != nil || sqlDB.Ping() != nil {
		dbStatus = "Down"
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return c.JSON(fiber.Map{
		"uptime":        uptime,
		"uptime_seconds": time.Since(serverStartTime).Seconds(),
		"web_server":    "Operational",
		"database":      dbStatus,
		"fcm_service":   "Operational", // Mocking FCM status
		"memory_mb":     memStats.Alloc / 1024 / 1024,
	})
}
