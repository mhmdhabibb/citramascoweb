package reservation

import (
	"citramascoweb-backend/internal/dto"
	"citramascoweb-backend/internal/modules/channel"
	"citramascoweb-backend/internal/modules/guest"
	"citramascoweb-backend/internal/modules/rooms"
	"citramascoweb-backend/internal/modules/rooms/unit"
	"time"

	"gorm.io/gorm"
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
	Code              string            `gorm:"type:varchar(191);index" json:"code"`
	BookingCode       string            `gorm:"type:varchar(191);index" json:"booking_code"`
	GuestId           *string           `gorm:"type:varchar(191);index" json:"guest_id"`
	Guest             *guest.Guest      `gorm:"foreignKey:GuestId;references:Id" json:"guest,omitempty"`
	OtaId             *string           `gorm:"type:varchar(191);index" json:"ota_id"`
	CreatedByUserId   *string           `gorm:"type:varchar(191)" json:"created_by_user_id"`
	FullName          string            `json:"full_name"`
	Email             string            `json:"email"`
	RoomId            string            `gorm:"type:varchar(191)" json:"room_id"`
	Room              rooms.Room        `gorm:"foreignKey:RoomId;references:Id" json:"room"`
	RoomUnitId        *string           `gorm:"type:varchar(191)" json:"room_unit_id"`
	RoomUnit          *unit.RoomUnit    `gorm:"foreignKey:RoomUnitId;references:Id" json:"room_unit"`
	ChannelId         *string           `gorm:"type:varchar(191)" json:"channel_id"`
	Channel           *channel.Channel  `gorm:"foreignKey:ChannelId;references:Id" json:"channel"`
	Status            ReservationStatus `gorm:"type:varchar(191);default:'pending'" json:"status"`
	TransactionStatus TransactionStatus `gorm:"type:varchar(191);default:'unpaid'" json:"transaction_status"`
	PaymentMethod     string            `gorm:"type:varchar(191);default:'bank_transfer'" json:"payment_method"`
	NumberOfAdult     int               `json:"number_of_adult"`
	NumberOfChildren  *int              `json:"number_of_children" gorm:"default:0"`

	CheckinDate    *dto.CustomDate `json:"checkin_date"`
	CheckoutDate   *dto.CustomDate `json:"checkout_date"`
	Price          int             `json:"price"`
	TotalNight     int             `json:"total_night"`
	TotalPrice     int             `json:"total_price"`
	TotalAmount    float64         `gorm:"type:decimal(12,2);default:0" json:"total_amount"`
	DiscountAmount float64         `gorm:"type:decimal(12,2);default:0" json:"discount_amount"`
	TaxAmount      float64         `gorm:"type:decimal(12,2);default:0" json:"tax_amount"`
	NetAmount      float64         `gorm:"type:decimal(12,2);default:0" json:"net_amount"`
	LeadTimeDays   int             `gorm:"default:0" json:"lead_time_days"`
	Deposit        int             `gorm:"default:0" json:"deposit"`
	IsOffer        *bool           `gorm:"type:boolean;default:false" json:"is_offer"`
	OfferCode      *string         `gorm:"type:varchar(191)" json:"offer_code"`

	IsEarlyCheckin       *bool      `gorm:"type:boolean;default:false" json:"is_early_checkin"`
	EstimatedArrivalTime *string    `gorm:"type:varchar(50)" json:"estimated_arrival_time"`
	ActualCheckinAt      *time.Time `json:"actual_checkin_at"`
	ActualCheckoutAt     *time.Time `json:"actual_checkout_at"`

	CancellationDate   *time.Time `json:"cancellation_date"`
	CancellationReason *string    `gorm:"type:text" json:"cancellation_reason"`

	Items []ReservationItem `gorm:"foreignKey:ReservationId" json:"items,omitempty"`
	Logs  []ReservationLog  `gorm:"foreignKey:ReservationId" json:"logs,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type ReservationItem struct {
	Id            string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	ReservationId string    `gorm:"type:varchar(191);index;not null" json:"reservation_id"`
	RoomUnitId    string    `gorm:"type:varchar(191);index;not null" json:"room_unit_id"`
	StayDate      time.Time `gorm:"type:date;not null" json:"stay_date"`
	RateApplied   float64   `gorm:"type:decimal(12,2);not null" json:"rate_applied"`
	CreatedAt     time.Time `json:"created_at"`
}

func (ReservationItem) TableName() string {
	return "reservation_items"
}

type ReservationLog struct {
	Id            string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	ReservationId string    `gorm:"type:varchar(191);index" json:"reservation_id"`
	Action        string    `gorm:"type:varchar(50)" json:"action"` // "check_in" | "check_out"
	Timestamp     time.Time `json:"timestamp"`
	IsEarly       bool      `gorm:"default:false" json:"is_early"`
	Notes         string    `json:"notes"`
	HandledById   *string   `gorm:"type:varchar(191)" json:"handled_by_id"`
	CreatedAt     time.Time `json:"created_at"`
}
