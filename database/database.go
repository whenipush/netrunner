package database

import (
	"log"
	"netrunner/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Connect to database SQLite
	var err error
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Database connected successfully!")
	}

	// Автоматическая миграция модели
	DB.AutoMigrate(&models.Host{}, &models.Group{})

}
