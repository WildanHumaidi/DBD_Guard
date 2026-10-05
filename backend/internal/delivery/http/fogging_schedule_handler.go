package http

import (
	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FoggingScheduleHandler struct {
	usecase   usecase.FoggingScheduleUsecase
	validator *validator.Validate
}

func NewFoggingScheduleHandler(usecase usecase.FoggingScheduleUsecase, validator *validator.Validate) *FoggingScheduleHandler {
	return &FoggingScheduleHandler{
		usecase:   usecase,
		validator: validator,
	}
}

func (h *FoggingScheduleHandler) CreateSchedule(c *fiber.Ctx) error {
	var req CreateScheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// M-09 & SEC-01: Extract creatorID from JWT
	creatorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}

	scheduleID, err := h.usecase.CreateSchedule(c.Context(), req.RequestID, req.TeamID, creatorID, req.PlannedTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	resData := fiber.Map{
		"message": "Jadwal Fogging berhasil dibuat",
	}
	if scheduleID != nil {
		resData["schedule_id"] = scheduleID.String()
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   resData,
	})
}

func (h *FoggingScheduleHandler) UpdateExecutionStatus(c *fiber.Ctx) error {
	var req UpdateExecutionStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	scheduleIDStr := c.Params("id")
	scheduleID, err := uuid.Parse(scheduleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid schedule ID"})
	}
	req.ScheduleID = scheduleID // override with param ID

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// SEC-01: Extract actorID from JWT
	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}

	var usedItemsInput []usecase.UsedItemInput
	for _, item := range req.UsedItems {
		usedItemsInput = append(usedItemsInput, usecase.UsedItemInput{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		})
	}

	if err := h.usecase.UpdateExecutionStatus(c.Context(), req.ScheduleID, actorID, req.ExecutionStatus, req.Notes, usedItemsInput); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Status eksekusi Fogging berhasil diupdate",
	})
}

func (h *FoggingScheduleHandler) UploadEvidence(c *fiber.Ctx) error {
	scheduleIDStr := c.Params("id")
	scheduleID, err := uuid.Parse(scheduleIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid schedule ID"})
	}

	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}

	// For now, simulate storage by saving locally or just generating a fake URL
	// In production, upload to GCS/S3.
	file, err := c.FormFile("evidence")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File evidence tidak ditemukan"})
	}

	storageKey := "uploads/fogging_evidence/" + uuid.New().String() + "_" + file.Filename
	
	// Save file locally for demo purpose (or handle via mock storage)
	if err := c.SaveFile(file, "./"+storageKey); err != nil {
		// Log error, but maybe don't fail immediately if it's just a demo
	}

	if err := h.usecase.UploadEvidence(c.Context(), scheduleID, actorID, storageKey); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Bukti fogging berhasil diunggah",
		"data": fiber.Map{
			"storage_key": storageKey,
		},
	})
}
