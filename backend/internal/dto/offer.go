package dto

type CreateOfferRequest struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	Image       *string `json:"image" form:"-"`
	Price       *int    `json:"price" form:"price"`
	Discount    *int    `json:"discount" form:"discount"`
	Description *string `json:"description" form:"description"`
}

type UpdateOfferRequest struct {
	Title       string  `json:"title" form:"title"`
	Image       *string `json:"image" form:"-"`
	Price       *int    `json:"price" form:"price"`
	Discount    *int    `json:"discount" form:"discount"`
	Description *string `json:"description" form:"description"`
	Status      *string `json:"status" form:"status"`
}
