package reservation

import (
	"citramascoweb-backend/internal/dto"
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
	Id            string            `gorm:"type:varchar(191);primaryKey" json:"id"`
	Code          string            `gorm:"type:varchar(191)" json:"code"`
	FullName      string            `json:"full_name"`
	Email         string            `json:"email"`
	RoomId        string            `gorm:"type:varchar(191)" json:"room_id"`
	Room          rooms.Room        `gorm:"foreignKey:RoomId;references:Id" json:"room"`
	Status        ReservationStatus `gorm:"type:varchar(191);default:'pending'" json:"status"`
	NumberOfGuest int               `json:"number_of_guest"`

	CheckinDate  *dto.CustomDate `json:"checkin_date"`
	CheckoutDate *dto.CustomDate `json:"checkout_date"`
	Price        int       `json:"price"`
	TotalNight   int       `json:"total_night"`
	TotalPrice   int       `json:"total_price"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
