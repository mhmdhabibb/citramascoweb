package reservation

import (
	"citramascoweb-backend/internal/modules/auth"
	"citramascoweb-backend/internal/modules/rooms"
	"time"
)

type ReservationStatus string

const (
	ReservationStatusPending    ReservationStatus = "pending"
	ReservationStatusApproved   ReservationStatus = "approved"
	ReservationStatusRejected   ReservationStatus = "rejected"
	ReservationStatusCancel     ReservationStatus = "cancel"
	ReservationStatusCheckedIn  ReservationStatus = "checked-in"
	ReservationStatusCheckedOut ReservationStatus = "checked-out"
)

type Reservation struct {
	Id     string            `gorm:"type:varchar(191);primaryKey" json:"id"`
	Code   string            `gorm:"type:varchar(191)" json:"code"`    // Ditambahkan varchar(191) agar aman di-index
	UserId string            `gorm:"type:varchar(191)" json:"user_id"` // FIX: Ditambahkan varchar(191)
	User   auth.User         `gorm:"foreignKey:UserId;references:Id" json:"user"`
	RoomId string            `gorm:"type:varchar(191)" json:"room_id"` // FIX: Ditambahkan varchar(191)
	Room   rooms.Room        `gorm:"foreignKey:RoomId;references:Id" json:"room"`
	Status ReservationStatus `gorm:"type:varchar(191);default:'pending'" json:"status"` // Ditambahkan tipe varchar

	CheckinDate  time.Time `json:"checkin_date"`
	CheckoutDate time.Time `json:"checkout_date"`
	Price        int       `json:"price"`
	TotalNight   int       `json:"total_night"`
	TotalPrice   int       `json:"total_price"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
