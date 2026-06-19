package config

import (
	"citramascoweb-backend/internal/modules/auth"
	"citramascoweb-backend/internal/modules/offer"
	"citramascoweb-backend/internal/modules/rooms"
	"citramascoweb-backend/internal/modules/rooms/category"
	"citramascoweb-backend/internal/modules/rooms/reservation"
	"citramascoweb-backend/internal/modules/rooms/types"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	host := Config("DB_HOST")
	user := Config("DB_USER")
	password := Config("DB_PASSWORD")
	port := Config("DB_PORT")
	dbname := Config("DB_NAME")

	// Rangkai menjadi format DSN PostgreSQL yang resmi
	dsn := fmt.Sprintf("host=%s user=%s password=%s port=%s dbname=%s sslmode=require",
		host, user, password, port, dbname)
	if dsn == "" {
		log.Fatalf("Database connection failed: DB_URL environment variable is missing or empty")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
