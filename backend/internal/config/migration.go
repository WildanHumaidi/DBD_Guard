package config

import (
	"log"
	"os"

	"github.com/dbd-guard/backend/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// RunAutoMigration creates or updates the database schema based on the domain structs.
func RunAutoMigration() {
	if DB == nil {
		log.Println("Warning: Cannot run migration, database connection is nil.")
		return
	}

	// Enable PostGIS extension if it doesn't exist
	err := DB.Exec("CREATE EXTENSION IF NOT EXISTS postgis;").Error
	if err != nil {
		log.Printf("Warning: Failed to create PostGIS extension: %v. Ignore if PostGIS is already enabled.", err)
	}

	log.Println("Running AutoMigration...")

	err = DB.AutoMigrate(
		&domain.Role{},
		&domain.MasterPuskesmas{},
		&domain.MasterKelurahan{},
		&domain.User{},
		&domain.Reporter{},
		&domain.CitizenReport{},
		&domain.ReportLocation{},
		&domain.ReportPhoto{},
		&domain.Verification{},
		&domain.PEResult{},
		&domain.RiskScore{},
		&domain.FoggingRequest{},
		&domain.FoggingSchedule{},
		&domain.FoggingEvidence{},
		&domain.InventoryItem{},
		&domain.StockRequest{},
		&domain.StockTransfer{},
		&domain.InventoryUsageLog{},
		&domain.AuditLog{},
		&domain.KlbConfiguration{},
		&domain.KlbAlert{},
		&domain.SystemSetting{},
	)

	if err != nil {
		log.Fatalf("Error running auto migration: %v", err)
	}

	// Create GIST index on geom::geography for spatial performance < 200ms (Workstream E)
	err = DB.Exec("CREATE INDEX IF NOT EXISTS idx_report_locations_geom_geography ON report_locations USING GIST ((geom::geography));").Error
	if err != nil {
		log.Printf("Warning: Failed to create GIST geography index: %v", err)
	}

	seedAdmin()
	seedSystemSettings()

	log.Println("AutoMigration completed successfully.")
}


// MigrateDeveloperRole ensures RoleAdminDeveloper exists idempotently across fresh and existing databases,
// without creating duplicates or breaking existing user references.
func MigrateDeveloperRole(db *gorm.DB) (*domain.Role, error) {
	var devRole domain.Role
	err := db.Transaction(func(tx *gorm.DB) error {
		var legacyRole domain.Role
		errDev := tx.Where("role_name = ?", domain.RoleAdminDeveloper).First(&devRole).Error
		errLegacy := tx.Where("role_name = ?", "Admin").First(&legacyRole).Error

		if errLegacy == nil && errDev != nil {
			// Existing database with legacy role "Admin": rename to domain.RoleAdminDeveloper
			legacyRole.RoleName = domain.RoleAdminDeveloper
			legacyRole.PermissionSet = datatypes.JSON(`["all"]`)
			if err := tx.Save(&legacyRole).Error; err != nil {
				return err
			}
			devRole = legacyRole
		} else if errLegacy == nil && errDev == nil {
			// Both exist: reassign all users from legacy role to devRole, then delete legacy role
			if err := tx.Model(&domain.User{}).Where("role_id = ?", legacyRole.ID).Update("role_id", devRole.ID).Error; err != nil {
				return err
			}
			if err := tx.Delete(&legacyRole).Error; err != nil {
				return err
			}
		} else if errLegacy != nil && errDev != nil {
			// Fresh database: create domain.RoleAdminDeveloper
			devRole = domain.Role{
				RoleName:          domain.RoleAdminDeveloper,
				PermissionSet:     datatypes.JSON(`["all"]`),
				AreaScopeRequired: false,
			}
			if err := tx.Create(&devRole).Error; err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &devRole, nil
}

func seedAdmin() {
	// Get admin password from environment variable
	adminPassword := os.Getenv("ADMIN_DEFAULT_PASSWORD")
	if adminPassword == "" {
		adminPassword = "password123"
		log.Println("WARNING: ADMIN_DEFAULT_PASSWORD tidak di-set, menggunakan default. Ganti segera untuk production!")
	}

	// Hash password at runtime
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}
	hashedPassword := string(hashedBytes)

	devRole, err := MigrateDeveloperRole(DB)
	if err != nil {
		log.Printf("Failed to migrate developer role: %v", err)
		return
	}

	// Create Admin User
	admin := domain.User{
		Name:         "Super Admin",
		Email:        "admin@dbdguard.id",
		PasswordHash: hashedPassword,
		RoleID:       devRole.ID,
		ActiveStatus: true,
	}

	var existingAdmin domain.User
	if err := DB.Where("email = ?", admin.Email).First(&existingAdmin).Error; err != nil {
		if err := DB.Create(&admin).Error; err != nil {
			log.Printf("Failed to seed admin: %v", err)
		} else {
			log.Println("Successfully seeded default admin user.")
		}
	} else {
		existingAdmin.PasswordHash = hashedPassword
		existingAdmin.RoleID = devRole.ID
		if err := DB.Save(&existingAdmin).Error; err != nil {
			log.Printf("Failed to update admin user: %v", err)
		} else {
			log.Println("Successfully updated admin user from environment.")
		}
	}
}

func seedSystemSettings() {
	var count int64
	DB.Model(&domain.SystemSetting{}).Count(&count)
	if count == 0 {
		defaultSettings := domain.SystemSetting{
			HeroTitle:             "Sistem Pelaporan & Pemantauan DBD Kota Depok",
			AnnouncementText:      "Waspada Musim Hujan! Segera bersihkan sarang nyamuk dengan 3M Plus.",
			IsCitizenPortalActive: true,
		}
		if err := DB.Create(&defaultSettings).Error; err != nil {
			log.Printf("Failed to seed system settings: %v", err)
		} else {
			log.Println("Successfully seeded default system settings.")
		}
	}
}
