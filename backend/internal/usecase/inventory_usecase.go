package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/dbd-guard/backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryUsecase interface {
	CreateItem(ctx context.Context, itemType string, quantity float64, unit string, minStock float64, warehouse string) error
	GetItems(ctx context.Context) ([]domain.InventoryItem, error)
	UpdateStock(ctx context.Context, itemID, actorID uuid.UUID, quantityDelta float64, reason string) error
	UpdateStockTx(ctx context.Context, tx *gorm.DB, itemID, actorID uuid.UUID, quantityDelta float64, reason string) error

	// Stock Request & Transfer
	GetStockRequests(ctx context.Context, warehouse string) ([]domain.StockRequest, error)
	CreateStockRequest(ctx context.Context, itemID uuid.UUID, warehouse string, quantity float64, notes string) error
	ApproveStockRequest(ctx context.Context, actorID, requestID uuid.UUID, status string) error
	CreateStockTransfer(ctx context.Context, actorID, itemID uuid.UUID, fromWarehouse, toWarehouse string, quantity float64, notes string) error
	GetStockTransfers(ctx context.Context, toWarehouse string) ([]domain.StockTransfer, error)
	ReceiveStockTransfer(ctx context.Context, actorID, transferID uuid.UUID, isAccepted bool, notes string) error

	// Fogging Usage
	LogFoggingUsage(ctx context.Context, teamID, scheduleID, itemID uuid.UUID, quantity float64, notes string) error
	ApproveFoggingUsage(ctx context.Context, approverID, usageID uuid.UUID, isApproved bool, reason string) error
	GetPendingUsageLogs(ctx context.Context) ([]domain.InventoryUsageLog, error)
}

type inventoryUsecase struct {
	inventoryRepo repository.InventoryRepository
	auditRepo     repository.AuditRepository
}

func NewInventoryUsecase(inventoryRepo repository.InventoryRepository, auditRepo repository.AuditRepository) InventoryUsecase {
	return &inventoryUsecase{
		inventoryRepo: inventoryRepo,
		auditRepo:     auditRepo,
	}
}

func (u *inventoryUsecase) CreateItem(ctx context.Context, itemType string, quantity float64, unit string, minStock float64, warehouse string) error {
	item := &domain.InventoryItem{
		ItemType:     itemType,
		Quantity:     quantity,
		Unit:         unit,
		MinimumStock: minStock,
		Warehouse:    warehouse,
	}
	return u.inventoryRepo.CreateItem(ctx, item)
}

func (u *inventoryUsecase) GetItems(ctx context.Context) ([]domain.InventoryItem, error) {
	return u.inventoryRepo.GetItems(ctx)
}

func (u *inventoryUsecase) UpdateStock(ctx context.Context, itemID, actorID uuid.UUID, quantityDelta float64, reason string) error {
	if quantityDelta == 0 {
		return errors.New("perubahan stok tidak boleh 0")
	}

	return u.inventoryRepo.WithTransaction(func(tx *gorm.DB) error {
		return u.UpdateStockTx(ctx, tx, itemID, actorID, quantityDelta, reason)
	})
}

func (u *inventoryUsecase) UpdateStockTx(ctx context.Context, tx *gorm.DB, itemID, actorID uuid.UUID, quantityDelta float64, reason string) error {
	if quantityDelta == 0 {
		return errors.New("perubahan stok tidak boleh 0")
	}

	// 1. Get item with lock
	item, err := u.inventoryRepo.GetItemByIDForUpdate(ctx, tx, itemID)
	if err != nil {
		return err
	}

	if item.Quantity+quantityDelta < 0 {
		return errors.New("stok tidak mencukupi untuk dikurangi")
	}

	// 2. Update stock
	if err := u.inventoryRepo.UpdateStock(ctx, tx, itemID, quantityDelta); err != nil {
		return err
	}

	// 3. Audit Trail
	auditDetails := map[string]interface{}{
		"quantity_delta": quantityDelta,
		"reason":         reason,
	}

	if err := u.auditRepo.LogAction(ctx, tx, actorID, "UPDATE_INVENTORY_STOCK", "InventoryItem", itemID, auditDetails); err != nil {
		return err
	}

	return nil
}

func (u *inventoryUsecase) GetStockRequests(ctx context.Context, warehouse string) ([]domain.StockRequest, error) {
	return u.inventoryRepo.GetStockRequests(ctx, warehouse)
}

func (u *inventoryUsecase) GetStockTransfers(ctx context.Context, toWarehouse string) ([]domain.StockTransfer, error) {
	return u.inventoryRepo.GetStockTransfers(ctx, toWarehouse)
}

func (u *inventoryUsecase) CreateStockRequest(ctx context.Context, itemID uuid.UUID, warehouse string, quantity float64, notes string) error {
	req := &domain.StockRequest{
		ItemID:    itemID,
		Warehouse: warehouse,
		Quantity:  quantity,
		Status:    "PENDING",
		Notes:     notes,
	}
	return u.inventoryRepo.CreateStockRequest(ctx, req)
}

