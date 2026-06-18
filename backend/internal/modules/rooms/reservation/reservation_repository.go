package reservation

import "gorm.io/gorm"

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
}

type reservationRepo struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepositoryInterface {
	return &reservationRepo{db: db}
}

func (r *reservationRepo) GetAll() ([]Reservation, error) {
	var reservations []Reservation
	err := r.db.Preload("User").Preload("Room").Find(&reservations).Error
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

func (r *reservationRepo) GetById(id string) (*Reservation, error) {
	var reservation Reservation
	err := r.db.Preload("User").Preload("Room").Where("id = ?", id).First(&reservation).Error
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
