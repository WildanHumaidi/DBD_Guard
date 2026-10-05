package http

import (
	"github.com/dbd-guard/backend/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserManagementHandler struct {
	userMgmtUsecase usecase.UserManagementUsecase
}

func NewUserManagementHandler(userMgmtUsecase usecase.UserManagementUsecase) *UserManagementHandler {
	return &UserManagementHandler{
		userMgmtUsecase: userMgmtUsecase,
	}
}

func (h *UserManagementHandler) GetPuskesmasAccounts(c *fiber.Ctx) error {
	accounts, err := h.userMgmtUsecase.GetPuskesmasAccounts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(accounts)
}

func (h *UserManagementHandler) CreatePuskesmasAccount(c *fiber.Ctx) error {
	var req CreatePuskesmasAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" || req.Role == "" || req.PuskesmasID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, Email, Password, Role, and PuskesmasID are required"})
	}

	puskesmasID, err := uuid.Parse(req.PuskesmasID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid PuskesmasID format"})
	}

	actorIDStr := c.Locals("user_id")
	actorID, ok := actorIDStr.(uuid.UUID)
	if !ok {
		actorID = uuid.Nil
	}

	err = h.userMgmtUsecase.CreatePuskesmasAccount(c.Context(), actorID, req.Name, req.Email, req.Contact, req.Password, req.Role, puskesmasID)
	if err != nil {
		if err.Error() == "email sudah terdaftar" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Account created successfully"})
}

func (h *UserManagementHandler) UpdatePuskesmasAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	var req UpdatePuskesmasAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	var updatePuskesmasID *uuid.UUID
	if req.PuskesmasID != "" {
		parsed, err := uuid.Parse(req.PuskesmasID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid PuskesmasID format"})
		}
		updatePuskesmasID = &parsed
	}

	actorIDStr := c.Locals("user_id")
	actorID, ok := actorIDStr.(uuid.UUID)
	if !ok {
		actorID = uuid.Nil
	}

	err = h.userMgmtUsecase.UpdatePuskesmasAccount(c.Context(), actorID, userID, req.Name, req.Email, req.Contact, req.ActiveStatus, req.Role, updatePuskesmasID, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Account updated successfully"})
}

func (h *UserManagementHandler) DeletePuskesmasAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	actorIDStr := c.Locals("user_id")
	actorID, ok := actorIDStr.(uuid.UUID)
	if !ok {
		actorID = uuid.Nil
	}

	err = h.userMgmtUsecase.DeletePuskesmasAccount(c.Context(), actorID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Account deleted successfully"})
}

// Dinkes Account Handlers (For Admin Developer)
func (h *UserManagementHandler) GetDinkesAccounts(c *fiber.Ctx) error {
	accounts, err := h.userMgmtUsecase.GetDinkesAccounts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(accounts)
}

func (h *UserManagementHandler) CreateDinkesAccount(c *fiber.Ctx) error {
	var req CreateDinkesAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name, Email, and Password are required"})
	}

	actorIDStr := c.Locals("user_id")
	actorID, ok := actorIDStr.(uuid.UUID)
	if !ok {
		actorID = uuid.Nil
	}

	err := h.userMgmtUsecase.CreateDinkesAccount(c.Context(), actorID, req.Name, req.Email, req.Contact, req.Password)
	if err != nil {
		if err.Error() == "email sudah terdaftar" {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Dinkes Account created successfully"})
}

func (h *UserManagementHandler) UpdateDinkesAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	var req UpdateDinkesAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	actorIDStr := c.Locals("user_id")
	actorID, ok := actorIDStr.(uuid.UUID)
	if !ok {
		actorID = uuid.Nil
	}

	err = h.userMgmtUsecase.UpdateDinkesAccount(c.Context(), actorID, userID, req.Name, req.Email, req.Contact, req.ActiveStatus, req.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Dinkes Account updated successfully"})
}

func (h *UserManagementHandler) DeleteDinkesAccount(c *fiber.Ctx) error {
	idParam := c.Params("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	actorIDStr := c.Locals("user_id")
	actorID, ok := actorIDStr.(uuid.UUID)
	if !ok {
		actorID = uuid.Nil
	}

	err = h.userMgmtUsecase.DeleteDinkesAccount(c.Context(), actorID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Dinkes Account deleted successfully"})
}

func (h *UserManagementHandler) GetTimFoggingAccounts(c *fiber.Ctx) error {
	accounts, err := h.userMgmtUsecase.GetTimFoggingAccounts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(accounts)
}
