package http

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/domain/dto"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReportHandler struct {
	reportUsecase usecase.ReportUsecase
}

func NewReportHandler(reportUsecase usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{
		reportUsecase: reportUsecase,
	}
}

func (h *ReportHandler) UploadPhoto(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Gagal membaca file foto. Pastikan field bernama 'photo'.",
		})
	}

	// Call usecase to save file
	path, err := h.reportUsecase.UploadPhoto(c.Context(), fileHeader)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal menyimpan foto",
			"message": err.Error(),
		})
	}

	// Normalize path using forward slash
	path = filepath.ToSlash(path)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data": fiber.Map{
			"url": "/" + path,
		},
	})
}

func (h *ReportHandler) SubmitCitizenReport(c *fiber.Ctx) error {
	var req dto.SubmitReportRequest

	// Parse body inputs using BodyParser
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Error parsing form: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid form data",
		})
	}

	// For multipart forms or query/form params, check form value if not already set by BodyParser
	if !req.AnonymousFlag {
		anonVal := c.FormValue("anonymous_flag")
		if anonVal == "true" || anonVal == "1" {
			req.AnonymousFlag = true
		}
	}

	// Parse multipart form to get files
	form, err := c.MultipartForm()
	if err == nil && form != nil {
		if files, ok := form.File["photos"]; ok {
			req.Photos = files
		}
	}

	// Basic validation
	if req.Category == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Kategori laporan tidak boleh kosong",
		})
	}
	// Validate coordinates specifically for Depok's bounds
	// Depok Bounding Box: Lat -6.4526 to -6.3475, Lng 106.7262 to 106.8778
	if req.Latitude < -6.4526 || req.Latitude > -6.3475 || req.Longitude < 106.7262 || req.Longitude > 106.8778 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Layanan hanya tersedia untuk area Kota Depok. Pastikan koordinat GPS Anda akurat."})
	}

	// Call Usecase
	ctx := c.Context()
	res, err := h.reportUsecase.SubmitCitizenReport(ctx, req)
	if err != nil {
		log.Printf("Error submitting report: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Terjadi kesalahan internal saat memproses laporan",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   res,
	})
}

func (h *ReportHandler) GetReports(c *fiber.Ctx) error {
	ctx := c.Context()
	
	userRole, _ := c.Locals("role").(string)
	locationID, _ := middleware.GetLocationID(c)

	reports, err := h.reportUsecase.GetReports(ctx, userRole, locationID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   reports,
	})
}

func (h *ReportHandler) TrackReport(c *fiber.Ctx) error {
	ctx := c.Context()
	refNo := c.Params("refNo")
	
	report, err := h.reportUsecase.TrackReport(ctx, refNo)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "Laporan tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   report,
	})
}

func (h *ReportHandler) TrackReportsByPhone(c *fiber.Ctx) error {
	ctx := c.Context()
	phone := c.Params("phone")

	if phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nomor handphone tidak boleh kosong",
		})
	}

	summaries, err := h.reportUsecase.TrackReportsByPhone(ctx, phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal Server Error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   summaries,
	})
}

func (h *ReportHandler) GetReportByID(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")
	
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid report ID",
		})
	}

	report, err := h.reportUsecase.GetReportByID(ctx, id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "Laporan tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   report,
	})
}


type appealRequest struct {
	Reason string `json:"reason"`
}

func (h *ReportHandler) SubmitAppeal(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid report ID format",
		})
	}

	var req appealRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if req.Reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Alasan banding tidak boleh kosong",
		})
	}

	err = h.reportUsecase.SubmitAppeal(c.Context(), id, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal memproses banding",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Banding berhasil diajukan",
	})
}

type feedbackRequest struct {
	Rating int    `json:"rating"`
	Notes  string `json:"notes"`
}

func (h *ReportHandler) SubmitFeedback(c *fiber.Ctx) error {
	refNo := c.Params("refNo")
	if refNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nomor referensi laporan tidak boleh kosong",
		})
	}

	var req feedbackRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	err := h.reportUsecase.SubmitFeedback(c.Context(), refNo, req.Rating, req.Notes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal mensubmit feedback",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Terima kasih, evaluasi/feedback Anda berhasil disimpan",
	})
}

