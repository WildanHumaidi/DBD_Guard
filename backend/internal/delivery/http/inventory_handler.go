package http

import (
	"github.com/dbd-guard/backend/internal/delivery/http/middleware"
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	usecase   usecase.InventoryUsecase
	validator *validator.Validate
}

func NewInventoryHandler(usecase usecase.InventoryUsecase, validator *validator.Validate) *InventoryHandler {
	return &InventoryHandler{
		usecase:   usecase,
		validator: validator,
	}
}

func (h *InventoryHandler) CreateItem(c *fiber.Ctx) error {
	var req CreateInventoryItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.usecase.CreateItem(c.Context(), req.ItemType, req.Quantity, req.Unit, req.MinimumStock, req.Warehouse); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Item inventaris berhasil ditambahkan",
	})
}

// UI-07: Wrap response with standard wrapper for consistency
func (h *InventoryHandler) GetItems(c *fiber.Ctx) error {
	items, err := h.usecase.GetItems(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   items,
	})
}

func (h *InventoryHandler) UpdateStock(c *fiber.Ctx) error {
	var req UpdateStockRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	itemIDStr := c.Params("id")
	itemID, err := uuid.Parse(itemIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid item ID"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// SEC-01: Get actor ID from JWT context
	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized: " + err.Error()})
	}

	if err := h.usecase.UpdateStock(c.Context(), itemID, actorID, req.Quantity, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Stok berhasil diupdate",
	})
}

func (h *InventoryHandler) GetStockRequests(c *fiber.Ctx) error {
	warehouse := c.Query("warehouse")
	requests, err := h.usecase.GetStockRequests(c.Context(), warehouse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   requests,
	})
}

func (h *InventoryHandler) CreateStockRequest(c *fiber.Ctx) error {
	var req CreateStockRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	itemID, _ := uuid.Parse(req.ItemID)

	if err := h.usecase.CreateStockRequest(c.Context(), itemID, req.Warehouse, req.Quantity, req.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Permintaan stok berhasil diajukan",
	})
}

func (h *InventoryHandler) ApproveStockRequest(c *fiber.Ctx) error {
	var req ApproveStockRequestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	reqIDStr := c.Params("id")
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request ID"})
	}

	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.usecase.ApproveStockRequest(c.Context(), actorID, reqID, req.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Status permintaan stok berhasil diperbarui",
	})
}

func (h *InventoryHandler) CreateStockTransfer(c *fiber.Ctx) error {
	var req CreateStockTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	itemID, _ := uuid.Parse(req.ItemID)

	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.usecase.CreateStockTransfer(c.Context(), actorID, itemID, req.FromWarehouse, req.ToWarehouse, req.Quantity, req.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Alokasi silang berhasil dilakukan",
	})
}

func (h *InventoryHandler) GetStockTransfers(c *fiber.Ctx) error {
	warehouse := c.Query("warehouse")
	transfers, err := h.usecase.GetStockTransfers(c.Context(), warehouse)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   transfers,
	})
}

// Additional handlers
type ReceiveStockTransferRequest struct {
	IsAccepted bool   `json:"is_accepted"`
	Notes      string `json:"notes"`
}

func (h *InventoryHandler) ReceiveStockTransfer(c *fiber.Ctx) error {
	var req ReceiveStockTransferRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	transferIDStr := c.Params("id")
	transferID, err := uuid.Parse(transferIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid transfer ID"})
	}

	actorID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.usecase.ReceiveStockTransfer(c.Context(), actorID, transferID, req.IsAccepted, req.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Penerimaan transfer stok berhasil diproses",
	})
}

type LogFoggingUsageRequest struct {
	ScheduleID string  `json:"schedule_id" validate:"required"`
	ItemID     string  `json:"item_id" validate:"required"`
	Quantity   float64 `json:"quantity" validate:"required,gt=0"`
	Notes      string  `json:"notes"`
}

func (h *InventoryHandler) LogFoggingUsage(c *fiber.Ctx) error {
	var req LogFoggingUsageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	scheduleID, _ := uuid.Parse(req.ScheduleID)
	itemID, _ := uuid.Parse(req.ItemID)

	teamID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.usecase.LogFoggingUsage(c.Context(), teamID, scheduleID, itemID, req.Quantity, req.Notes); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Penggunaan stok fogging berhasil dicatat dan menunggu approval",
	})
}

type ApproveFoggingUsageRequest struct {
	IsApproved bool   `json:"is_approved"`
	Reason     string `json:"reason"`
}

func (h *InventoryHandler) ApproveFoggingUsage(c *fiber.Ctx) error {
	var req ApproveFoggingUsageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	usageIDStr := c.Params("id")
	usageID, err := uuid.Parse(usageIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid usage ID"})
	}

	approverID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.usecase.ApproveFoggingUsage(c.Context(), approverID, usageID, req.IsApproved, req.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Status log penggunaan stok fogging berhasil diperbarui",
	})
}

func (h *InventoryHandler) GetPendingUsageLogs(c *fiber.Ctx) error {
	logs, err := h.usecase.GetPendingUsageLogs(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   logs,
	})
}
