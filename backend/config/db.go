package config

import (
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

	db, err := gorm.Open(postgres.New(postgres.Config{

		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: false,
	})

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	return db
}
