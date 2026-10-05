package repository

import (
	"context"
	"errors"

	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	CreateItem(ctx context.Context, item *domain.InventoryItem) error
	CreateItemWithTx(ctx context.Context, tx *gorm.DB, item *domain.InventoryItem) error
	GetItems(ctx context.Context) ([]domain.InventoryItem, error)
	GetItemByID(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error)
	GetItemByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.InventoryItem, error)
	GetItemByTypeAndWarehouse(ctx context.Context, itemType, warehouse string) (*domain.InventoryItem, error)
	UpdateStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantityDelta float64) error
	ReserveStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error
	ReleaseReservedStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error
	FindDonorWarehouse(ctx context.Context, tx *gorm.DB, itemType string, requiredAmount float64, excludeWarehouse string) (*domain.InventoryItem, error)
	WithTransaction(fn func(tx *gorm.DB) error) error

	// Stock Request & Transfer
	CreateStockRequest(ctx context.Context, req *domain.StockRequest) error
	GetStockRequests(ctx context.Context, warehouse string) ([]domain.StockRequest, error)
	GetStockRequestByID(ctx context.Context, id uuid.UUID) (*domain.StockRequest, error)
	UpdateStockRequestStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error
	CreateStockTransfer(ctx context.Context, tx *gorm.DB, transfer *domain.StockTransfer) error
	GetStockTransferByID(ctx context.Context, id uuid.UUID) (*domain.StockTransfer, error)
	GetStockTransfers(ctx context.Context, toWarehouse string) ([]domain.StockTransfer, error)
	UpdateStockTransferStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error

	// Inventory Usage Log
	CreateInventoryUsageLog(ctx context.Context, tx *gorm.DB, log *domain.InventoryUsageLog) error
	GetInventoryUsageLogByID(ctx context.Context, id uuid.UUID) (*domain.InventoryUsageLog, error)
	UpdateInventoryUsageLogStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error
	GetPendingUsageLogs(ctx context.Context) ([]domain.InventoryUsageLog, error)
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) CreateItem(ctx context.Context, item *domain.InventoryItem) error {
	return r.CreateItemWithTx(ctx, nil, item)
}

func (r *inventoryRepository) CreateItemWithTx(ctx context.Context, tx *gorm.DB, item *domain.InventoryItem) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Create(item).Error
}

func (r *inventoryRepository) GetItems(ctx context.Context) ([]domain.InventoryItem, error) {
	var items []domain.InventoryItem
	err := r.db.WithContext(ctx).Find(&items).Error
	return items, err
}

func (r *inventoryRepository) GetItemByID(ctx context.Context, id uuid.UUID) (*domain.InventoryItem, error) {
	var item domain.InventoryItem
	err := r.db.WithContext(ctx).First(&item, "item_id = ?", id).Error
	return &item, err
}

func (r *inventoryRepository) GetItemByIDForUpdate(ctx context.Context, tx *gorm.DB, id uuid.UUID) (*domain.InventoryItem, error) {
	var item domain.InventoryItem
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	// clause.Locking{Strength: "UPDATE"} is equivalent to FOR UPDATE
	err := execDB.WithContext(ctx).Set("gorm:query_option", "FOR UPDATE").First(&item, "item_id = ?", id).Error
	return &item, err
}

func (r *inventoryRepository) GetItemByTypeAndWarehouse(ctx context.Context, itemType, warehouse string) (*domain.InventoryItem, error) {
	var item domain.InventoryItem
	err := r.db.WithContext(ctx).Where("item_type = ? AND warehouse = ?", itemType, warehouse).First(&item).Error
	return &item, err
}

func (r *inventoryRepository) UpdateStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantityDelta float64) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}

	// Uses gorm.Expr to atomically update stock to prevent race conditions
	return execDB.WithContext(ctx).
		Model(&domain.InventoryItem{}).
		Where("item_id = ?", itemID).
		Update("quantity", gorm.Expr("quantity + ?", quantityDelta)).Error
}

func (r *inventoryRepository) ReserveStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}

	return execDB.WithContext(ctx).
		Model(&domain.InventoryItem{}).
		Where("item_id = ?", itemID).
		Update("reserved_quantity", gorm.Expr("reserved_quantity + ?", quantity)).Error
}

