package rooms

import (
	"citramascoweb-backend/internal/modules/offer"
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

	// Fetch all rooms
	err := r.db.Preload("Type").Preload("Category").Find(&rooms).Error
	if err != nil {
		return nil, err
	}

	// Default date range if not specified
	if checkinDate.IsZero() {
		checkinDate = time.Now()
	}
	if checkoutDate.IsZero() {
		checkoutDate = checkinDate.Add(24 * time.Hour)
	}

	// Query available room IDs for the period
	var availableRooms []Room
	subQuery := r.db.Table("reservations").
		Select("room_id").
		Where("room_id IS NOT NULL AND room_id != ''").
		Where("checkin_date < ? AND checkout_date > ? AND status IN ?",
			checkoutDate,
			checkinDate,
			[]string{"pending", "approved", "checked-in"},
		)

	// Fetch all rooms not in overlapping reservations
	err = r.db.Select("id").Table("rooms").Where("id NOT IN (?)", subQuery).Find(&availableRooms).Error
	if err != nil {
		return nil, err
	}

	availableMap := make(map[string]bool)
	for _, room := range availableRooms {
		availableMap[room.Id] = true
	}

	// Populate virtual fields: Status and AvailabilityStatus
	for i := range rooms {
		if availableMap[rooms[i].Id] {
			rooms[i].Status = "tersedia"
			rooms[i].AvailabilityStatus = "tersedia"
		} else {
			rooms[i].Status = "tidak tersedia"
			rooms[i].AvailabilityStatus = "tidak tersedia"
		}
	}

	// Helper helper match function to handle various status naming conventions
	isStatusMatch := func(roomVal, queryVal string) bool {
		q := strings.ToLower(queryVal)
		rv := strings.ToLower(roomVal)

		if q == "available" || q == "tersedia" {
			return rv == "tersedia" || rv == "available"
		}
		if q == "occupied" || q == "tidak tersedia" || q == "tidak_tersedia" {
			return rv == "tidak tersedia" || rv == "tidak_tersedia" || rv == "occupied"
		}
		return rv == q
	}

	// Perform in-memory filtering if filter params are set
	var filteredRooms []Room
	for _, room := range rooms {
		matchesStatus := true
		matchesAvailability := true

		if status != "" {
			matchesStatus = isStatusMatch(room.Status, status)
		}
		if availabilityStatus != "" {
			matchesAvailability = isStatusMatch(room.AvailabilityStatus, availabilityStatus)
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
