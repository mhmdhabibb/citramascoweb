package offer

import "time"

type OfferStatus string

const (
	OfferStatusActive    OfferStatus = "active"
	OfferStatusArchieved OfferStatus = "archieved"
	OfferStatusDraft     OfferStatus = "draft"
)

type Offer struct {
	Id          string      `json:"id" gorm:"primaryKey"`
	Title       string      `json:"title"`
	Image       string      `json:"image"`
	Status      OfferStatus `gorm:"default:'draft'"`
	Price       *int        `json:"price"`
	Discount    *int        `json:"discount"`
	Code        string      `json:"code"`
	Discounted  *int        `json:"discounteed"`
	ValidStart  *time.Time  `json:"valid_start"`
	ValidEnd    *time.Time  `json:"valid_end"`
	MaxQuota    *int        `json:"max_quota"`
	Description *string     `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
