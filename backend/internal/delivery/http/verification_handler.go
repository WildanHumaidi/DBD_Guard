package http

import (
	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VerificationHandler struct {
	usecase  usecase.VerificationUsecase
	validate *validator.Validate
}

func NewVerificationHandler(u usecase.VerificationUsecase, v *validator.Validate) *VerificationHandler {
	return &VerificationHandler{
		usecase:  u,
		validate: v,
	}
}

func (h *VerificationHandler) VerifyReport(c *fiber.Ctx) error {
	var req VerifyReportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	reportIDStr := c.Params("id")
	reportID, err := uuid.Parse(reportIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid report ID"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// SEC-01: Get actor ID, role, and location ID from JWT context
	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}
	
	actorRole, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: missing role"})
	}
	
	locationID, _ := middleware.GetLocationID(c)

	if err := h.usecase.VerifyReport(c.Context(), reportID, actorID, actorRole, locationID, req.Decision, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Verifikasi laporan berhasil diproses",
	})
}
