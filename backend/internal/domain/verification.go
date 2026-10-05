package domain

import (
	"time"

	"github.com/google/uuid"
)

// Verification represents the VERIFICATIONS table
type Verification struct {
	ID             uuid.UUID `gorm:"column:verification_id;type:uuid;default:gen_random_uuid();primaryKey" json:"verification_id"`
	ReportID       uuid.UUID `gorm:"type:uuid;not null" json:"report_id"`
	ActorID        uuid.UUID `gorm:"type:uuid;not null" json:"actor_id"` // User who verified (Puskesmas)
	Decision       string    `gorm:"type:varchar(50);not null" json:"decision"` // "APPROVED", "REJECTED"
	Reason         string    `gorm:"type:text" json:"reason"` // Reason is mandatory for rejection per business rules
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Report CitizenReport `json:"report,omitempty"`
	Actor  User          `json:"actor,omitempty"`
}
