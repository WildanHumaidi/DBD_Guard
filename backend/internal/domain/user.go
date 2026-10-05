package domain

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

const (
	RoleWarga            = "ROLE_WARGA"
	RoleKader            = "ROLE_KADER"
	RolePetugasPuskesmas = "ROLE_PETUGAS_PUSKESMAS"
	RoleKepalaPuskesmas  = "ROLE_KEPALA_PUSKESMAS"
	RoleTimFogging       = "ROLE_TIM_FOGGING"
	RoleDinkes           = "ROLE_DINKES"
	RoleAdminDeveloper   = "ROLE_ADMIN_DEVELOPER"
)

// Role represents the ROLES table
type Role struct {
	ID                uuid.UUID      `gorm:"column:role_id;type:uuid;default:gen_random_uuid();primaryKey" json:"role_id"`
	RoleName          string         `gorm:"type:varchar(100);not null" json:"role_name"`
	PermissionSet     datatypes.JSON `gorm:"type:jsonb" json:"permission_set"`
	AreaScopeRequired bool           `gorm:"not null;default:false" json:"area_scope_required"`
}

// User represents the USERS table
type User struct {
	ID           uuid.UUID `gorm:"column:user_id;type:uuid;default:gen_random_uuid();primaryKey" json:"user_id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Contact      string    `gorm:"type:varchar(50)" json:"contact"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:varchar(255);not null;default:''" json:"-"`
	RoleID       uuid.UUID  `gorm:"type:uuid;not null" json:"role_id"`
	LocationID   *uuid.UUID `gorm:"type:uuid" json:"location_id"`
	ActiveStatus bool       `gorm:"not null;default:true" json:"active_status"`

	// Relations
	Role Role `json:"role,omitempty"`
}
