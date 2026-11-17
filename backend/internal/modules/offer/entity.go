package offer

import (
	"citramascoweb-backend/internal/dto"
	"time"

	"gorm.io/gorm"
)

type OfferStatus string

const (
	OfferStatusActive    OfferStatus = "active"
	OfferStatusArchieved OfferStatus = "archieved"
	OfferStatusDraft     OfferStatus = "draft"
)

type Offer struct {
	Id            string          `json:"id" gorm:"primaryKey"`
	Title         string          `json:"title"`
	SubTitle      string          `json:"sub_title"`
	Image         string          `json:"image"`
	Status        OfferStatus     `gorm:"default:'draft'" json:"status"`
	DiscountType  string          `gorm:"type:varchar(50);default:'percentage'" json:"discount_type"` // percentage, fixed
	DiscountValue float64         `gorm:"type:decimal(12,2);default:0" json:"discount_value"`
	MinNights     int             `gorm:"default:1" json:"min_nights"`
	Price         *int            `json:"price"`
	Discount      *int            `json:"discount"`
	Code          string          `json:"code"`
	Discounted    *int            `json:"discounteed"`
	ValidStart    *dto.CustomDate `json:"valid_start"`
	ValidEnd      *dto.CustomDate `json:"valid_end"`
	MaxQuota      *int            `json:"max_quota"`
	Description   *string         `json:"description"`
	IsActive      bool            `gorm:"default:true" json:"is_active"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

type OfferUsage struct {
	Id              string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	ReservationId   string    `gorm:"type:varchar(191);index;not null" json:"reservation_id"`
	GuestId         *string   `gorm:"type:varchar(191);index" json:"guest_id"`
	OfferId         string    `gorm:"type:varchar(191);index;not null" json:"offer_id"`
	DiscountApplied float64   `gorm:"type:decimal(12,2);not null" json:"discount_applied"`
	CreatedAt       time.Time `json:"created_at"`
}

func (OfferUsage) TableName() string {
	return "offer_usages"
}
