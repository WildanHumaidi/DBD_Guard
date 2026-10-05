package http

import (
	"strconv"

	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type AuditHandler struct {
	auditUsecase usecase.AuditUsecase
}

func NewAuditHandler(auditUsecase usecase.AuditUsecase) *AuditHandler {
	return &AuditHandler{auditUsecase}
}

func (h *AuditHandler) GetRecentLogs(c *fiber.Ctx) error {
	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	logs, err := h.auditUsecase.GetRecentLogs(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(logs)
}