func (u *inventoryUsecase) ApproveStockRequest(ctx context.Context, actorID, requestID uuid.UUID, status string) error {
	if status != "APPROVED" && status != "REJECTED" {
		return errors.New("status must be APPROVED or REJECTED")
	}

	return u.inventoryRepo.WithTransaction(func(tx *gorm.DB) error {
		req, err := u.inventoryRepo.GetStockRequestByID(ctx, requestID)
		if err != nil {
			return err
		}

		if req.Status != "PENDING" {
			return errors.New("request already processed")
		}

		if err := u.inventoryRepo.UpdateStockRequestStatus(ctx, tx, requestID, status); err != nil {
			return err
		}

		if status == "APPROVED" {
			targetItem, err := u.inventoryRepo.GetItemByID(ctx, req.ItemID)
			if err != nil {
				return errors.New("item permintaan tidak ditemukan")
			}

			// Find source stock: prioritize Dinkes Central Warehouse, fallback to donor warehouse
			var sourceItem *domain.InventoryItem
			sourceItem, err = u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, targetItem.ItemType, "Gudang Utama Dinkes")
			if err != nil {
				sourceItem, err = u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, targetItem.ItemType, "Gudang Dinkes")
			}
			if err != nil {
				sourceItem, err = u.inventoryRepo.FindDonorWarehouse(ctx, tx, targetItem.ItemType, req.Quantity, req.Warehouse)
				if err != nil {
					return errors.New("stok di gudang pusat maupun donor tidak mencukupi untuk memenuhi permintaan")
				}
			}

			if (sourceItem.Quantity - sourceItem.ReservedQuantity) < req.Quantity {
				return errors.New("stok pengirim tidak cukup untuk memenuhi permintaan (mempertimbangkan reservasi)")
			}

			// Initiate StockTransfer (IN_TRANSIT) according to single accounting model
			transferNotes := fmt.Sprintf("Auto-dispatch transfer untuk Stock Request %s", req.ID)
			if req.Notes != "" {
				transferNotes += fmt.Sprintf(" (%s)", req.Notes)
			}

			transfer := &domain.StockTransfer{
				ItemID:        sourceItem.ItemID,
				FromWarehouse: sourceItem.Warehouse,
				ToWarehouse:   req.Warehouse,
				Quantity:      req.Quantity,
				Status:        "IN_TRANSIT",
				Notes:         transferNotes,
			}

			if err := u.inventoryRepo.CreateStockTransfer(ctx, tx, transfer); err != nil {
				return err
			}

			// Decrease source stock immediately to conserve total inventory
			if err := u.inventoryRepo.UpdateStock(ctx, tx, sourceItem.ItemID, -req.Quantity); err != nil {
				return err
			}
		}

		auditDetails := map[string]interface{}{
			"request_id": requestID,
			"status":     status,
		}
		if err := u.auditRepo.LogAction(ctx, tx, actorID, "PROCESS_STOCK_REQUEST", "StockRequest", requestID, auditDetails); err != nil {
			return err
		}

		return nil
	})
}

func (u *inventoryUsecase) CreateStockTransfer(ctx context.Context, actorID, itemID uuid.UUID, fromWarehouse, toWarehouse string, quantity float64, notes string) error {
	return u.inventoryRepo.WithTransaction(func(tx *gorm.DB) error {
		// Verify source item with row lock
		sourceItem, err := u.inventoryRepo.GetItemByIDForUpdate(ctx, tx, itemID)
		if err != nil {
			return errors.New("item asal tidak ditemukan")
		}
		
		if sourceItem.Warehouse != fromWarehouse {
			return errors.New("item asal bukan dari warehouse pengirim")
		}

		if (sourceItem.Quantity - sourceItem.ReservedQuantity) < quantity {
			return errors.New("stok pengirim tidak cukup (mempertimbangkan reservasi)")
		}

		// Create transfer record
		transfer := &domain.StockTransfer{
			ItemID:        itemID,
			FromWarehouse: fromWarehouse,
			ToWarehouse:   toWarehouse,
			Quantity:      quantity,
			Status:        "IN_TRANSIT",
			Notes:         notes,
		}
		
		if err := u.inventoryRepo.CreateStockTransfer(ctx, tx, transfer); err != nil {
			return err
		}

		// Decrease source stock
		if err := u.inventoryRepo.UpdateStock(ctx, tx, itemID, -quantity); err != nil {
			return err
		}

		// Do NOT increase destination stock yet (Wait for ReceiveStockTransfer)

		// Audit
		auditDetails := map[string]interface{}{
			"from_warehouse": fromWarehouse,
			"to_warehouse":   toWarehouse,
			"quantity":       quantity,
			"notes":          notes,
		}
		if err := u.auditRepo.LogAction(ctx, tx, actorID, "CREATE_STOCK_TRANSFER", "StockTransfer", transfer.ID, auditDetails); err != nil {
			return err
		}

		return nil
	})
}

