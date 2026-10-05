package http

type CreateInventoryItemRequest struct {
	ItemType     string  `json:"item_type" validate:"required"`
	Quantity     float64 `json:"quantity" validate:"required,min=0"`
	Unit         string  `json:"unit" validate:"required"`
	MinimumStock float64 `json:"minimum_stock" validate:"required,min=0"`
	Warehouse    string  `json:"warehouse" validate:"required"`
}

// UpdateStockRequest represents the payload for adjusting inventory stock.
type UpdateStockRequest struct {
	Quantity float64 `json:"quantity" validate:"required"` // positive to add, negative to subtract
	Reason   string  `json:"reason" validate:"required"`
}

type CreateStockRequestRequest struct {
	ItemID    string  `json:"item_id" validate:"required,uuid"`
	Warehouse string  `json:"warehouse" validate:"required"`
	Quantity  float64 `json:"quantity" validate:"required,min=0.1"`
	Notes     string  `json:"notes"`
}

type ApproveStockRequestRequest struct {
	Status string `json:"status" validate:"required,oneof=APPROVED REJECTED"`
}

type CreateStockTransferRequest struct {
	ItemID        string  `json:"item_id" validate:"required,uuid"`
	FromWarehouse string  `json:"from_warehouse" validate:"required"`
	ToWarehouse   string  `json:"to_warehouse" validate:"required"`
	Quantity      float64 `json:"quantity" validate:"required,min=0.1"`
	Notes         string  `json:"notes"`
}
