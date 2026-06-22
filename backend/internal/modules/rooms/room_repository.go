package rooms

import (
	"citramascoweb-backend/internal/modules/offer"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type RoomRepositoryInterface interface {
	GetAll() ([]Room, error)
	GetById(id string) (*Room, error)
	Create(room *Room) error
	Update(room *Room, id string) error
	Delete(id string) error
	GetLastRoom() (*Room, error)
	FilerByCategory(categoryId string) ([]Room, error)
	FilterByType(typeId string) ([]Room, error)
	Filter(status, availabilityStatus string, checkinDate, checkoutDate time.Time) ([]Room, error)
	GetOfferByCode(code string) (*offer.Offer, error)
	CountRoomsWithOfferCode(code string, excludeRoomId string) (int64, error)
	DecrementOfferQuota(code string) error
	IncrementOfferQuota(code string) error
	UpdateStatus(id string, status RoomStatus) error
}

type roomRepo struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepositoryInterface {
	return &roomRepo{db: db}
}

func (r *roomRepo) GetAll() ([]Room, error) {
	var rooms []Room

	err := r.db.Preload("Category").Preload("Type").Find(&rooms).Error
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *roomRepo) GetById(id string) (*Room, error) {
	var room Room
	err := r.db.Preload("Category").Preload("Type").Where("id = ?", id).First(&room).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepo) Create(room *Room) error {
	return r.db.Omit("Category", "Type").Create(room).Error
}

func (r *roomRepo) Update(room *Room, id string) error {
	return r.db.Model(&Room{}).Where("id = ?", id).Omit("Category", "Type").Updates(room).Error
}

func (r *roomRepo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Room{}).Error
}

func (r *roomRepo) GetLastRoom() (*Room, error) {
	var room Room
	err := r.db.Order("created_at desc").First(&room).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &room, nil
}

func (r *roomRepo) FilerByCategory(categoryId string) ([]Room, error) {

	var rooms []Room

	err := r.db.Preload("Type").Preload("Category").Where("category_id = ?", categoryId).Find(&rooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil

}

func (r *roomRepo) FilterByType(typeId string) ([]Room, error) {

	var rooms []Room

	err := r.db.Preload("Type").Preload("Category").Where("type_id = ?", typeId).Find(&rooms).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return rooms, nil

}

func (r *roomRepo) Filter(status, availabilityStatus string, checkinDate, checkoutDate time.Time) ([]Room, error) {
	var rooms []Room

	// 1. Fetch all rooms (Termasuk status fisiknya dari DB jika nanti kamu simpan status fisik ke DB)
	err := r.db.Preload("Type").Preload("Category").Find(&rooms).Error
	if err != nil {
		return nil, err
	}

	// Default date range jika tidak dispesifikasikan
	if checkinDate.IsZero() {
		checkinDate = time.Now()
	}
	if checkoutDate.IsZero() {
		checkoutDate = checkinDate.Add(24 * time.Hour)
	}

	// 2. Query ID kamar yang OVERLAP / SEDANG TERISI pada periode tersebut
	var occupiedRoomIds []string
	err = r.db.Table("reservations").
		Select("room_id").
		Where("room_id IS NOT NULL AND room_id != ''").
		// Logika overlap: Checkin kustomer baru < Checkout reservasi lama DAN Checkout kustomer baru > Checkin reservasi lama
		Where("checkin_date < ? AND checkout_date > ? AND status IN ?",
			checkoutDate,
			checkinDate,
			[]string{"pending", "approved", "checked-in"},
		).
		Distinct("room_id").
		Pluck("room_id", &occupiedRoomIds).Error // Menggunakan Pluck langsung ke slice string agar lebih efisien

	if err != nil {
		return nil, err
	}

	// Buat map pencarian untuk kamar yang terisi (occupied)
	occupiedMap := make(map[string]bool)
	for _, id := range occupiedRoomIds {
		occupiedMap[id] = true
	}

	// 3. Populasi field virtual menggunakan ENUM yang baru
	for i := range rooms {
		// Set Availability Status berdasarkan transaksi booking
		if occupiedMap[rooms[i].Id] {
			rooms[i].AvailabilityStatus = RoomAvailabilityOccupied
		} else {
			rooms[i].AvailabilityStatus = RoomAvailabilityAvailable
		}

		// Set Status Fisik (Contoh logika: jika di DB tidak maintenance, default ke Active)
		// Jika ke depannya kamu simpan status fisik ke DB, bagian ini tinggal menyesuaikan kolom DB-mu
		if rooms[i].Status == "" {
			rooms[i].Status = RoomStatusActive
		}
	}

	// 4. Helper match function menggunakan Enum string
	isStatusMatch := func(roomVal, queryVal string) bool {
		q := strings.ToLower(queryVal)
		rv := strings.ToLower(roomVal)

		// Normalisasi pencarian bahasa Indonesia ke Enum internasional kamu
		if q == "tersedia" {
			q = "available"
		} else if q == "tidak tersedia" || q == "tidak_tersedia" {
			q = "occupied"
		}

		return rv == q
	}

	// 5. In-memory filtering berdasarkan parameter filter
	var filteredRooms []Room
	for _, room := range rooms {
		matchesStatus := true
		matchesAvailability := true

		if status != "" {
			matchesStatus = isStatusMatch(string(room.Status), status)
		}
		if availabilityStatus != "" {
			matchesAvailability = isStatusMatch(string(room.AvailabilityStatus), availabilityStatus)
		}

		if matchesStatus && matchesAvailability {
			filteredRooms = append(filteredRooms, room)
		}
	}

	return filteredRooms, nil
}

func (r *roomRepo) GetOfferByCode(code string) (*offer.Offer, error) {
	var o offer.Offer
	err := r.db.Where("code = ?", code).First(&o).Error
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *roomRepo) CountRoomsWithOfferCode(code string, excludeRoomId string) (int64, error) {
	var count int64
	query := r.db.Model(&Room{}).Where("offer_code = ?", code)
	if excludeRoomId != "" {
		query = query.Where("id != ?", excludeRoomId)
	}
	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *roomRepo) DecrementOfferQuota(code string) error {
	return r.db.Model(&offer.Offer{}).Where("code = ? AND max_quota IS NOT NULL AND max_quota > 0", code).UpdateColumn("max_quota", gorm.Expr("max_quota - ?", 1)).Error
}

func (r *roomRepo) IncrementOfferQuota(code string) error {
	return r.db.Model(&offer.Offer{}).Where("code = ? AND max_quota IS NOT NULL", code).UpdateColumn("max_quota", gorm.Expr("max_quota + ?", 1)).Error
}

func (r *roomRepo) UpdateStatus(id string, status RoomStatus) error {
	log.Printf("[DEBUG][ROOM-REPO] Memulai update status fisik kamar ID: %s menjadi '%s'", id, status)

	result := r.db.Model(&Room{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		log.Printf("[ERROR][ROOM-REPO] Gagal update status kamar ID: %s. Error: %v", id, result.Error)
		return result.Error
	}

	log.Printf("[DEBUG][ROOM-REPO] Sukses update status kamar ID: %s. Rows affected: %d", id, result.RowsAffected)
	return nil
}
