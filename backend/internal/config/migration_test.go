package config_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dbd-guard/backend/internal/config"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getTestDB(t *testing.T) *gorm.DB {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5433"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}
	pass := os.Getenv("DB_PASSWORD")
	if pass == "" {
		pass = "secret"
	}
	name := os.Getenv("DB_NAME")
	if name == "" {
		name = "dbd_guard"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", host, user, pass, name, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("Skipping DB test: database not reachable: %v", err)
	}
	return db
}

func TestMigrateDeveloperRole_FreshAndExisting(t *testing.T) {
	db := getTestDB(t)

	// Use a transaction rolled back at the end of the test to avoid polluting the DB
	tx := db.Begin()
	defer tx.Rollback()

	// 1. Test Idempotency on existing database
	devRole, err := config.MigrateDeveloperRole(tx)
	if err != nil {
		t.Fatalf("MigrateDeveloperRole failed on existing database: %v", err)
	}
	if devRole.RoleName != domain.RoleAdminDeveloper {
		t.Errorf("expected role name %s, got %s", domain.RoleAdminDeveloper, devRole.RoleName)
	}

	var countAdmin int64
	tx.Model(&domain.Role{}).Where("role_name = 'Admin'").Count(&countAdmin)
	if countAdmin != 0 {
		t.Errorf("expected 0 legacy Admin roles, found %d", countAdmin)
	}

	// 2. Scenario: Existing database has both "Admin" and "ROLE_ADMIN_DEVELOPER"
	// Create legacy "Admin" role and assign a test user to it
	legacyRole := domain.Role{
		ID:                uuid.New(),
		RoleName:          "Admin",
		PermissionSet:     datatypes.JSON(`["all"]`),
		AreaScopeRequired: false,
	}
	if err := tx.Create(&legacyRole).Error; err != nil {
		t.Fatalf("failed to create legacy Admin role: %v", err)
	}

	testUser := domain.User{
		ID:           uuid.New(),
		Name:         "User Under Legacy Admin",
		Email:        "legacy_admin_user@test.com",
		RoleID:       legacyRole.ID,
		ActiveStatus: true,
	}
	if err := tx.Create(&testUser).Error; err != nil {
		t.Fatalf("failed to create test user under legacy Admin: %v", err)
	}

	// Now run MigrateDeveloperRole: it should detect dual roles, reassign user to devRole, and remove "Admin"
	migratedRole, err := config.MigrateDeveloperRole(tx)
	if err != nil {
		t.Fatalf("MigrateDeveloperRole failed on dual role scenario: %v", err)
	}
	if migratedRole.RoleName != domain.RoleAdminDeveloper {
		t.Errorf("expected role name %s, got %s", domain.RoleAdminDeveloper, migratedRole.RoleName)
	}

	// Verify legacy role "Admin" was deleted
	tx.Model(&domain.Role{}).Where("role_name = 'Admin'").Count(&countAdmin)
	if countAdmin != 0 {
		t.Errorf("expected 0 legacy Admin roles after dual migration, found %d", countAdmin)
	}

	// Verify test user was reassigned to devRole without foreign key violation
	var updatedUser domain.User
	if err := tx.Preload("Role").Where("email = ?", "legacy_admin_user@test.com").First(&updatedUser).Error; err != nil {
		t.Fatalf("failed to fetch user after migration: %v", err)
	}
	if updatedUser.RoleID != devRole.ID {
		t.Errorf("expected user role ID %s, got %s", devRole.ID, updatedUser.RoleID)
	}
	if updatedUser.Role.RoleName != domain.RoleAdminDeveloper {
		t.Errorf("expected user role name %s, got %s", domain.RoleAdminDeveloper, updatedUser.Role.RoleName)
	}
}

func TestAutoMigrateForeignKeysAndRelationships(t *testing.T) {
	db := getTestDB(t)

	tx := db.Begin()
	defer tx.Rollback()

	// Verify relationships can be queried without foreign key violation
	var puskesmasList []domain.MasterPuskesmas
	if err := tx.Find(&puskesmasList).Error; err != nil {
		t.Fatalf("failed to query puskesmas: %v", err)
	}

	var kelurahanList []domain.MasterKelurahan
	if err := tx.Preload("Puskesmas").Find(&kelurahanList).Error; err != nil {
		t.Fatalf("failed to query kelurahan with puskesmas relation: %v", err)
	}

	var reports []domain.CitizenReport
	if err := tx.Preload("Reporter").Preload("Location").Limit(5).Find(&reports).Error; err != nil {
		t.Fatalf("failed to preload citizen report relations: %v", err)
	}
}
