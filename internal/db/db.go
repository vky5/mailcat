package db

import (
	"log"

	"github.com/vky5/mailcat/internal/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("mailcat.db"), &gorm.Config{})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	err = DB.AutoMigrate(&models.Account{}, &models.Email{}) // automatically creates or updates the database to matches the go struct defined here
	if err != nil {
		log.Fatalf("failed to migrate models: %v", err)
	}
}
