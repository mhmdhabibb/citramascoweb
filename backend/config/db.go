package config

import (
	"citramascoweb-backend/internal/modules/auth"
	"citramascoweb-backend/internal/modules/category"
	"citramascoweb-backend/internal/modules/types"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := Config("DB_URL")

	if dsn == "" {
		log.Fatalf("Database connection failed: DB_URL environment variable is missing or empty")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	db.AutoMigrate(&auth.User{}, &category.Category{}, &types.Types{})

	return db
}
