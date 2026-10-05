package http

import (
	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FoggingHandler struct {
	usecase  usecase.FoggingUsecase
	validate *validator.Validate
}

func NewFoggingHandler(u usecase.FoggingUsecase, v *validator.Validate) *FoggingHandler {
	return &FoggingHandler{
		usecase:  u,
		validate: v,
	}
}

func (h *FoggingHandler) ApproveFogging(c *fiber.Ctx) error {
	var req ApproveFoggingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	requestIDStr := c.Params("id")
	requestID, err := uuid.Parse(requestIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid fogging request ID"})
	}

	if err := h.validate.Struct(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// SEC-01: Get approver ID from JWT context
	approverID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}

	if err := h.usecase.ApproveFoggingRequest(c.Context(), requestID, approverID, req.Decision, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Keputusan fogging berhasil diproses",
	})
}
