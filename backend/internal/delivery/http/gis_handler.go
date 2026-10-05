package http

import (
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

type GisHandler struct {
	gisUsecase usecase.GisUsecase
}

func NewGisHandler(gisUsecase usecase.GisUsecase) *GisHandler {
	return &GisHandler{gisUsecase: gisUsecase}
}

func (h *GisHandler) GetHeatmap(c *fiber.Ctx) error {
	ctx := c.Context()
	
	geoJSON, err := h.gisUsecase.GetHeatmapData(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Heatmap data retrieved successfully",
		"data":    geoJSON,
	})
}

func (h *GisHandler) GetRiskScore(c *fiber.Ctx) error {
	ctx := c.Context()
	
	riskScores, err := h.gisUsecase.GetRiskScores(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Risk scores retrieved successfully",
		"data":    riskScores,
	})
}
