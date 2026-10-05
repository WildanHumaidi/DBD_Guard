package domain

import (
	"time"

	"github.com/google/uuid"
)

// KlbConfiguration represents the KLB_CONFIGURATIONS table
type KlbConfiguration struct {
	ID             uuid.UUID `gorm:"column:config_id;type:uuid;default:gen_random_uuid();primaryKey" json:"config_id"`
	RegionScope    string    `gorm:"type:varchar(50);not null;default:'KELURAHAN'" json:"region_scope"` // 'KELURAHAN' or 'KECAMATAN'
	MinCases       int       `gorm:"type:int;not null;default:5" json:"min_cases"`
	TimeframeDays  int       `gorm:"type:int;not null;default:7" json:"timeframe_days"`
	IsAutoTrigger  bool      `gorm:"not null;default:true" json:"is_auto_trigger"`
	LastUpdatedBy  uuid.UUID `gorm:"type:uuid" json:"last_updated_by,omitempty"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// KlbAlert represents the KLB_ALERTS table
type KlbAlert struct {
	ID             uuid.UUID `gorm:"column:alert_id;type:uuid;default:gen_random_uuid();primaryKey" json:"alert_id"`
	AreaName       string    `gorm:"type:varchar(100);not null" json:"area_name"` // Name of Kelurahan/Kecamatan
	AreaType       string    `gorm:"type:varchar(50);not null" json:"area_type"` // 'KELURAHAN' or 'KECAMATAN'
	CasesCount     int       `gorm:"type:int;not null" json:"cases_count"`
	TriggeredAt    time.Time `gorm:"autoCreateTime" json:"triggered_at"`
	ResolvedAt     *time.Time `gorm:"type:timestamp" json:"resolved_at,omitempty"`
	IsActive       bool      `gorm:"not null;default:true" json:"is_active"`
}
