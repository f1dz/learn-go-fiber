package migration

import (
	"fiber-api/models"
	"log"

	"gorm.io/gorm"
)

func MigrateDB(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
	)

	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}
