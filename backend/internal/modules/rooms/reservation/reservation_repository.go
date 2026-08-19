package reservation

import (
	"citramascoweb-backend/internal/modules/rooms"
	"time"

	"gorm.io/gorm"
)

type ReservationRepositoryInterface interface {
	GetAll() ([]Reservation, error)
	GetById(id string) (*Reservation, error)
	Create(reservation *Reservation) error
	Update(reservation *Reservation, id string) error
	Delete(id string) error
	GetByStatus(status string) ([]Reservation, error)
	ApproveReservation(id string) error
	CancelReservation(id string) error
	RejectReservation(id string) error
	CheckAvailability(checkInDate time.Time, checkOutDate time.Time) ([]rooms.Room, error)
	CheckIn(id string, roomId string) error
	CheckOut(id string, roomId string) error
}

type reservationRepo struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepositoryInterface {
	return &reservationRepo{db: db}
}

func (r *reservationRepo) GetAll() ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("Room").Preload("RoomUnit").Preload("RoomUnit.Floor").Preload("Channel").Preload("Logs", func(db *gorm.DB) *gorm.DB {
		return db.Order("timestamp asc")
	}).Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepo) GetById(id string) (*Reservation, error) {
	var reservation Reservation
	err := r.db.Preload("Room").Preload("RoomUnit").Preload("RoomUnit.Floor").Preload("Channel").Preload("Logs", func(db *gorm.DB) *gorm.DB {
		return db.Order("timestamp asc")
	}).Where("id = ?", id).First(&reservation).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepo) Create(reservation *Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *reservationRepo) Update(reservation *Reservation, id string) error {
	return r.db.Model(&Reservation{}).Where("id = ?", id).Updates(reservation).Error
}

func (r *reservationRepo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&Reservation{}).Error
}

func (r *reservationRepo) GetByStatus(status string) ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("User").Preload("Room").Preload("RoomUnit").Preload("RoomUnit.Floor").Preload("Channel").Preload("Logs", func(db *gorm.DB) *gorm.DB {
		return db.Order("timestamp asc")
	}).Where("status = ?", status).Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepo) ApproveReservation(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var currentRes Reservation
		if err := tx.Where("id = ?", id).First(&currentRes).Error; err != nil {
			return err
		}

		assignedUnitId := currentRes.RoomUnitId

		// Jika belum punya kamar fisik, cari kamar fisik kosong secara berurutan (Sequential Auto-Assign)
		if assignedUnitId == nil || *assignedUnitId == "" {
			var autoUnit struct {
				Id string
			}
			err := tx.Table("room_units").
				Select("room_units.id").
				Joins("LEFT JOIN floors ON floors.id = room_units.floor_id").
				Where("room_units.room_id = ? AND room_units.status = ?", currentRes.RoomId, "available").
				Order("floors.floor_number ASC, room_units.room_number ASC").
				First(&autoUnit).Error

			if err == nil && autoUnit.Id != "" {
				assignedUnitId = &autoUnit.Id
			}
		}

		updateData := map[string]interface{}{
			"status": ReservationStatusApproved,
		}
		if assignedUnitId != nil && *assignedUnitId != "" {
			updateData["room_unit_id"] = *assignedUnitId
		}

		return tx.Model(&Reservation{}).Where("id = ?", id).Updates(updateData).Error
	})
}

func (r *reservationRepo) CancelReservation(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var currentRes Reservation
		if err := tx.Where("id = ?", id).First(&currentRes).Error; err == nil {
			// Jika ada unit kamar fisik yang sudah ter-assign, kembalikan statusnya menjadi available
			if currentRes.RoomUnitId != nil && *currentRes.RoomUnitId != "" {
				tx.Table("room_units").Where("id = ?", *currentRes.RoomUnitId).Update("status", "available")
			}
		}

		// Lepaskan room_unit_id dan ubah status reservasi menjadi cancel
		return tx.Model(&Reservation{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":       ReservationStatusCancel,
			"room_unit_id": nil,
		}).Error
	})
}

func (r *reservationRepo) RejectReservation(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var currentRes Reservation
		if err := tx.Where("id = ?", id).First(&currentRes).Error; err == nil {
			// Jika ada unit kamar fisik yang sudah ter-assign, kembalikan statusnya menjadi available
			if currentRes.RoomUnitId != nil && *currentRes.RoomUnitId != "" {
				tx.Table("room_units").Where("id = ?", *currentRes.RoomUnitId).Update("status", "available")
			}
		}

		// Lepaskan room_unit_id dan ubah status reservasi menjadi rejected
		return tx.Model(&Reservation{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":       ReservationStatusRejected,
			"room_unit_id": nil,
		}).Error
	})
}

