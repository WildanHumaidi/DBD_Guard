package main

import (
	"fmt"
	"log"

	"github.com/dbd-guard/backend/internal/config"
	"github.com/dbd-guard/backend/internal/domain"
	"github.com/google/uuid"
)

func main() {
	config.LoadConfig()
	config.ConnectDB()
	db := config.DB
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}

	inventories := []domain.InventoryItem{
		{ItemID: uuid.New(), ItemType: "ABATE", Quantity: 50, MinimumStock: 20, Unit: "Kg", Warehouse: "Puskesmas Pancoran Mas"},
		{ItemID: uuid.New(), ItemType: "INSEKTISIDA", Quantity: 150, MinimumStock: 50, Unit: "Liter", Warehouse: "Puskesmas Pancoran Mas"},
		{ItemID: uuid.New(), ItemType: "ALAT FOGGING", Quantity: 5, MinimumStock: 2, Unit: "Unit", Warehouse: "Puskesmas Pancoran Mas"},
	}

	for _, inv := range inventories {
		// Check if exists first to avoid duplicates
		var existing domain.InventoryItem
		if err := db.Where("item_type = ? AND warehouse = ?", inv.ItemType, inv.Warehouse).First(&existing).Error; err != nil {
			db.Create(&inv)
			fmt.Printf("Inserted %s for %s\n", inv.ItemType, inv.Warehouse)
		} else {
			fmt.Printf("Already exists %s for %s\n", inv.ItemType, inv.Warehouse)
		}
	}
	fmt.Println("Done")
}
