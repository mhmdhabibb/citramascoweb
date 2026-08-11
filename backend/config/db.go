package config

import (
	"citramascoweb-backend/internal/modules/rooms/reservation"
	"citramascoweb-backend/internal/modules/finance"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := Config("DB_URL")

	// Rangkai menjadi format DSN PostgreSQL yang resmi

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

	db.AutoMigrate(&reservation.Reservation{})
	db.AutoMigrate(&finance.ChartOfAccount{}, &finance.CashTransaction{}, &finance.GeneralJournal{}, &finance.Invoice{})

	seedCOA(db)

	return db
}

func seedCOA(db *gorm.DB) {
	coas := []finance.ChartOfAccount{
		{Code: "1-100", Name: "Kas Finance", Type: "Asset"},
		{Code: "1-101", Name: "Kas Front Desk", Type: "Asset"},
		{Code: "1-102", Name: "Bank BCA", Type: "Asset"},
		{Code: "4-100", Name: "Pendapatan Sewa", Type: "Revenue"},
	}

	for _, coa := range coas {
		var existing finance.ChartOfAccount
		if err := db.Where("code = ?", coa.Code).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				coa.CreatedAt = time.Now()
				coa.UpdatedAt = time.Now()
				db.Create(&coa)
			}
		}
	}
}
