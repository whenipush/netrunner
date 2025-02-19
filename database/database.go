package database

import (
	"log"
	"netrunner/models"
	"netrunner/parser"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var PostgreDB *gorm.DB

func Connect() {
	// Connect to database SQLite
	var err error
	DB, err = gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Database connected successfully!")
	}
	dsn := "user=postgres password=postgres dbname=NetRunnerVulns port=5432 sslmode=disable"
	PostgreDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//PostgreDB, err = gorm.Open(postgres.Open("vulndatabase.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	} else {
		log.Println("Database connected successfully!")
	}
	// Автоматическая миграция модели
	DB.AutoMigrate(&models.Host{}, &models.Group{}, &models.TaskStatus{})
	if err := PostgreDB.AutoMigrate(&models.Vulnerability{}, &models.Exploits{}, &models.Description{}, &parser.CPE{}, &models.CWE{}, &models.Solutions{}, &models.Workarounds{}, &parser.CVSS{}, &parser.CVSS3{}); err != nil {
		log.Printf("Error migrating: %v", err)
	}
}
