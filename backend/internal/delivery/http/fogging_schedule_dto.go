package http

import (
	"time"

	"github.com/google/uuid"
)

type CreateScheduleRequest struct {
	RequestID   uuid.UUID `json:"request_id" validate:"required"`
	TeamID      uuid.UUID `json:"team_id" validate:"required"`
	PlannedTime time.Time `json:"planned_time" validate:"required"`
}

type UsedInventoryItem struct {
	ItemID   uuid.UUID `json:"item_id" validate:"required"`
	Quantity float64   `json:"quantity" validate:"required,gt=0"`
}

type UpdateExecutionStatusRequest struct {
	ScheduleID      uuid.UUID           `params:"id" validate:"required"`
	ExecutionStatus string              `json:"execution_status" validate:"required,oneof=SCHEDULED ON_THE_WAY IN_PROGRESS COMPLETED FAILED"`
	Notes           string              `json:"notes"`
	UsedItems       []UsedInventoryItem `json:"used_items" validate:"dive"` // Only used when COMPLETED
}