func (r *reservationRepo) CheckAvailability(checkInDate time.Time, checkOutDate time.Time) ([]rooms.Room, error) {
	var availableRooms []rooms.Room

	// Subquery to find room IDs that have overlapping reservations in "pending", "approved", or "checked-in" status
	subQuery := r.db.Model(&Reservation{}).
		Select("room_id").
		Where("room_id IS NOT NULL AND room_id != ''").
		Where("checkin_date < ? AND checkout_date > ? AND status IN ?",
			checkOutDate,
			checkInDate,
			[]ReservationStatus{ReservationStatusPending, ReservationStatusApproved, ReservationStatusCheckedIn},
		)

	// Fetch all rooms that are not in the subquery result
	err := r.db.Preload("Category").Preload("Type").
		Where("id NOT IN (?)", subQuery).
		Where("LOWER(status) != ?", "maintenance"). // Pengunci utama agar kamar rusak tidak bisa di-reservasi
		Find(&availableRooms).Error

	if err != nil {
		return nil, err
	}
	return availableRooms, nil
}

func (r *reservationRepo) CheckIn(id string, roomId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		isEarly := now.Hour() < 14

		// 1. Ambil data reservasi saat ini untuk cek apakah sudah ada unit fisik
		var currentRes Reservation
		if err := tx.Where("id = ?", id).First(&currentRes).Error; err != nil {
			return err
		}

		assignedUnitId := currentRes.RoomUnitId

		// 2. Jika belum ada unit fisik, lakukan Auto-Assign berurutan (Sequential Allocation)
		if assignedUnitId == nil || *assignedUnitId == "" {
			var autoUnit struct {
				Id string
			}
			err := tx.Table("room_units").
				Select("room_units.id").
				Joins("LEFT JOIN floors ON floors.id = room_units.floor_id").
				Where("room_units.room_id = ? AND room_units.status = ?", roomId, "available").
				Order("floors.floor_number ASC, room_units.room_number ASC").
				First(&autoUnit).Error

			if err == nil && autoUnit.Id != "" {
				assignedUnitId = &autoUnit.Id
			}
		}

		// 3. Ubah status reservasi menjadi checked-in dan catat jam check-in riil
		updateData := map[string]interface{}{
			"status":            ReservationStatusCheckedIn,
			"actual_checkin_at": now,
			"is_early_checkin":  isEarly,
		}
		if assignedUnitId != nil && *assignedUnitId != "" {
			updateData["room_unit_id"] = *assignedUnitId
		}

		if err := tx.Model(&Reservation{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
			return err
		}

		// 4. Simpan Riwayat Log Check In Harian
		logEntry := ReservationLog{
			Id:            time.Now().Format("20060102150405.000000"),
			ReservationId: id,
			Action:        "check_in",
			Timestamp:     now,
			IsEarly:       isEarly,
			Notes:         "Tamu Check In",
			CreatedAt:     now,
		}
		if err := tx.Create(&logEntry).Error; err != nil {
			return err
		}

		// 5. Ubah status kamar fisik unit terkait menjadi Occupied
		if assignedUnitId != nil && *assignedUnitId != "" {
			tx.Table("room_units").Where("id = ?", *assignedUnitId).Update("status", "occupied")
		}

		// 6. Ubah status katalog kamar terkait menjadi Occupied
		if err := tx.Table("rooms").Where("id = ?", roomId).Update("status", "Occupied").Error; err != nil {
			return err
		}
		return nil
	})
}

// Implementasi fungsi CheckOut
func (r *reservationRepo) CheckOut(id string, roomId string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		var currentRes Reservation
		if err := tx.Where("id = ?", id).First(&currentRes).Error; err == nil {
			// Kembalikan status unit kamar fisik menjadi Available
			if currentRes.RoomUnitId != nil && *currentRes.RoomUnitId != "" {
				tx.Table("room_units").Where("id = ?", *currentRes.RoomUnitId).Update("status", "available")
			}
		}

		// 1. Ubah status reservasi menjadi checked-out dan catat jam check-out riil
		if err := tx.Model(&Reservation{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status":             ReservationStatusCheckedOut,
			"actual_checkout_at": now,
		}).Error; err != nil {
			return err
		}

		// 2. Simpan Riwayat Log Check Out Harian
		logEntry := ReservationLog{
			Id:            time.Now().Format("20060102150405.000000"),
			ReservationId: id,
			Action:        "check_out",
			Timestamp:     now,
			Notes:         "Tamu Check Out",
			CreatedAt:     now,
		}
		if err := tx.Create(&logEntry).Error; err != nil {
			return err
		}

		// 3. Kembalikan status kamar katalog terkait menjadi Available
		if err := tx.Table("rooms").Where("id = ?", roomId).Update("status", "Available").Error; err != nil {
			return err
		}
		return nil
	})
}
