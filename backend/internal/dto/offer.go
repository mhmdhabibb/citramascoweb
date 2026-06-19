package dto

import "time"

type CreateOfferRequest struct {
	Title       string     `json:"title" form:"title" binding:"required"`
	Image       *string    `json:"image" form:"-"`
	Price       *int       `json:"price" form:"price"`
	Discount    *int       `json:"discount" form:"discount"`
	Code        string     `json:"code"`
	ValidStart  *time.Time `json:"valid_start" form:"valid_start"`
	ValidEnd    *time.Time `json:"valid_end" form:"valid_end"`
	MaxQuota    *int       `json:"max_quota" form:"max_quota"`
	Description *string    `json:"description" form:"description"`
}

type UpdateOfferRequest struct {
	Title       string     `json:"title" form:"title"`
	Image       *string    `json:"image" form:"-"`
	Price       *int       `json:"price" form:"price"`
	Discount    *int       `json:"discount" form:"discount"`
	ValidStart  *time.Time `json:"valid_start" form:"valid_start"`
	ValidEnd    *time.Time `json:"valid_end" form:"valid_end"`
	MaxQuota    *int       `json:"max_quota" form:"max_quota"`
	Description *string    `json:"description" form:"description"`
	Status      *string    `json:"status" form:"status"`
}
