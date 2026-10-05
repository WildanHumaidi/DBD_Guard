package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDB initializes the database connection
func ConnectDB() {
	if Env == nil {
		log.Fatal("Config is not loaded. Call LoadConfig() first.")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		Env.DBHost, Env.DBUser, Env.DBPassword, Env.DBName, Env.DBPort, Env.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Log the error but don't fatally crash during scaffolding if the DB isn't running yet.
		log.Printf("Warning: Failed to connect to database at %s:%s. \nError: %v\n", Env.DBHost, Env.DBPort, err)
		return
	}

	log.Println("Database connected successfully")
	DB = db
}