func (r *inventoryRepository) ReleaseReservedStock(ctx context.Context, tx *gorm.DB, itemID uuid.UUID, quantity float64) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}

	return execDB.WithContext(ctx).
		Model(&domain.InventoryItem{}).
		Where("item_id = ?", itemID).
		Update("reserved_quantity", gorm.Expr("reserved_quantity - ?", quantity)).Error
}



func (r *inventoryRepository) FindDonorWarehouse(ctx context.Context, tx *gorm.DB, itemType string, requiredAmount float64, excludeWarehouse string) (*domain.InventoryItem, error) {
	var item domain.InventoryItem
	execDB := r.db
	if tx != nil {
		execDB = tx
	}

	err := execDB.WithContext(ctx).
		Set("gorm:query_option", "FOR UPDATE").
		Where("item_type = ? AND warehouse != ? AND (quantity - reserved_quantity - minimum_stock) >= ?", itemType, excludeWarehouse, requiredAmount).
		First(&item).Error

	return &item, err
}

func (r *inventoryRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *inventoryRepository) CreateStockRequest(ctx context.Context, req *domain.StockRequest) error {
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *inventoryRepository) GetStockRequests(ctx context.Context, warehouse string) ([]domain.StockRequest, error) {
	var requests []domain.StockRequest
	query := r.db.WithContext(ctx)
	if warehouse != "" {
		query = query.Where("warehouse = ?", warehouse)
	}
	err := query.Order("created_at desc").Find(&requests).Error
	return requests, err
}

func (r *inventoryRepository) GetStockRequestByID(ctx context.Context, id uuid.UUID) (*domain.StockRequest, error) {
	var req domain.StockRequest
	err := r.db.WithContext(ctx).First(&req, "id = ?", id).Error
	return &req, err
}

func (r *inventoryRepository) UpdateStockRequestStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	res := execDB.WithContext(ctx).Model(&domain.StockRequest{}).Where("id = ? AND status = 'PENDING'", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("stock request tidak dapat diproses: status bukan PENDING atau sudah diproses sebelumnya")
	}
	return nil
}

func (r *inventoryRepository) CreateStockTransfer(ctx context.Context, tx *gorm.DB, transfer *domain.StockTransfer) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Create(transfer).Error
}

func (r *inventoryRepository) GetStockTransferByID(ctx context.Context, id uuid.UUID) (*domain.StockTransfer, error) {
	var transfer domain.StockTransfer
	err := r.db.WithContext(ctx).First(&transfer, "id = ?", id).Error
	return &transfer, err
}

func (r *inventoryRepository) GetStockTransfers(ctx context.Context, toWarehouse string) ([]domain.StockTransfer, error) {
	var transfers []domain.StockTransfer
	query := r.db.WithContext(ctx)
	if toWarehouse != "" {
		query = query.Where("to_warehouse = ?", toWarehouse)
	}
	err := query.Order("created_at desc").Find(&transfers).Error
	return transfers, err
}

func (r *inventoryRepository) UpdateStockTransferStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	res := execDB.WithContext(ctx).Model(&domain.StockTransfer{}).Where("id = ? AND status = 'IN_TRANSIT'", id).Update("status", status)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("transfer tidak dapat diproses: status bukan IN_TRANSIT")
	}
	return nil
}

func (r *inventoryRepository) CreateInventoryUsageLog(ctx context.Context, tx *gorm.DB, log *domain.InventoryUsageLog) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Create(log).Error
}

func (r *inventoryRepository) GetInventoryUsageLogByID(ctx context.Context, id uuid.UUID) (*domain.InventoryUsageLog, error) {
	var log domain.InventoryUsageLog
	err := r.db.WithContext(ctx).First(&log, "id = ?", id).Error
	return &log, err
}

func (r *inventoryRepository) UpdateInventoryUsageLogStatus(ctx context.Context, tx *gorm.DB, id uuid.UUID, status string) error {
	execDB := r.db
	if tx != nil {
		execDB = tx
	}
	return execDB.WithContext(ctx).Model(&domain.InventoryUsageLog{}).Where("id = ?", id).Update("status", status).Error
}

func (r *inventoryRepository) GetPendingUsageLogs(ctx context.Context) ([]domain.InventoryUsageLog, error) {
	var logs []domain.InventoryUsageLog
	err := r.db.WithContext(ctx).
		Where("status IN ?", []string{"RECORDED", "APPROVED", "PENDING_APPROVAL"}).
		Order("created_at desc").
		Find(&logs).Error
	return logs, err
}
