package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/dbd-guard/backend/internal/config"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func main() {
	fmt.Println("Starting Comprehensive Database Seeder...")

	// Load Configuration
	config.LoadConfig()

	// Connect to Database
	config.ConnectDB()
	db := config.DB
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}

	fmt.Println("Resetting database (Dropping all tables)...")
	// Drop all tables
	db.Migrator().DropTable(
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
		&domain.FoggingRequest{},
		&domain.InventoryItem{},
		&domain.StockRequest{},
		&domain.StockTransfer{},
		&domain.FoggingSchedule{},
		&domain.AuditLog{},
		&domain.RiskScore{},
	)

	// Auto Migrate All Domains
	fmt.Println("Running AutoMigrate...")
	err := db.AutoMigrate(
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
		&domain.FoggingRequest{},
		&domain.InventoryItem{},
		&domain.StockRequest{},
		&domain.StockTransfer{},
		&domain.FoggingSchedule{},
		&domain.AuditLog{},
		&domain.RiskScore{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// 1. Seed Roles
	fmt.Println("Seeding Roles...")
	roles := []domain.Role{
		{RoleName: domain.RoleAdminDeveloper, PermissionSet: datatypes.JSON(`["all"]`), AreaScopeRequired: false},
		{RoleName: domain.RoleDinkes, PermissionSet: datatypes.JSON(`["read_all", "approve_stock"]`), AreaScopeRequired: false},
		{RoleName: domain.RoleKepalaPuskesmas, PermissionSet: datatypes.JSON(`["verify", "pe", "fogging", "inventory"]`), AreaScopeRequired: true},
		{RoleName: domain.RolePetugasPuskesmas, PermissionSet: datatypes.JSON(`["verify", "pe", "inventory"]`), AreaScopeRequired: true},
		{RoleName: domain.RoleTimFogging, PermissionSet: datatypes.JSON(`["fogging"]`), AreaScopeRequired: true},
		{RoleName: domain.RoleKader, PermissionSet: datatypes.JSON(`["verify"]`), AreaScopeRequired: true},
		{RoleName: domain.RoleWarga, PermissionSet: datatypes.JSON(`["report"]`), AreaScopeRequired: false},
	}
	for i := range roles {
		db.Create(&roles[i])
	}
	roleMap := make(map[string]uuid.UUID)
	for _, r := range roles {
		roleMap[r.RoleName] = r.ID
	}

	// 2. Seed Master Data (Kecamatan, Puskesmas and Kelurahan)
	fmt.Println("Seeding Master Data (Kecamatan, Puskesmas, Kelurahan)...")
	kecamatans := []string{"Beji", "Bojongsari", "Cilodong", "Cimanggis", "Cinere", "Cipayung", "Limo", "Pancoran Mas", "Sawangan", "Sukmajaya", "Tapos"}

	var puskList []domain.MasterPuskesmas
	// 38 Puskesmas
	for i := 1; i <= 38; i++ {
		kec := kecamatans[i%len(kecamatans)]
		p := domain.MasterPuskesmas{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("Puskesmas %s %d", kec, i),
			Kecamatan: kec,
		}
		db.Create(&p)
		puskList = append(puskList, p)
	}

	var kelList []domain.MasterKelurahan
	// 63 Kelurahan
	for i := 1; i <= 63; i++ {
		pusk := puskList[i%len(puskList)] // assign evenly
		k := domain.MasterKelurahan{
			Name:        fmt.Sprintf("Kelurahan %d", i),
			Kecamatan:   pusk.Kecamatan,
			PuskesmasID: pusk.ID,
		}
		db.Create(&k)
		kelList = append(kelList, k)
	}

	// 3. Seed Users
	fmt.Println("Seeding Users (Admin, Dinkes, Puskesmas, Kelurahan)...")
	defaultHash := "$2a$10$t9asZm/hciyqysoaKyd68.te9Zr2e5Cy7L2JBeUsDxvDV0bn50HYu" // password123

	firstPuskID := puskList[0].ID

	// Primary Demo Users for each Role
	primaryDemoUsers := []domain.User{
		{Name: "Super Admin Developer", Email: "dev@dbdguard.id", Contact: "08111111111", RoleID: roleMap[domain.RoleAdminDeveloper], ActiveStatus: true, PasswordHash: defaultHash},
		{Name: "Super Admin Developer (Alt)", Email: "admin@dbdguard.id", Contact: "08111111118", RoleID: roleMap[domain.RoleAdminDeveloper], ActiveStatus: true, PasswordHash: defaultHash},
		{Name: "Dinas Kesehatan Kota Depok", Email: "admin@depok.go.id", Contact: "08111111112", RoleID: roleMap[domain.RoleDinkes], ActiveStatus: true, PasswordHash: defaultHash},
		{Name: "Kepala Puskesmas Beji", Email: "kepala.puskesmas@depok.go.id", Contact: "08111111113", RoleID: roleMap[domain.RoleKepalaPuskesmas], LocationID: &firstPuskID, ActiveStatus: true, PasswordHash: defaultHash},
		{Name: "Petugas PE Puskesmas Beji", Email: "petugas.puskesmas@depok.go.id", Contact: "08111111114", RoleID: roleMap[domain.RolePetugasPuskesmas], LocationID: &firstPuskID, ActiveStatus: true, PasswordHash: defaultHash},
		{Name: "Tim Fogging Puskesmas Beji", Email: "tim.fogging@depok.go.id", Contact: "08111111115", RoleID: roleMap[domain.RoleTimFogging], LocationID: &firstPuskID, ActiveStatus: true, PasswordHash: defaultHash},
		{Name: "Kader Wilayah Kelurahan Beji", Email: "kader@depok.go.id", Contact: "08111111116", RoleID: roleMap[domain.RoleKader], LocationID: &firstPuskID, ActiveStatus: true, PasswordHash: defaultHash},
		{Name: "Warga Pelapor Depok", Email: "warga@gmail.com", Contact: "08111111117", RoleID: roleMap[domain.RoleWarga], ActiveStatus: true, PasswordHash: defaultHash},
	}
	for i := range primaryDemoUsers {
		db.Create(&primaryDemoUsers[i])
	}

	// Users for each Puskesmas (Kepala, Petugas, Tim Fogging)
	var allPuskUsers []domain.User
	for idx, pusk := range puskList {
		puskID := pusk.ID
		users := []domain.User{
			{Name: fmt.Sprintf("Kepala %s", pusk.Name), Email: fmt.Sprintf("kepala.pusk%d@depok.go.id", idx+1), Contact: fmt.Sprintf("082000%05d", rand.Intn(99999)), RoleID: roleMap[domain.RoleKepalaPuskesmas], LocationID: &puskID, ActiveStatus: true, PasswordHash: defaultHash},
			{Name: fmt.Sprintf("Petugas PE %s", pusk.Name), Email: fmt.Sprintf("petugas.pusk%d@depok.go.id", idx+1), Contact: fmt.Sprintf("082100%05d", rand.Intn(99999)), RoleID: roleMap[domain.RolePetugasPuskesmas], LocationID: &puskID, ActiveStatus: true, PasswordHash: defaultHash},
			{Name: fmt.Sprintf("Tim Fogging %s", pusk.Name), Email: fmt.Sprintf("fogging.pusk%d@depok.go.id", idx+1), Contact: fmt.Sprintf("082200%05d", rand.Intn(99999)), RoleID: roleMap[domain.RoleTimFogging], LocationID: &puskID, ActiveStatus: true, PasswordHash: defaultHash},
		}
		for i := range users {
			db.Create(&users[i])
			allPuskUsers = append(allPuskUsers, users[i])
		}
	}

	// Users for each Kelurahan (Kader / RT/RW)
	var allKaderUsers []domain.User
	for idx, kel := range kelList {
		puskID := kel.PuskesmasID
		u := domain.User{
			Name:         fmt.Sprintf("Kader RT/RW %s", kel.Name),
			Email:        fmt.Sprintf("kader.kel%d@dbdguard.id", idx+1),
			Contact:      fmt.Sprintf("083000%05d", rand.Intn(99999)),
			RoleID:       roleMap[domain.RoleKader],
			LocationID:   &puskID,
			ActiveStatus: true,
			PasswordHash: defaultHash,
		}
		db.Create(&u)
		allKaderUsers = append(allKaderUsers, u)
	}

	// Some Warga Users
	for i := 1; i <= 10; i++ {
		warga := domain.User{
			Name:         fmt.Sprintf("Warga %d", i),
			Email:        fmt.Sprintf("warga%d@gmail.com", i),
			Contact:      fmt.Sprintf("084000%d", rand.Intn(99999)),
			RoleID:       roleMap[domain.RoleWarga],
			ActiveStatus: true,
			PasswordHash: defaultHash,
		}
		db.Create(&warga)
	}

	// 4. Seed Inventory
	fmt.Println("Seeding Inventories for all Puskesmas and Dinkes...")
	// Dinkes / Gudang Utama
	dinkesItems := []domain.InventoryItem{
		{ItemType: "ABATE", Quantity: 5000, MinimumStock: 1000, Unit: "Kg", Warehouse: "Gudang Utama Dinkes"},
		{ItemType: "INSEKTISIDA", Quantity: 10000, MinimumStock: 2000, Unit: "Liter", Warehouse: "Gudang Utama Dinkes"},
		{ItemType: "ALAT FOGGING", Quantity: 500, MinimumStock: 100, Unit: "Unit", Warehouse: "Gudang Utama Dinkes"},
		{ItemType: "SOLAR", Quantity: 15000, MinimumStock: 3000, Unit: "Liter", Warehouse: "Gudang Utama Dinkes"},
	}
	for i := range dinkesItems {
		db.Create(&dinkesItems[i])
	}

	// Puskesmas Inventories
	for _, pusk := range puskList {
		items := []domain.InventoryItem{
			{ItemType: "ABATE", Quantity: float64(50 + rand.Intn(100)), MinimumStock: 20, Unit: "Kg", Warehouse: pusk.Name},
			{ItemType: "INSEKTISIDA", Quantity: float64(30 + rand.Intn(200)), MinimumStock: 50, Unit: "Liter", Warehouse: pusk.Name},
			{ItemType: "ALAT FOGGING", Quantity: float64(5 + rand.Intn(10)), MinimumStock: 2, Unit: "Unit", Warehouse: pusk.Name},
		}
		for i := range items {
			db.Create(&items[i])
		}
	}

	// 5. Seed Reports & Scenarios
	fmt.Println("Seeding Citizen Reports & Scenarios...")
	// We'll create more dummy reports spread out
	for i := 1; i <= 100; i++ {
		kel := kelList[rand.Intn(len(kelList))]
		
		reporter := domain.Reporter{
			Name:          fmt.Sprintf("Pelapor %d", i),
			Contact:       fmt.Sprintf("0899%d", rand.Intn(99999)),
			AnonymousFlag: i%5 == 0,
		}
		db.Create(&reporter)

		status := "MENUNGGU_PE_KADER"
		var rating int
		var notes string

		// Distribute statuses realistically
		if i%10 == 0 {
			status = "REJECTED"
		} else if i%2 == 0 {
			status = "COMPLETED"
			rating = rand.Intn(3) + 3 // Rating 3 to 5
			notes = "Pelayanan sangat baik dan cepat tanggap."
		} else if i%3 == 0 {
			status = "FOGGING_SCHEDULED"
		} else if i%5 == 0 {
			status = "WAITING_FOGGING_APPROVAL"
		} else if i%7 == 0 {
			status = "PE_SCHEDULED"
		} else {
			status = "VERIFIED"
		}

		report := domain.CitizenReport{
			ReferenceNo:    fmt.Sprintf("DBD-DPK-%d%d", time.Now().Year(), 1000+i),
			Category:       "Kasus DBD",
			Description:    fmt.Sprintf("Laporan kasus DBD dari warga di kelurahan %s", kel.Name),
			Status:         status,
			ReporterID:     reporter.ID,
			FeedbackRating: rating,
			FeedbackNotes:  notes,
		}
		db.Create(&report)

		lat := -6.400000 + (rand.Float64() - 0.5) * 0.1
		lng := 106.820000 + (rand.Float64() - 0.5) * 0.1
		
		loc := domain.ReportLocation{
			ReportID:  report.ID,
			Latitude:  lat,
			Longitude: lng,
			Address:   fmt.Sprintf("Jl. Dummy No.%d, %s", i, kel.Name),
			Kelurahan: kel.Name,
			Kecamatan: kel.Kecamatan,
		}
		db.Omit("Geom").Create(&loc)

		// Create PE & Fogging for specific statuses
		if status == "WAITING_FOGGING_APPROVAL" || status == "FOGGING_SCHEDULED" || status == "COMPLETED" {
			pe := domain.PEResult{
				ReportID:       report.ID,
				OfficerID:      allPuskUsers[1].ID, // assign a random officer id (using first for simplicity)
				Recommendation: "FOGGING_FOKUS",
				LarvaeFound:    true,
				FeverCount:     rand.Intn(5) + 1,
				InspectionDate: time.Now().Add(-24 * time.Hour),
			}
			db.Create(&pe)

			appStat := "PENDING"
			if status == "FOGGING_SCHEDULED" || status == "COMPLETED" {
				appStat = "APPROVED"
			}
			
			fogReq := domain.FoggingRequest{
				PEID:           pe.ID,
				ApprovalStatus: appStat,
				ReservedStock:  10,
			}
			db.Create(&fogReq)

			if status == "FOGGING_SCHEDULED" || status == "COMPLETED" {
				execStat := "SCHEDULED"
				if status == "COMPLETED" {
					execStat = "COMPLETED"
				}
				sched := domain.FoggingSchedule{
					RequestID:       fogReq.ID,
					TeamID:          allPuskUsers[2].ID, // assign random team
					PlannedTime:     time.Now().Add(24 * time.Hour),
					ExecutionStatus: execStat,
				}
				db.Create(&sched)
			}
		}
	}

	fmt.Println("Seeding finished successfully! All 63 Kelurahans, 38 Puskesmas, Inventories, and Users per role are created.")
}
