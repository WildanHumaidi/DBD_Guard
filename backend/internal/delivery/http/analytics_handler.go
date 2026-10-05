package http

import (
	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	reportUsecase usecase.ReportUsecase
}

func NewAnalyticsHandler(reportUsecase usecase.ReportUsecase) *AnalyticsHandler {
	return &AnalyticsHandler{
		reportUsecase: reportUsecase,
	}
}

func (h *AnalyticsHandler) GetSLAMetrics(c *fiber.Ctx) error {
	role := middleware.GetUserRole(c)
	// Optionally check role if needed, e.g. only ADMIN_DINKES
	if role != "ADMIN_DINKES" && role != "DEVELOPER" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akses ditolak, khusus admin Dinkes",
		})
	}

	metrics, err := h.reportUsecase.GetSLAMetrics(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil metrik SLA",
		})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   metrics,
	})
}
