package reservation

import (
	"citramascoweb-backend/internal/dto"
	"citramascoweb-backend/internal/modules/channel"
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

type TransactionStatus string

const (
	TransactionStatusUnpaid      TransactionStatus = "unpaid"
	TransactionStatusDownPayment TransactionStatus = "down_payment"
	TransactionStatusPaid        TransactionStatus = "paid"
	TransactionStatusRefunded    TransactionStatus = "refunded"
)

type Reservation struct {
	Id                string            `gorm:"type:varchar(191);primaryKey" json:"id"`
	Code              string            `gorm:"type:varchar(191)" json:"code"`
	FullName          string            `json:"full_name"`
	Email             string            `json:"email"`
	RoomId            string            `gorm:"type:varchar(191)" json:"room_id"`
	Room              rooms.Room        `gorm:"foreignKey:RoomId;references:Id" json:"room"`
	ChannelId         *string           `gorm:"type:varchar(191)" json:"channel_id"`
	Channel           *channel.Channel  `gorm:"foreignKey:ChannelId;references:Id" json:"channel"`
	Status            ReservationStatus `gorm:"type:varchar(191);default:'pending'" json:"status"`
	TransactionStatus TransactionStatus `gorm:"type:varchar(191);default:'unpaid'" json:"transaction_status"`
	PaymentMethod     string            `gorm:"type:varchar(191);default:'bank_transfer'" json:"payment_method"`
	NumberOfAdult     int               `json:"number_of_adult"`
	NumberOfChildren  *int              `json:"number_of_children" gorm:"default:0"`

	CheckinDate  *dto.CustomDate `json:"checkin_date"`
	CheckoutDate *dto.CustomDate `json:"checkout_date"`
	Price        int             `json:"price"`
	TotalNight   int             `json:"total_night"`
	TotalPrice   int             `json:"total_price"`
	Deposit      int             `gorm:"default:0" json:"deposit"`
	IsOffer      *bool           `gorm:"type:boolean;default:false" json:"is_offer"`
	OfferCode    *string         `gorm:"type:varchar(191)" json:"offer_code"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
