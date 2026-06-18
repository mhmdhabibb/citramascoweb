package config

import (
	"citramascoweb-backend/internal/modules/auth"
	"citramascoweb-backend/internal/modules/offer"
	"citramascoweb-backend/internal/modules/rooms"
	"citramascoweb-backend/internal/modules/rooms/category"
	"citramascoweb-backend/internal/modules/rooms/reservation"
	"citramascoweb-backend/internal/modules/rooms/types"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := Config("DB_URL")

	if dsn == "" {
		log.Fatalf("Database connection failed: DB_URL environment variable is missing or empty")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	err = db.AutoMigrate(&auth.User{}, &category.Category{}, &types.Types{}, &rooms.Room{}, &reservation.Reservation{}, &offer.Offer{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	return db
}
