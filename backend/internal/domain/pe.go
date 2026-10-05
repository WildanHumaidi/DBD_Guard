package domain

import (
	"time"

	"github.com/google/uuid"
)

// PEResult represents the PE_RESULTS table (Penyelidikan Epidemiologi)
type PEResult struct {
	ID             uuid.UUID  `gorm:"column:pe_id;type:uuid;default:gen_random_uuid();primaryKey" json:"pe_id"`
	ReportID       uuid.UUID  `gorm:"type:uuid;not null" json:"report_id"`
	OfficerID      uuid.UUID  `gorm:"type:uuid;not null" json:"officer_id"` // User who conducted PE (Kader)
	ValidatorID    *uuid.UUID `gorm:"type:uuid" json:"validator_id"` // Puskesmas officer who validated
	Status         string     `gorm:"type:varchar(50);not null;default:'PENDING_VALIDATION'" json:"status"` // PENDING_VALIDATION, VALIDATED, REJECTED
	InspectionDate time.Time  `gorm:"type:date;not null" json:"inspection_date"`
	RadiusMeters   int        `gorm:"not null;default:0" json:"pe_radius"`
	LarvaeFound    bool       `gorm:"not null;default:false" json:"larva_found"`
	FeverCount     int        `gorm:"not null;default:0" json:"fever_count"`
	Recommendation string     `gorm:"type:text" json:"recommendation"`
	PsnEdukasiDone bool       `gorm:"not null;default:false" json:"psn_edukasi_done"`
	Notes          string     `gorm:"type:text" json:"notes"`
	EvidencePhotoURL string   `gorm:"type:varchar(255)" json:"evidence_photo_url"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Report         CitizenReport   `json:"report,omitempty"`
	Officer        User            `gorm:"foreignKey:OfficerID" json:"officer,omitempty"`
	Validator      *User           `gorm:"foreignKey:ValidatorID" json:"validator,omitempty"`
	FoggingRequest *FoggingRequest `gorm:"foreignKey:PEID" json:"fogging_request,omitempty"`
}

// RiskScore represents the RISK_SCORES table
type RiskScore struct {
	ID           uuid.UUID `gorm:"column:score_id;type:uuid;default:gen_random_uuid();primaryKey" json:"score_id"`
	AreaID       string    `gorm:"type:varchar(100);not null" json:"area_id"` // Can be Kelurahan ID
	Score        float64   `gorm:"type:double precision;not null" json:"score"`
	RiskLevel    string    `gorm:"type:varchar(50);not null" json:"risk_level"` // e.g., "TINGGI", "SEDANG", "RENDAH"
	CalculatedAt time.Time `gorm:"autoCreateTime" json:"calculated_at"`
}
