package main

import (
	"citramascoweb-backend/config"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 1. Definisikan Model User
type User struct {
	ID        string `gorm:"primaryKey"`
	FirstName string
	LastName  string
	Username  string `gorm:"uniqueIndex;not null"`
	Password  string `gorm:"not null"` // Akan menyimpan password yang sudah di-hash
	Role      string `gorm:"not null;default:'user'"`
}

// 2. Fungsi Helper untuk Hashing Password
func HashPassword(password string) (string, error) {
	// Cost 14 adalah standar yang baik untuk keseimbangan kecepatan & keamanan saat ini
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// 3. Fungsi Seeder Utama
func SeedAdminUser(db *gorm.DB) {
	// Data admin yang ingin kita masukkan
	adminData := User{
		ID:        uuid.New().String(),
		FirstName: "Admin",
		LastName:  "",
		Username:  "admin.citramas",
		Role:      "admin",
	}

	rawPassword := "LNG.Admin#Citramas!2026" // Password asli

	// Pengecekan Idempotency: Jangan buat jika username sudah ada
	var existingUser User
	err := db.Where("username = ?", adminData.Username).First(&existingUser).Error

	if err == nil {
		fmt.Println("⏩ Skip seeding: User admin dengan username", adminData.Username, "sudah ada.")
		return
	}

	if err != gorm.ErrRecordNotFound {
		// Terjadi error database selain "data tidak ditemukan"
		log.Fatalf("❌ Error saat mengecek database: %v", err)
	}

	// Hash password sebelum disimpan
	hashedPassword, err := HashPassword(rawPassword)
	if err != nil {
		log.Fatalf("❌ Gagal melakukan hash password: %v", err)
	}
	adminData.Password = hashedPassword

	// Simpan data ke database
	if err := db.Create(&adminData).Error; err != nil {
		log.Fatalf("❌ Gagal membuat akun admin: %v", err)
	}

	fmt.Println("✅ Berhasil! Akun admin berhasil di-seed ke database.")
}

func main() {

	dsn := config.Config("DB_URL")
	// Buka koneksi database (Contoh ini pakai SQLite file "app.db")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database")
	}

	// Auto-Migrasi (Membuat tabel otomatis jika belum ada)
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi tabel")
	}

	// Jalankan fungsi seeder
	SeedAdminUser(db)
}
