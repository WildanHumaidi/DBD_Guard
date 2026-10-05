package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// AuditLog represents the AUDIT_LOGS table (Business Rule: 6)
type AuditLog struct {
	ID         uuid.UUID      `gorm:"column:log_id;type:uuid;default:gen_random_uuid();primaryKey" json:"log_id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"` // Who performed the action
	Action     string         `gorm:"type:varchar(100);not null" json:"action"` // e.g., "UPDATE_STATUS", "APPROVE_FOGGING"
	EntityType string         `gorm:"type:varchar(100);not null" json:"entity_type"` // e.g., "CITIZEN_REPORT", "FOGGING_REQUEST"
	EntityID   uuid.UUID      `gorm:"type:uuid;not null" json:"entity_id"`
	Changes    datatypes.JSON `gorm:"type:jsonb" json:"changes"` // What changed (before/after)
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	User User `json:"user,omitempty"`
}
