package http

import (
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/gofiber/fiber/v2"
)

type PublicHandler struct {
	reportRepo repository.ReportRepository
}

func NewPublicHandler(reportRepo repository.ReportRepository) *PublicHandler {
	return &PublicHandler{reportRepo: reportRepo}
}

// GetPublicStats returns aggregated statistics accessible without authentication
func (h *PublicHandler) GetPublicStats(c *fiber.Ctx) error {
	ctx := c.Context()

	reports, err := h.reportRepo.GetReports(ctx, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Gagal memuat statistik",
		})
	}

	totalReports := len(reports)
	totalCompleted := 0
	totalRejected := 0
	totalInProgress := 0
	totalFeedbackRating := 0
	feedbackCount := 0

	for _, r := range reports {
		switch r.Status {
		case "COMPLETED":
			totalCompleted++
			if r.FeedbackRating > 0 {
				totalFeedbackRating += r.FeedbackRating
				feedbackCount++
			}
		case "REJECTED":
			totalRejected++
		default:
			totalInProgress++
		}
	}

	avgRating := 0.0
	if feedbackCount > 0 {
		avgRating = float64(totalFeedbackRating) / float64(feedbackCount)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Statistik publik berhasil dimuat",
		"data": fiber.Map{
			"total_reports":     totalReports,
			"total_completed":   totalCompleted,
			"total_in_progress": totalInProgress,
			"total_rejected":    totalRejected,
			"avg_feedback":      avgRating,
			"feedback_count":    feedbackCount,
		},
	})
}