func (u *inventoryUsecase) ReceiveStockTransfer(ctx context.Context, actorID, transferID uuid.UUID, isAccepted bool, notes string) error {
	return u.inventoryRepo.WithTransaction(func(tx *gorm.DB) error {
		transfer, err := u.inventoryRepo.GetStockTransferByID(ctx, transferID)
		if err != nil {
			return errors.New("data transfer tidak ditemukan")
		}

		if transfer.Status != "IN_TRANSIT" {
			return errors.New("hanya transfer dengan status IN_TRANSIT yang dapat diproses")
		}

		newStatus := "COMPLETED"
		if !isAccepted {
			newStatus = "REJECTED"
		}

		// Atomically update status from IN_TRANSIT to newStatus
		if err := u.inventoryRepo.UpdateStockTransferStatus(ctx, tx, transferID, newStatus); err != nil {
			return err
		}

		if !isAccepted {
			// Refund to source warehouse
			if err := u.inventoryRepo.UpdateStock(ctx, tx, transfer.ItemID, transfer.Quantity); err != nil {
				return err
			}
		} else {
			// Increase destination stock
			sourceItem, err := u.inventoryRepo.GetItemByID(ctx, transfer.ItemID)
			if err != nil {
				return err
			}

			destItem, err := u.inventoryRepo.GetItemByTypeAndWarehouse(ctx, sourceItem.ItemType, transfer.ToWarehouse)
			if err != nil {
				// Create new item for dest
				newItem := &domain.InventoryItem{
					ItemType:     sourceItem.ItemType,
					Quantity:     transfer.Quantity,
					Unit:         sourceItem.Unit,
					MinimumStock: sourceItem.MinimumStock,
					Warehouse:    transfer.ToWarehouse,
				}
				if err := u.inventoryRepo.CreateItemWithTx(ctx, tx, newItem); err != nil {
					return err
				}
			} else {
				if err := u.inventoryRepo.UpdateStock(ctx, tx, destItem.ItemID, transfer.Quantity); err != nil {
					return err
				}
			}
		}

		// Audit
		auditDetails := map[string]interface{}{
			"transfer_id": transferID,
			"status":      newStatus,
			"notes":       notes,
		}
		return u.auditRepo.LogAction(ctx, tx, actorID, "RECEIVE_STOCK_TRANSFER", "StockTransfer", transferID, auditDetails)
	})
}

func (u *inventoryUsecase) LogFoggingUsage(ctx context.Context, teamID, scheduleID, itemID uuid.UUID, quantity float64, notes string) error {
	if quantity <= 0 {
		return errors.New("jumlah logistik yang dilaporkan harus lebih besar dari 0")
	}

	return u.inventoryRepo.WithTransaction(func(tx *gorm.DB) error {
		usage := &domain.InventoryUsageLog{
			ItemID:     itemID,
			TeamID:     teamID,
			ScheduleID: scheduleID,
			Quantity:   quantity,
			Status:     "PENDING_APPROVAL",
			Notes:      notes,
		}

		return u.inventoryRepo.CreateInventoryUsageLog(ctx, tx, usage)
	})
}

func (u *inventoryUsecase) ApproveFoggingUsage(ctx context.Context, approverID, usageID uuid.UUID, isApproved bool, reason string) error {
	return u.inventoryRepo.WithTransaction(func(tx *gorm.DB) error {
		usage, err := u.inventoryRepo.GetInventoryUsageLogByID(ctx, usageID)
		if err != nil {
			return errors.New("log penggunaan tidak ditemukan")
		}

		if usage.Status != "PENDING_APPROVAL" {
			return errors.New("log penggunaan sudah diproses")
		}

		if isApproved {
			// Update status to APPROVED (audited review; stok sudah terpotong atomik saat selesai)
			if err := u.inventoryRepo.UpdateInventoryUsageLogStatus(ctx, tx, usageID, "APPROVED"); err != nil {
				return err
			}
		} else {
			if reason == "" {
				return errors.New("alasan penolakan wajib diisi")
			}
			if err := u.inventoryRepo.UpdateInventoryUsageLogStatus(ctx, tx, usageID, "REJECTED"); err != nil {
				return err
			}
			// if rejected, reserved stock is maintained until PE / schedule is handled separately?
			// Usually if usage is rejected, the team needs to resubmit. So we don't release reserved stock yet.
		}

		auditDetails := map[string]interface{}{
			"usage_id": usageID,
			"approved": isApproved,
			"reason":   reason,
		}
		return u.auditRepo.LogAction(ctx, tx, approverID, "APPROVE_FOGGING_USAGE", "InventoryUsageLog", usageID, auditDetails)
	})
}

func (u *inventoryUsecase) GetPendingUsageLogs(ctx context.Context) ([]domain.InventoryUsageLog, error) {
	return u.inventoryRepo.GetPendingUsageLogs(ctx)
}
