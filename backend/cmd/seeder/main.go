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
	fmt.Println("================================================================")
	fmt.Println("🚀 Menjalankan Database Seeder Idempoten - DBD Guard (Kota Depok)")
	fmt.Println("   Skenario: Laporan -> Verifikasi -> PE -> Rekomendasi -> Fogging")
	fmt.Println("   Label: DATA SIMULASI (100% Fiktif untuk PKL Diskominfo)")
	fmt.Println("================================================================")

	// 1. Load Configuration & Safety Checks
	config.LoadConfig()

	// Pengaman Lingkungan Seeder: Hanya boleh dijalankan di environment 'development' dan host localhost
	if config.Env.AppEnv != "development" {
		log.Fatalf("PENGAMAN SEEDER DITOLAK: Seeder hanya diizinkan berjalan pada APP_ENV='development' (terdeteksi: APP_ENV='%s')", config.Env.AppEnv)
	}

	allowedHosts := map[string]bool{
		"localhost": true,
		"127.0.0.1": true,
		"::1":       true,
		"":          true,
	}
	if !allowedHosts[config.Env.DBHost] {
		log.Fatalf("PENGAMAN SEEDER DITOLAK: Seeder hanya diizinkan berjalan pada host database lokal ('localhost' atau '127.0.0.1', terdeteksi DB_HOST='%s')", config.Env.DBHost)
	}

	config.ConnectDB()
	db := config.DB
	if db == nil {
		log.Fatal("Gagal menghubungkan ke database PostgreSQL/PostGIS")
	}

	fmt.Println("\n[1/7] Mereset tabel database secara bersih...")
	// Drop tables in dependency order
	db.Migrator().DropTable(
		&domain.FoggingEvidence{},
		&domain.FoggingSchedule{},
		&domain.StockTransfer{},
		&domain.StockRequest{},
		&domain.InventoryItem{},
		&domain.FoggingRequest{},
		&domain.PEResult{},
		&domain.Verification{},
		&domain.ReportPhoto{},
		&domain.ReportLocation{},
		&domain.CitizenReport{},
		&domain.Reporter{},
		&domain.AuditLog{},
		&domain.RiskScore{},
		&domain.KlbAlert{},
		&domain.KlbConfiguration{},
		&domain.User{},
		&domain.MasterKelurahan{},
		&domain.MasterPuskesmas{},
		&domain.Role{},
	)

	fmt.Println("[2/7] Menjalankan AutoMigrate seluruh tabel domain...")
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
		&domain.FoggingEvidence{},
		&domain.AuditLog{},
		&domain.RiskScore{},
		&domain.KlbConfiguration{},
		&domain.KlbAlert{},
	)
	if err != nil {
		log.Fatalf("AutoMigrate gagal: %v", err)
	}

	// 2. Seed Roles
	fmt.Println("[3/7] Membuat 7 Peran Pengguna (RBAC)...")
	roles := []domain.Role{
		{RoleName: domain.RoleAdminDeveloper, PermissionSet: datatypes.JSON(`["all"]`), AreaScopeRequired: false},
		{RoleName: domain.RoleDinkes, PermissionSet: datatypes.JSON(`["read_all", "approve_stock"]`), AreaScopeRequired: false},
		{RoleName: domain.RoleKepalaPuskesmas, PermissionSet: datatypes.JSON(`["verify", "pe", "fogging", "inventory"]`), AreaScopeRequired: true},
		{RoleName: domain.RolePetugasPuskesmas, PermissionSet: datatypes.JSON(`["verify", "pe", "inventory"]`), AreaScopeRequired: true},
		{RoleName: domain.RoleTimFogging, PermissionSet: datatypes.JSON(`["fogging"]`), AreaScopeRequired: true},
		{RoleName: domain.RoleKader, PermissionSet: datatypes.JSON(`["verify"]`), AreaScopeRequired: true},
		{RoleName: domain.RoleWarga, PermissionSet: datatypes.JSON(`["report"]`), AreaScopeRequired: false},
	}
	roleMap := make(map[string]uuid.UUID)
	for i := range roles {
		db.Create(&roles[i])
		roleMap[roles[i].RoleName] = roles[i].ID
	}

	// 3. Seed Master Data (11 Kecamatan, 38 Puskesmas, 63 Kelurahan Kota Depok)
	fmt.Println("[4/7] Memuat Master Data Wilayah Kota Depok...")
	kecamatans := []string{
		"Beji", "Bojongsari", "Cilodong", "Cimanggis", "Cinere", 
		"Cipayung", "Limo", "Pancoran Mas", "Sawangan", "Sukmajaya", "Tapos",
	}

	// Create 38 Puskesmas
	var puskList []domain.MasterPuskesmas
	puskBejiID := uuid.New()
	puskBeji := domain.MasterPuskesmas{
		ID:        puskBejiID,
		Name:      "Puskesmas Beji",
		Kecamatan: "Beji",
	}
	db.Create(&puskBeji)
	puskList = append(puskList, puskBeji)

	puskCipayungID := uuid.New()
	puskCipayung := domain.MasterPuskesmas{
		ID:        puskCipayungID,
		Name:      "Puskesmas Cipayung",
		Kecamatan: "Cipayung",
	}
	db.Create(&puskCipayung)
	puskList = append(puskList, puskCipayung)

	for i := 3; i <= 38; i++ {
		kec := kecamatans[i%len(kecamatans)]
		p := domain.MasterPuskesmas{
			ID:        uuid.New(),
			Name:      fmt.Sprintf("Puskesmas %s %d", kec, i),
			Kecamatan: kec,
		}
		db.Create(&p)
		puskList = append(puskList, p)
	}

	// Create 63 Kelurahan
	kelurahanNames := []string{
		// Beji
		"Beji", "Beji Timur", "Kemiri Muka", "Kukusan", "Pondok Cina", "Tanah Baru",
		// Bojongsari
		"Bojongsari", "Bojongsari Baru", "Curug", "Duren Mekar", "Duren Seribu", "Pondok Petir", "Serua",
		// Cilodong
		"Cilodong", "Jatimulya", "Kalibaru", "Kalimulya", "Sukamaju",
		// Cimanggis
		"Cisalak Pasar", "Curug Cimanggis", "Harjamukti", "Mekarsari", "Pasir Gunung Selatan", "Tugu",
		// Cinere
		"Cinere", "Gandul", "Pangkalan Jati", "Pangkalan Jati Baru",
		// Cipayung
		"Bojong Pondok Terong", "Cipayung", "Cipayung Jaya", "Pondok Jaya", "Ratujaya",
		// Limo
		"Grogol", "Krukut", "Limo", "Meruyung",
		// Pancoran Mas
		"Depok", "Depok Jaya", "Mampang", "Pancoran Mas", "Rangkapan Jaya", "Rangkapan Jaya Baru",
		// Sawangan
		"Bedahan", "Cinangka", "Kedaung", "Pasir Putih", "Pengasinan", "Sawangan", "Sawangan Baru",
		// Sukmajaya
		"Abadijaya", "Bakti Jaya", "Cisalak", "Mekar Jaya", "Sukmajaya", "Tirtajaya",
		// Tapos
		"Cilangkap", "Cimpaeun", "Jatijajar", "Leuwinanggung", "Sukamaju Baru", "Sukatani", "Tapos",
	}

	var kelList []domain.MasterKelurahan
	for idx, name := range kelurahanNames {
		pusk := puskList[0] // Default to Beji for Beji kelurahans
		if idx >= 6 && idx < len(puskList) {
			pusk = puskList[idx%len(puskList)]
		}
		k := domain.MasterKelurahan{
			Name:        name,
			Kecamatan:   pusk.Kecamatan,
			PuskesmasID: pusk.ID,
		}
		db.Create(&k)
		kelList = append(kelList, k)
	}

	// 4. Seed Primary Demo Users
	fmt.Println("[5/7] Mendaftarkan Akun Demo 7 Peran Utama...")
	defaultHash := "$2a$10$t9asZm/hciyqysoaKyd68.te9Zr2e5Cy7L2JBeUsDxvDV0bn50HYu" // password123

	superAdmin := domain.User{
		Name:         "SIMULASI - Super Admin Developer",
		Email:        "dev@dbdguard.id",
		Contact:      "08111111111",
		RoleID:       roleMap[domain.RoleAdminDeveloper],
		ActiveStatus: true,
		PasswordHash: defaultHash,
	}
	db.Create(&superAdmin)

	dinkesUser := domain.User{
		Name:         "SIMULASI - Dinas Kesehatan Kota Depok",
		Email:        "admin@depok.go.id",
		Contact:      "08111111112",
		RoleID:       roleMap[domain.RoleDinkes],
		ActiveStatus: true,
		PasswordHash: defaultHash,
	}
	db.Create(&dinkesUser)

	kepalaPusk := domain.User{
		Name:         "SIMULASI - Dr. H. Kepala Puskesmas Beji",
		Email:        "kepala.puskesmas@depok.go.id",
		Contact:      "08111111113",
		RoleID:       roleMap[domain.RoleKepalaPuskesmas],
		LocationID:   &puskBejiID,
		ActiveStatus: true,
		PasswordHash: defaultHash,
	}
	db.Create(&kepalaPusk)

	petugasPE := domain.User{
		Name:         "SIMULASI - Petugas PE Puskesmas Beji",
		Email:        "petugas.puskesmas@depok.go.id",
		Contact:      "08111111114",
		RoleID:       roleMap[domain.RolePetugasPuskesmas],
		LocationID:   &puskBejiID,
		ActiveStatus: true,
		PasswordHash: defaultHash,
	}
	db.Create(&petugasPE)

	timFogging := domain.User{
		Name:         "SIMULASI - Tim Fogging Reaksi Cepat Beji",
		Email:        "tim.fogging@depok.go.id",
		Contact:      "08111111115",
		RoleID:       roleMap[domain.RoleTimFogging],
		LocationID:   &puskBejiID,
		ActiveStatus: true,
		PasswordHash: defaultHash,
	}
	db.Create(&timFogging)

	kaderBeji := domain.User{
		Name:         "SIMULASI - Ibu Kader RW 04 Kelurahan Beji",
		Email:        "kader@depok.go.id",
		Contact:      "08111111116",
		RoleID:       roleMap[domain.RoleKader],
		LocationID:   &puskBejiID,
		ActiveStatus: true,
		PasswordHash: defaultHash,
	}
	db.Create(&kaderBeji)

	wargaUser := domain.User{
		Name:         "SIMULASI - Bpk. Hendra Wijaya (Warga Pelapor)",
		Email:        "warga@gmail.com",
		Contact:      "08111111117",
		RoleID:       roleMap[domain.RoleWarga],
		ActiveStatus: true,
		PasswordHash: defaultHash,
	}
	db.Create(&wargaUser)

	// 5. Seed Inventories (Gudang Dinkes & Puskesmas)
	fmt.Println("[6/7] Menyiapkan Inventaris Insektisida & Logistik Fogging...")
	dinkesItems := []domain.InventoryItem{
		{ItemType: "ABATE", Quantity: 5000, MinimumStock: 1000, Unit: "Kg", Warehouse: "Gudang Utama Dinkes"},
		{ItemType: "INSEKTISIDA", Quantity: 10000, MinimumStock: 2000, Unit: "Liter", Warehouse: "Gudang Utama Dinkes"},
		{ItemType: "ALAT FOGGING", Quantity: 500, MinimumStock: 100, Unit: "Unit", Warehouse: "Gudang Utama Dinkes"},
		{ItemType: "SOLAR", Quantity: 15000, MinimumStock: 3000, Unit: "Liter", Warehouse: "Gudang Utama Dinkes"},
	}
	for i := range dinkesItems {
		db.Create(&dinkesItems[i])
	}

	// Puskesmas Beji Inventory: Insektisida 23 Liter (telah terpotong 2L dari awal 25L)
	invBejiInsektisida := domain.InventoryItem{
		ItemType:     "INSEKTISIDA",
		Quantity:     23.0,
		MinimumStock: 10.0,
		Unit:         "Liter",
		Warehouse:    puskBeji.Name,
	}
	db.Create(&invBejiInsektisida)

	db.Create(&domain.InventoryItem{
		ItemType:     "ABATE",
		Quantity:     60.0,
		MinimumStock: 15.0,
		Unit:         "Kg",
		Warehouse:    puskBeji.Name,
	})
	db.Create(&domain.InventoryItem{
		ItemType:     "ALAT FOGGING",
		Quantity:     6.0,
		MinimumStock: 2.0,
		Unit:         "Unit",
		Warehouse:    puskBeji.Name,
	})

	// Puskesmas Cipayung Inventory: Insektisida 5 Liter (STOK KRITIS < 10L untuk simulasi SPK Aturan 2)
	db.Create(&domain.InventoryItem{
		ItemType:     "INSEKTISIDA",
		Quantity:     5.0, // < 10 Liter
		MinimumStock: 20.0,
		Unit:         "Liter",
		Warehouse:    puskCipayung.Name,
	})

	// 6. Skenario Demo Lengkap & Wilayah Siaga EWS
	fmt.Println("[7/7] Menyusun Skenario Alur Lengkap, Kelurahan Siaga, & Penolakan Lintas Wilayah...")

	// -------------------------------------------------------------
	// SKENARIO 1: ALUR END-TO-END LENGKAP TERPOTONG STOK
	// Warga Lapor -> Verifikasi Kader -> PE -> SPK Fogging Fokus -> Approval -> Fogging Selesai -> Stok Terpotong
	// -------------------------------------------------------------
	reporterDemo := domain.Reporter{
		Name:          "SIMULASI - Bpk. Hendra Wijaya",
		Contact:       "08123456789",
		AnonymousFlag: false,
	}
	db.Create(&reporterDemo)

	reportDemo := domain.CitizenReport{
		ReferenceNo:    "DBD-DPK-2026-DEMO01",
		Category:       "Kasus Demam Berdarah",
		Description:    "SIMULASI - Ditemukan anak demam tinggi 3 hari beruntun dan terdapat genangan drum di belakang rumah warga RT 02.",
		Status:         "COMPLETED",
		ReporterID:     reporterDemo.ID,
		InvestigatorID: &petugasPE.ID,
		FeedbackRating: 5,
		FeedbackNotes:  "SIMULASI - Penanganan tim sangat sigap, fogging telah selesai dilaksanakan.",
		CreatedAt:      time.Now().Add(-48 * time.Hour),
		UpdatedAt:      time.Now().Add(-2 * time.Hour),
	}
	db.Create(&reportDemo)

	locDemo := domain.ReportLocation{
		ReportID:  reportDemo.ID,
		Latitude:  -6.372800,
		Longitude: 106.831500,
		Address:   "SIMULASI - Jl. Margonda Raya No. 100, RT 02 / RW 04",
		RT:        "02",
		RW:        "04",
		Kelurahan: "Beji",
		Kecamatan: "Beji",
	}
	db.Omit("Geom").Create(&locDemo)

	// Verification record by Kader
	db.Create(&domain.Verification{
		ReportID:  reportDemo.ID,
		ActorID:   kaderBeji.ID,
		Decision:  "APPROVED",
		Reason:    "SIMULASI - Laporan valid, kader telah meninjau lokasi dan memverifikasi gejala klinis warga.",
		CreatedAt: time.Now().Add(-42 * time.Hour),
	})

	// PE Result (Larvae = true, Fever = 3 -> SPK Recommends FOGGING_FOKUS)
	peDemo := domain.PEResult{
		ReportID:         reportDemo.ID,
		OfficerID:        petugasPE.ID,
		ValidatorID:      &kepalaPusk.ID,
		Status:           "VALIDATED",
		InspectionDate:   time.Now().Add(-30 * time.Hour),
		RadiusMeters:     100,
		LarvaeFound:      true,
		FeverCount:       3,
		Recommendation:   "FOGGING_FOKUS",
		PsnEdukasiDone:   true,
		Notes:            "SIMULASI - Ditemukan jentik positif di 2 penampungan dan 3 warga sekitar demam. Stok puskesmas memadai (25 L). SPK merekomendasikan Fogging Fokus segera.",
		EvidencePhotoURL: "https://images.unsplash.com/photo-1584036561566-baf8f5f1b144?auto=format&fit=crop&w=600&q=80",
		CreatedAt:        time.Now().Add(-30 * time.Hour),
	}
	db.Create(&peDemo)

	// Fogging Request
	fogReqDemo := domain.FoggingRequest{
		PEID:           peDemo.ID,
		ApprovalStatus: "APPROVED",
		ReservedStock:  2.0,
		CreatedAt:      time.Now().Add(-24 * time.Hour),
	}
	db.Create(&fogReqDemo)

	// Fogging Schedule
	schedDemo := domain.FoggingSchedule{
		RequestID:       fogReqDemo.ID,
		TeamID:          timFogging.ID,
		PlannedTime:     time.Now().Add(-6 * time.Hour),
		ExecutionStatus: "COMPLETED",
		Notes:           "SIMULASI - Fogging fokus selesai dalam radius 100m. Pemakaian insektisida 2 Liter tercatat.",
		CreatedAt:       time.Now().Add(-20 * time.Hour),
		UpdatedAt:       time.Now().Add(-2 * time.Hour),
	}
	db.Create(&schedDemo)

	// Skenario Jadwal Aktif Tim Fogging: DBD-DPK-2026-FOG01 (Siap Didemokan OTW -> Selesai)
	reportFogReady := domain.CitizenReport{
		ReferenceNo:    "DBD-DPK-2026-FOG01",
		Category:       "Kasus Demam Berdarah",
		Description:    "SIMULASI - Kasus demam berdarah klaster RT 03 Beji telah divalidasi SPK dan dijadwalkan untuk fogging fokus.",
		Status:         "FOGGING_SCHEDULED",
		ReporterID:     reporterDemo.ID,
		InvestigatorID: &petugasPE.ID,
		CreatedAt:      time.Now().Add(-24 * time.Hour),
		UpdatedAt:      time.Now().Add(-1 * time.Hour),
	}
	db.Create(&reportFogReady)

	locFogReady := domain.ReportLocation{
		ReportID:  reportFogReady.ID,
		Latitude:  -6.374200,
		Longitude: 106.829000,
		Address:   "SIMULASI - Jl. KH. M. Usman No. 12, RT 03 / RW 04",
		RT:        "03",
		RW:        "04",
		Kelurahan: "Beji",
		Kecamatan: "Beji",
	}
	db.Omit("Geom").Create(&locFogReady)

	peFogReady := domain.PEResult{
		ReportID:         reportFogReady.ID,
		OfficerID:        petugasPE.ID,
		ValidatorID:      &kepalaPusk.ID,
		Status:           "VALIDATED",
		InspectionDate:   time.Now().Add(-12 * time.Hour),
		RadiusMeters:     100,
		LarvaeFound:      true,
		FeverCount:       4,
		Recommendation:   "FOGGING_FOKUS",
		PsnEdukasiDone:   true,
		Notes:            "SIMULASI - Validasi SPK Aturan #1 terpicu: Jentik positif dan demam 4 orang.",
		EvidencePhotoURL: "https://images.unsplash.com/photo-1584036561566-baf8f5f1b144?auto=format&fit=crop&w=600&q=80",
		CreatedAt:        time.Now().Add(-12 * time.Hour),
	}
	db.Create(&peFogReady)

	fogReqReady := domain.FoggingRequest{
		PEID:           peFogReady.ID,
		ApprovalStatus: "APPROVED",
		ReservedStock:  2.0,
		CreatedAt:      time.Now().Add(-6 * time.Hour),
	}
	db.Create(&fogReqReady)

	schedReady := domain.FoggingSchedule{
		RequestID:       fogReqReady.ID,
		TeamID:          timFogging.ID,
		PlannedTime:     time.Now().Add(2 * time.Hour),
		ExecutionStatus: "SCHEDULED",
		Notes:           "SIMULASI - Jadwal fogging fokus aktif untuk demonstrasi tim pelaksana lapangan.",
		CreatedAt:       time.Now().Add(-6 * time.Hour),
		UpdatedAt:       time.Now().Add(-1 * time.Hour),
	}
	db.Create(&schedReady)

	// Log Audit Trail for Scenario 1
	auditLogs := []domain.AuditLog{
		{
			UserID:     wargaUser.ID,
			Action:     "SUBMIT_CITIZEN_REPORT",
			EntityType: "CITIZEN_REPORT",
			EntityID:   reportDemo.ID,
			Changes:    datatypes.JSON(`{"reference_no": "DBD-DPK-2026-DEMO01", "status": "MENUNGGU_VERIFIKASI"}`),
			CreatedAt:  time.Now().Add(-48 * time.Hour),
		},
		{
			UserID:     kaderBeji.ID,
			Action:     "VERIFY_REPORT_VALID",
			EntityType: "CITIZEN_REPORT",
			EntityID:   reportDemo.ID,
			Changes:    datatypes.JSON(`{"verified_by": "Kader Beji", "status": "MENUNGGU_PE_KADER"}`),
			CreatedAt:  time.Now().Add(-42 * time.Hour),
		},
		{
			UserID:     petugasPE.ID,
			Action:     "SUBMIT_PE_RESULT",
			EntityType: "PE_RESULT",
			EntityID:   peDemo.ID,
			Changes:    datatypes.JSON(`{"larva_found": true, "fever_count": 3, "dss_recommendation": "FOGGING_FOKUS"}`),
			CreatedAt:  time.Now().Add(-30 * time.Hour),
		},
		{
			UserID:     kepalaPusk.ID,
			Action:     "APPROVE_FOGGING_RECOMMENDATION",
			EntityType: "FOGGING_REQUEST",
			EntityID:   fogReqDemo.ID,
			Changes:    datatypes.JSON(`{"approved": true, "dss_decision": "FOGGING_FOKUS", "allocated_stock": 2.0}`),
			CreatedAt:  time.Now().Add(-24 * time.Hour),
		},
		{
			UserID:     timFogging.ID,
			Action:     "EXECUTE_FOGGING_COMPLETED",
			EntityType: "FOGGING_SCHEDULE",
			EntityID:   schedDemo.ID,
			Changes:    datatypes.JSON(`{"execution_status": "COMPLETED", "chemical_used_liters": 2.0, "warehouse": "Puskesmas Beji"}`),
			CreatedAt:  time.Now().Add(-2 * time.Hour),
		},
		{
			UserID:     timFogging.ID,
			Action:     "DEDUCT_INVENTORY_STOCK",
			EntityType: "INVENTORY_ITEM",
			EntityID:   invBejiInsektisida.ItemID,
			Changes:    datatypes.JSON(`{"previous_qty": 25.0, "deducted": 2.0, "new_qty": 23.0, "item": "INSEKTISIDA"}`),
			CreatedAt:  time.Now().Add(-2 * time.Hour),
		},
	}
	for i := range auditLogs {
		db.Create(&auditLogs[i])
	}

	// -------------------------------------------------------------
	// SKENARIO 2: KELURAHAN BERSTATUS SIAGA (EWS / KLB ACTIVE)
	// Kelurahan Beji memiliki 6 laporan aktif terdistribusi, memicu alert KLB
	// -------------------------------------------------------------
	fmt.Println("   -> Mempersiapkan Skenario Kelurahan Siaga (Beji - 6 Kasus Aktif)...")
	for i := 1; i <= 6; i++ {
		rep := domain.Reporter{
			Name:          fmt.Sprintf("SIMULASI - Pelapor Beji Siaga %d", i),
			Contact:       fmt.Sprintf("0812999%04d", i),
			AnonymousFlag: false,
		}
		db.Create(&rep)

		rpt := domain.CitizenReport{
			ReferenceNo: fmt.Sprintf("DBD-DPK-2026-KLB%02d", i),
			Category:    "Kasus Demam Berdarah",
			Description: fmt.Sprintf("SIMULASI KLB - Laporan kasus DBD warga Beji lingkungan RW 0%d.", (i%4)+1),
			Status:      "MENUNGGU_PE_KADER",
			ReporterID:  rep.ID,
			CreatedAt:   time.Now().Add(-time.Duration(i*6) * time.Hour),
		}
		db.Create(&rpt)

		loc := domain.ReportLocation{
			ReportID:  rpt.ID,
			Latitude:  -6.371000 + (float64(i) * 0.002),
			Longitude: 106.828000 + (float64(i) * 0.002),
			Address:   fmt.Sprintf("SIMULASI - Jl. Melati Lingkungan RT 0%d / RW 0%d", i, (i%4)+1),
			RT:        fmt.Sprintf("0%d", i),
			RW:        fmt.Sprintf("0%d", (i%4)+1),
			Kelurahan: "Beji",
			Kecamatan: "Beji",
		}
		db.Omit("Geom").Create(&loc)
	}

	// Active KLB Alert Record
	db.Create(&domain.KlbConfiguration{
		RegionScope:   "KELURAHAN",
		MinCases:      5,
		TimeframeDays: 7,
		IsAutoTrigger: true,
		LastUpdatedBy: dinkesUser.ID,
	})

	db.Create(&domain.KlbAlert{
		AreaName:    "Beji",
		AreaType:    "KELURAHAN",
		CasesCount:  6,
		IsActive:    true,
		TriggeredAt: time.Now().Add(-12 * time.Hour),
	})

	// -------------------------------------------------------------
	// SKENARIO 3: SKENARIO DITOLAK LINTAS WILAYAH
	// Laporan di perbatasan wilayah kerja kader, ditolak resmi dengan alasan
	// -------------------------------------------------------------
	fmt.Println("   -> Mempersiapkan Skenario Ditolak Lintas Wilayah...")
	repReject := domain.Reporter{
		Name:          "SIMULASI - Ibu Siti Rahma (Warga Perbatasan)",
		Contact:       "087811122233",
		AnonymousFlag: false,
	}
	db.Create(&repReject)

	rptReject := domain.CitizenReport{
		ReferenceNo: "DBD-DPK-2026-REJECT01",
		Category:    "Genangan Jentik Liar",
		Description: "SIMULASI LINTAS WILAYAH - Genangan air di area tanah kosong perbatasan Kelurahan Beji dan wilayah kabupaten tetangga.",
		Status:      "REJECTED",
		ReporterID:  repReject.ID,
		CreatedAt:   time.Now().Add(-20 * time.Hour),
		UpdatedAt:   time.Now().Add(-15 * time.Hour),
	}
	db.Create(&rptReject)

	locReject := domain.ReportLocation{
		ReportID:  rptReject.ID,
		Latitude:  -6.388500,
		Longitude: 106.812000,
		Address:   "SIMULASI - Jl. Perbatasan Sasak Panjang RT 01 / RW 01 (Luar Batas Wilayah)",
		RT:        "01",
		RW:        "01",
		Kelurahan: "Beji",
		Kecamatan: "Beji",
	}
	db.Omit("Geom").Create(&locReject)

	db.Create(&domain.Verification{
		ReportID:  rptReject.ID,
		ActorID:   kaderBeji.ID,
		Decision:  "REJECTED",
		Reason:    "SIMULASI - Ditolak: Lokasi laporan berada di luar wilayah yurisdiksi kerja Kader Kelurahan Beji (Lintas Wilayah Administratif). Laporan dialihkan ke kanal wilayah terkait.",
		CreatedAt: time.Now().Add(-15 * time.Hour),
	})

	db.Create(&domain.AuditLog{
		UserID:     kaderBeji.ID,
		Action:     "REJECT_REPORT_CROSS_BORDER",
		EntityType: "CITIZEN_REPORT",
		EntityID:   rptReject.ID,
		Changes:    datatypes.JSON(`{"reason": "Di luar yurisdiksi kerja Kader RW 04 Beji", "status": "REJECTED"}`),
		CreatedAt:  time.Now().Add(-15 * time.Hour),
	})

	// -------------------------------------------------------------
	// SKENARIO 4: LAPORAN SIAP VALIDASI SPK DI PUSKESMAS (Untuk Uji Coba Langsung di Layar SPK)
	// Status: MENUNGGU_VALIDASI_PUSKESMAS dengan PE Result lengkap
	// -------------------------------------------------------------
	repSPK := domain.Reporter{
		Name:          "SIMULASI - Bpk. Ahmad Yani (Uji SPK)",
		Contact:       "085600011122",
		AnonymousFlag: false,
	}
	db.Create(&repSPK)

	rptSPK := domain.CitizenReport{
		ReferenceNo: "DBD-DPK-2026-SPK01",
		Category:    "Kasus Demam Berdarah",
		Description: "SIMULASI SPK - Pemeriksaan jentik berkala oleh kader menunjukkan indikator positif dengan 4 orang warga demam.",
		Status:      "MENUNGGU_VALIDASI_PUSKESMAS",
		ReporterID:  repSPK.ID,
		CreatedAt:   time.Now().Add(-8 * time.Hour),
	}
	db.Create(&rptSPK)

	locSPK := domain.ReportLocation{
		ReportID:  rptSPK.ID,
		Latitude:  -6.375000,
		Longitude: 106.825000,
		Address:   "SIMULASI - Jl. Nusantara Raya No. 45, RT 03 / RW 02, Beji",
		RT:        "03",
		RW:        "02",
		Kelurahan: "Beji",
		Kecamatan: "Beji",
	}
	db.Omit("Geom").Create(&locSPK)

	db.Create(&domain.PEResult{
		ReportID:         rptSPK.ID,
		OfficerID:        kaderBeji.ID,
		Status:           "PENDING_VALIDATION",
		InspectionDate:   time.Now().Add(-4 * time.Hour),
		RadiusMeters:     100,
		LarvaeFound:      true,
		FeverCount:       4,
		Recommendation:   "FOGGING_FOKUS",
		PsnEdukasiDone:   true,
		Notes:            "SIMULASI - Ditemukan jentik di bak terbuka, 4 warga demam tinggi. Memenuhi kriteria Fogging Fokus.",
		EvidencePhotoURL: "https://images.unsplash.com/photo-1584036561566-baf8f5f1b144?auto=format&fit=crop&w=600&q=80",
		CreatedAt:        time.Now().Add(-4 * time.Hour),
	})

	// -------------------------------------------------------------
	// LAPORAN SINTETIS TAMBAHAN DI 11 KECAMATAN (Agar Peta & Statistik Penuh)
	// -------------------------------------------------------------
	fmt.Println("   -> Mengisi 30 laporan sintetis sebaran 11 Kecamatan Kota Depok...")
	for i := 1; i <= 30; i++ {
		kel := kelList[rand.Intn(len(kelList))]
		rep := domain.Reporter{
			Name:          fmt.Sprintf("SIMULASI - Warga %s %d", kel.Kecamatan, i),
			Contact:       fmt.Sprintf("0813%06d", rand.Intn(999999)),
			AnonymousFlag: i%6 == 0,
		}
		db.Create(&rep)

		status := "MENUNGGU_PE_KADER"
		if i%4 == 0 {
			status = "COMPLETED"
		} else if i%5 == 0 {
			status = "REJECTED"
		}

		rpt := domain.CitizenReport{
			ReferenceNo: fmt.Sprintf("DBD-DPK-2026-%04d", 200+i),
			Category:    "Gejala Demam Berdarah",
			Description: fmt.Sprintf("SIMULASI - Laporan deteksi dini DBD di wilayah Kelurahan %s.", kel.Name),
			Status:      status,
			ReporterID:  rep.ID,
			CreatedAt:   time.Now().Add(-time.Duration(rand.Intn(168)) * time.Hour),
		}
		db.Create(&rpt)

		lat := -6.402484 + (rand.Float64()-0.5)*0.09
		lng := 106.794243 + (rand.Float64()-0.5)*0.09

		loc := domain.ReportLocation{
			ReportID:  rpt.ID,
			Latitude:  lat,
			Longitude: lng,
			Address:   fmt.Sprintf("SIMULASI - Jl. Melati Indah No. %d, %s", i, kel.Name),
			RT:        fmt.Sprintf("0%d", (i%5)+1),
			RW:        fmt.Sprintf("0%d", (i%4)+1),
			Kelurahan: kel.Name,
			Kecamatan: kel.Kecamatan,
		}
		db.Omit("Geom").Create(&loc)
	}

	fmt.Println("\n================================================================")
	fmt.Println("✅ Database Seeder Berhasil Dijalankan Secara Idempoten!")
	fmt.Println("   1. Skenario End-to-End Fogging & Potong Stok: DBD-DPK-2026-DEMO01")
	fmt.Println("   2. Skenario Kelurahan Siaga (EWS KLB): Kelurahan Beji (6 Kasus)")
	fmt.Println("   3. Skenario Ditolak Lintas Wilayah: DBD-DPK-2026-REJECT01")
	fmt.Println("   4. Kasus Siap Validasi SPK: DBD-DPK-2026-SPK01")
	fmt.Println("   5. Akun Demo Utama: dev@dbdguard.id, admin@depok.go.id, dll.")
	fmt.Println("================================================================")
}
