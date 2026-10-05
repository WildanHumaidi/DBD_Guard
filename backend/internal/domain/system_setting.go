package domain

import (
	"time"

	"github.com/google/uuid"
)

type SystemSetting struct {
	ID                    uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	HeroTitle             string    `gorm:"type:varchar(255);not null"`
	AnnouncementText      string    `gorm:"type:text"`
	BannerImagePath       string    `gorm:"type:varchar(255)"`
	IsCitizenPortalActive bool      `gorm:"default:true"`
	UpdatedAt             time.Time
}
