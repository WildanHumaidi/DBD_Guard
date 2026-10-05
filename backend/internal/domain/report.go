package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrUnauthorizedJurisdiction returned when an actor attempts to access/claim a report outside their assigned jurisdiction.
	ErrUnauthorizedJurisdiction = errors.New("akses ditolak: laporan di luar yurisdiksi wilayah kerja")
	// ErrMissingLocationContext returned when an actor lacks a valid location_id context required for territorial authorization.
	ErrMissingLocationContext = errors.New("akses ditolak: akun tidak memiliki wilayah kerja yang valid")
)

// Reporter represents the REPORTER table
type Reporter struct {
	ID            uuid.UUID `gorm:"column:reporter_id;type:uuid;default:gen_random_uuid();primaryKey" json:"reporter_id"`
	Name          string    `gorm:"type:varchar(255);not null" json:"name"`
	Contact       string    `gorm:"type:varchar(50)" json:"contact"`
	AnonymousFlag bool      `gorm:"not null;default:false" json:"anonymous_flag"`
}

// CitizenReport represents the CITIZEN_REPORTS table
type CitizenReport struct {
	ID          uuid.UUID `gorm:"column:report_id;type:uuid;default:gen_random_uuid();primaryKey" json:"report_id"`
	ReferenceNo string    `gorm:"type:varchar(50);unique;not null" json:"reference_no"`
	Category    string    `gorm:"type:varchar(100);not null" json:"category"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:varchar(50);not null" json:"status"`
	ReporterID          uuid.UUID  `gorm:"type:uuid;not null" json:"reporter_id"`
	InvestigatorID      *uuid.UUID `gorm:"type:uuid" json:"investigator_id"`
	IsPossibleDuplicate bool       `gorm:"not null;default:false" json:"is_possible_duplicate"`
	MergedWithID        *uuid.UUID `gorm:"type:uuid" json:"merged_with_id,omitempty"`
	FeedbackRating      int       `gorm:"type:int" json:"feedback_rating,omitempty"`
	FeedbackNotes       string    `gorm:"type:text" json:"feedback_notes,omitempty"`
	FeedbackRequested   bool      `gorm:"not null;default:false" json:"feedback_requested"`
	AppealCount         int       `gorm:"type:int;not null;default:0" json:"appeal_count"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Reporter Reporter       `json:"reporter,omitempty"`
	Location ReportLocation `gorm:"foreignKey:ReportID" json:"location,omitempty"`
	Photos   []ReportPhoto  `gorm:"foreignKey:ReportID" json:"photos,omitempty"`
	PEResult *PEResult      `gorm:"foreignKey:ReportID" json:"pe_result,omitempty"`
}

// ReportLocation represents the REPORT_LOCATIONS table
type ReportLocation struct {
	ID         uuid.UUID `gorm:"column:location_id;type:uuid;default:gen_random_uuid();primaryKey" json:"location_id"`
	ReportID   uuid.UUID `gorm:"type:uuid;not null" json:"report_id"`
	Latitude   float64   `gorm:"type:double precision;not null" json:"latitude"`
	Longitude  float64   `gorm:"type:double precision;not null" json:"longitude"`
	Geom       string    `gorm:"type:geometry(Point,4326)" json:"geom"` // PostGIS geometry
	Address    string    `gorm:"type:text" json:"address"`
	RT         string    `gorm:"type:varchar(10)" json:"rt"`
	RW         string    `gorm:"type:varchar(10)" json:"rw"`
	Kelurahan  string    `gorm:"type:varchar(100)" json:"kelurahan"`
	Kecamatan  string    `gorm:"type:varchar(100)" json:"kecamatan"`
}

// ReportPhoto represents the REPORT_PHOTOS table
type ReportPhoto struct {
	ID         uuid.UUID `gorm:"column:photo_id;type:uuid;default:gen_random_uuid();primaryKey" json:"photo_id"`
	ReportID   uuid.UUID `gorm:"type:uuid;not null" json:"report_id"`
	StorageKey string    `gorm:"type:varchar(255);not null" json:"storage_key"`
	Hash       string    `gorm:"type:varchar(255)" json:"hash"`
}
