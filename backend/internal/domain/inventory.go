package domain

import (
	"time"

	"github.com/google/uuid"
)

// InventoryItem represents the INVENTORY_ITEMS table
type InventoryItem struct {
	ItemID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"item_id"`
	ItemType     string    `gorm:"type:varchar(100);not null" json:"item_type"` // e.g., "Insektisida", "Solar", "Alat Fogging"
	Quantity         float64   `gorm:"type:double precision;not null;default:0" json:"quantity"`
	ReservedQuantity float64   `gorm:"type:double precision;not null;default:0" json:"reserved_quantity"`
	Unit             string    `gorm:"type:varchar(50);not null" json:"unit"` // e.g., "Liter", "Botol", "Unit"
	MinimumStock float64   `gorm:"type:double precision;not null;default:0" json:"minimum_stock"`
	Warehouse    string    `gorm:"type:varchar(255);not null" json:"warehouse"` // Which Puskesmas warehouse
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// StockRequest represents a request for stock from Puskesmas to Dinkes
type StockRequest struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ItemID    uuid.UUID `gorm:"type:uuid;not null" json:"item_id"`
	Warehouse string    `gorm:"type:varchar(255);not null" json:"warehouse"` // The Puskesmas making the request
	Quantity  float64   `gorm:"type:double precision;not null" json:"quantity"`
	Status    string    `gorm:"type:varchar(50);not null;default:'PENDING'" json:"status"` // PENDING, APPROVED, REJECTED
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// StockTransfer represents a cross-allocation of stock between Puskesmas by Dinkes
type StockTransfer struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ItemID        uuid.UUID `gorm:"type:uuid;not null" json:"item_id"`
	FromWarehouse string    `gorm:"type:varchar(255);not null" json:"from_warehouse"`
	ToWarehouse   string    `gorm:"type:varchar(255);not null" json:"to_warehouse"`
	Quantity      float64   `gorm:"type:double precision;not null" json:"quantity"`
	Status        string    `gorm:"type:varchar(50);not null;default:'IN_TRANSIT'" json:"status"` // IN_TRANSIT, COMPLETED, REJECTED
	Notes         string    `gorm:"type:text" json:"notes"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// InventoryUsageLog represents the usage of stock by Fogging Teams
type InventoryUsageLog struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ItemID    uuid.UUID `gorm:"type:uuid;not null" json:"item_id"`
	TeamID    uuid.UUID `gorm:"type:uuid;not null" json:"team_id"`
	ScheduleID uuid.UUID `gorm:"type:uuid;not null" json:"schedule_id"`
	Quantity  float64   `gorm:"type:double precision;not null" json:"quantity"`
	Status    string    `gorm:"type:varchar(50);not null;default:'PENDING_APPROVAL'" json:"status"` // PENDING_APPROVAL, APPROVED, REJECTED
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
