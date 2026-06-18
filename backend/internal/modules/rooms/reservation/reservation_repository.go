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
}

type reservationRepo struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepositoryInterface {
	return &reservationRepo{db: db}
}

func (r *reservationRepo) GetAll() ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("Room").Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepo) GetById(id string) (*Reservation, error) {
	var reservation Reservation
	err := r.db.Preload("Room").Where("id = ?", id).First(&reservation).Error
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
	err := r.db.Preload("User").Preload("Room").Where("status = ?", status).Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepo) ApproveReservation(id string) error {
	return r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", ReservationStatusApproved).Error
}

func (r *reservationRepo) CancelReservation(id string) error {
	return r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", ReservationStatusCancel).Error
}

func (r *reservationRepo) RejectReservation(id string) error {
	return r.db.Model(&Reservation{}).Where("id = ?", id).Update("status", ReservationStatusRejected).Error
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
		Find(&availableRooms).Error

	if err != nil {
		return nil, err
	}
	return availableRooms, nil
}
