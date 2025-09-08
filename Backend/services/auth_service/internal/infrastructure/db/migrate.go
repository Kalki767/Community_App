package db

import (
	"log"

	"auth/internal/domain/entity"
	"gorm.io/gorm"
)

// RunMigrations performs GORM auto-migrations
func RunMigrations(db *gorm.DB) {
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(
		&entity.User{},
		&entity.Session{},
	); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Database migration completed successfully.")
}
