package main

import (
	"fmt"
	"log"

	"github.com/dbd-guard/backend/internal/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=secret dbname=dbd_guard port=5433 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	var primaryEmails = []string{
		"dev@dbdguard.id",
		"admin@dbdguard.id",
		"admin@depok.go.id",
		"kepala.puskesmas@depok.go.id",
		"petugas.puskesmas@depok.go.id",
		"tim.fogging@depok.go.id",
		"kader@depok.go.id",
		"warga@gmail.com",
	}

	fmt.Println("=== AKUN DEMO UTAMA SISTEM DBD GUARD ===")
	for _, email := range primaryEmails {
		var u domain.User
		if err := db.Preload("Role").Where("email = ?", email).First(&u).Error; err == nil {
			locName := "-"
			if u.LocationID != nil {
				var pusk domain.MasterPuskesmas
				if err := db.Where("id = ?", *u.LocationID).First(&pusk).Error; err == nil {
					locName = pusk.Name
				}
			}
			fmt.Printf("Role: %-23s | Email: %-30s | Name: %-32s | Wilayah: %s\n", u.Role.RoleName, u.Email, u.Name, locName)
		} else {
			fmt.Printf("Email: %s NOT FOUND\n", email)
		}
	}
}
