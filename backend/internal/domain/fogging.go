package domain

import (
	"time"

	"github.com/google/uuid"
)

// FoggingRequest represents the FOGGING_REQUESTS table
type FoggingRequest struct {
	ID             uuid.UUID `gorm:"column:request_id;type:uuid;default:gen_random_uuid();primaryKey" json:"request_id"`
	PEID           uuid.UUID `gorm:"type:uuid;not null" json:"pe_id"`
	RiskScoreID    *uuid.UUID `gorm:"type:uuid" json:"risk_score_id"` // Can be nullable if based on other criteria
	ApprovalStatus string    `gorm:"type:varchar(50);not null" json:"approval_status"` // e.g., "PENDING", "APPROVED", "REJECTED"
	Reason         string    `gorm:"type:text" json:"reason"` // Mandatory for rejection
	ReservedStock  float64   `gorm:"type:double precision;not null;default:0" json:"reserved_stock"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	PE        PEResult         `json:"pe,omitempty"`
	RiskScore RiskScore        `json:"risk_score,omitempty"`
	Schedule  *FoggingSchedule `gorm:"foreignKey:RequestID" json:"schedule,omitempty"`
}

// FoggingSchedule represents the FOGGING_SCHEDULES table
type FoggingSchedule struct {
	ID              uuid.UUID `gorm:"column:schedule_id;type:uuid;default:gen_random_uuid();primaryKey" json:"schedule_id"`
	RequestID       uuid.UUID `gorm:"type:uuid;not null" json:"request_id"`
	TeamID          uuid.UUID `gorm:"type:uuid;not null" json:"team_id"` // Team responsible
	PlannedTime     time.Time `gorm:"not null" json:"planned_time"`
	ExecutionStatus string    `gorm:"type:varchar(50);not null" json:"execution_status"` // e.g., "SCHEDULED", "ON_THE_WAY", "IN_PROGRESS", "COMPLETED", "FAILED"
	Notes           string    `gorm:"type:text" json:"notes"` // E.g., failed due to rain
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Request FoggingRequest `json:"request,omitempty"`
	Team    User           `json:"team,omitempty"` // Assuming a Team is represented or linked via User/Group
	Evidence []FoggingEvidence `gorm:"foreignKey:ScheduleID" json:"evidence,omitempty"`
}

// FoggingEvidence represents the FOGGING_EVIDENCE table
type FoggingEvidence struct {
	ID         uuid.UUID `gorm:"column:evidence_id;type:uuid;default:gen_random_uuid();primaryKey" json:"evidence_id"`
	ScheduleID uuid.UUID `gorm:"type:uuid;not null" json:"schedule_id"`
	StorageKey string    `gorm:"type:varchar(255);not null" json:"storage_key"`
	UploadedBy uuid.UUID `gorm:"type:uuid;not null" json:"uploaded_by"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
