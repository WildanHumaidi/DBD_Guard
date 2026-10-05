package http

import (
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type EwsKlbHandler struct {
	usecase usecase.EwsKlbUsecase
}

func NewEwsKlbHandler(u usecase.EwsKlbUsecase) *EwsKlbHandler {
	return &EwsKlbHandler{usecase: u}
}

func (h *EwsKlbHandler) GetConfig(c *fiber.Ctx) error {
	config, err := h.usecase.GetConfiguration(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"data": config,
	})
}

func (h *EwsKlbHandler) UpdateConfig(c *fiber.Ctx) error {
	var req usecase.UpdateKlbConfigRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.usecase.UpdateConfiguration(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Configuration updated successfully",
	})
}

func (h *EwsKlbHandler) GetActiveAlerts(c *fiber.Ctx) error {
	alerts, err := h.usecase.GetActiveKLBAlerts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"data": alerts,
	})
}

func (h *EwsKlbHandler) GetEwsWarnings(c *fiber.Ctx) error {
	kelurahan := c.Query("kelurahan")
	rw := c.Query("rw")

	if kelurahan == "" || rw == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "kelurahan and rw query parameters are required",
		})
	}

	reports, err := h.usecase.GetEwsForRw(c.Context(), kelurahan, rw)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"data": reports,
	})
}
