package domain

import (
	"github.com/google/uuid"
)

// MasterPuskesmas represents a Puskesmas in Depok
type MasterPuskesmas struct {
	ID        uuid.UUID `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Kecamatan string    `gorm:"type:varchar(100);not null" json:"kecamatan"`
}

// MasterKelurahan represents a Kelurahan mapped to a specific Puskesmas
type MasterKelurahan struct {
	ID          uuid.UUID       `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string          `gorm:"type:varchar(100);not null;unique" json:"name"`
	Kecamatan   string          `gorm:"type:varchar(100);not null" json:"kecamatan"`
	PuskesmasID uuid.UUID       `gorm:"type:uuid;not null" json:"puskesmas_id"`
	
	// Relations
	Puskesmas MasterPuskesmas `gorm:"foreignKey:PuskesmasID;references:ID" json:"puskesmas,omitempty"`
}