func (h *ReportHandler) ClaimReport(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid report ID format",
		})
	}

	// Get actor ID from JWT middleware
	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	locationID, _ := middleware.GetLocationID(c)

	err = h.reportUsecase.ClaimReport(c.Context(), id, actorID, locationID)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedJurisdiction) || errors.Is(err, domain.ErrMissingLocationContext) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   "Akses ditolak: Laporan berada di luar yurisdiksi wilayah kerja Anda",
				"message": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "laporan sudah diklaim") || strings.Contains(err.Error(), "status tidak valid") || strings.Contains(err.Error(), "tidak dalam status") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error":   "Konflik klaim tugas",
				"message": err.Error(),
			})
		}
		if strings.Contains(err.Error(), "record not found") || strings.Contains(err.Error(), "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   "Laporan tidak ditemukan",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal mengklaim tugas laporan",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengklaim tugas laporan",
	})
}

func (h *ReportHandler) CancelReport(c *fiber.Ctx) error {
	refNo := c.Params("refNo")
	if refNo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nomor referensi laporan tidak boleh kosong",
		})
	}

	err := h.reportUsecase.CancelReport(c.Context(), refNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal membatalkan laporan",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Laporan berhasil dibatalkan",
	})
}

type returnTaskRequest struct {
	Reason string `json:"reason"`
}

func (h *ReportHandler) ReturnReportTask(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid report ID format",
		})
	}

	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req returnTaskRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if req.Reason == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Alasan pengembalian tugas tidak boleh kosong",
		})
	}

	err = h.reportUsecase.ReturnReportTask(c.Context(), id, actorID, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal mengembalikan tugas laporan",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tugas laporan berhasil dikembalikan ke sistem",
	})
}

func (h *ReportHandler) ForceUnassignReport(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid report ID format",
		})
	}

	err = h.reportUsecase.ForceUnassignReport(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Gagal mencabut tugas laporan",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Tugas laporan berhasil dicabut dari kader",
	})
}

type TakeoverReportRequest struct {
	Reason string `json:"reason"`
}

func (h *ReportHandler) TakeoverReport(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid report ID format",
		})
	}

	officerID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var req TakeoverReportRequest
	_ = c.BodyParser(&req)

	if err := h.reportUsecase.TakeoverReport(c.Context(), id, officerID, req.Reason); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Gagal mengambil alih tugas laporan",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Laporan berhasil diambil alih oleh Puskesmas",
	})
}

func (h *ReportHandler) GetProtectedPhoto(c *fiber.Ctx) error {
	userRole := middleware.GetUserRole(c)
	if userRole == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Akses ditolak: Autentikasi diperlukan untuk mengakses dokumen/foto medis",
		})
	}

	userLocationID, _ := middleware.GetLocationID(c)
	filename := c.Params("filename")
	if filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Nama file tidak boleh kosong",
		})
	}

	presignedURL, err := h.reportUsecase.GetProtectedPhotoPath(c.Context(), userRole, userLocationID, filename)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "path traversal") || strings.Contains(errMsg, "tidak valid") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errMsg})
		}
		if strings.Contains(errMsg, "tidak ditemukan") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": errMsg})
		}
		if strings.Contains(errMsg, "akses ditolak") || strings.Contains(errMsg, "tidak memiliki izin") {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": errMsg})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": errMsg})
	}

	return c.Redirect(presignedURL, fiber.StatusTemporaryRedirect)
}

type mergeDuplicateRequest struct {
	PrimaryReportID string `json:"primary_report_id"`
	Reason          string `json:"reason"`
}

func (h *ReportHandler) MergeDuplicateReport(c *fiber.Ctx) error {
	childIDStr := c.Params("id")
	childID, err := uuid.Parse(childIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid child report ID format"})
	}

	var req mergeDuplicateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON format"})
	}

	primaryID, err := uuid.Parse(req.PrimaryReportID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid primary report ID format"})
	}

	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	err = h.reportUsecase.MergeDuplicateReport(c.Context(), childID, primaryID, actorID, req.Reason)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Gagal menggabungkan laporan duplikat",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Laporan duplikat berhasil digabungkan",
	})
}

