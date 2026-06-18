package dto

type CreateOfferRequest struct {
	Title       string  `json:"title" binding:"required"`
	Image       *string `json:"image"`
	Price       *int    `json:"price"`
	Discount    *int    `json:"discount"`
	Description *string `json:"description"`
}

type UpdateOfferRequest struct {
	Title       string  `json:"title"`
	Image       *string `json:"image"`
	Price       *int    `json:"price"`
	Discount    *int    `json:"discount"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
}
