package http

import (
	"strings"

	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PEHandler struct {
	usecase  usecase.PEUsecase
	validate *validator.Validate
}

func NewPEHandler(u usecase.PEUsecase, v *validator.Validate) *PEHandler {
	return &PEHandler{
		usecase:  u,
		validate: v,
	}
}

func (h *PEHandler) SubmitPEResult(c *fiber.Ctx) error {
	var req SubmitPEResultRequest
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

	// SEC-01: Get investigator ID from JWT context
	investigatorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}

	peID, err := h.usecase.SubmitPEResult(c.Context(), reportID, investigatorID, req.PeRadius, req.LarvaFound, req.FeverCount, req.PsnEdukasiDone, req.Notes, req.EvidencePhotoUrl)
	if err != nil {
		if strings.Contains(err.Error(), "UNAUTHORIZED_PE_SUBMISSION") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	resData := fiber.Map{
		"message": "Hasil PE berhasil disimpan",
	}
	if peID != nil {
		resData["pe_id"] = peID.String()
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   resData,
	})
}

type ValidatePEReq struct {
	IsApproved       bool   `json:"is_approved"`
	Reason           string `json:"reason"`
	OverrideDecision string `json:"override_decision"`
	OverrideReason   string `json:"override_reason"`
}

func (h *PEHandler) ValidatePEResult(c *fiber.Ctx) error {
	var req ValidatePEReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	peIDStr := c.Params("id")
	peID, err := uuid.Parse(peIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid PE ID"})
	}

	validatorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	validatorRole := middleware.GetUserRole(c)
	if validatorRole == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: role missing"})
	}

	// P0-03: If override is requested, verify role authorization
	if req.OverrideDecision != "" {
		if validatorRole != domain.RoleKepalaPuskesmas && validatorRole != domain.RoleDinkes && validatorRole != domain.RoleAdminDeveloper {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak: Hanya Kepala Puskesmas, Dinas Kesehatan, atau Super Admin yang berwenang melakukan override rekomendasi DSS",
			})
		}
		if req.OverrideReason == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Alasan override keputusan DSS wajib diisi",
			})
		}
	}

	if err := h.usecase.ValidatePEResult(c.Context(), peID, validatorID, req.IsApproved, req.Reason, req.OverrideDecision, req.OverrideReason); err != nil {
		if strings.Contains(err.Error(), "pemisahan tugas") || strings.Contains(err.Error(), "Four-Eyes Principle") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "alasan") || strings.Contains(err.Error(), "tidak valid") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Hasil PE berhasil divalidasi",
	})
}
