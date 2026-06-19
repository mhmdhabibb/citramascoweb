package dto

type CreateRoomRequest struct {
	Name        string  `json:"name" form:"name" binding:"required"`
	Image       string  `json:"image" form:"-"`
	CategoryId  string  `json:"category_id" form:"category_id" binding:"required"`
	Price       int     `json:"price" form:"price" binding:"required"`
	Capacity    int     `json:"capacity" form:"capacity" binding:"required"`
	Size        int     `json:"size" form:"size" binding:"required"`
	TypeId      string  `json:"type_id" form:"type_id" binding:"required"`
	Description string  `json:"description" form:"description" binding:"required"`
	IsOffer     *bool   `json:"is_offer" form:"is_offer"`
	OfferCode   *string `json:"offer_code" form:"offer_code"`
}

type UpdateRoomRequest struct {
	Name        string  `json:"name" form:"name" `
	Image       string  `json:"image" form:"-"`
	CategoryId  string  `json:"category_id" form:"category_id" `
	Price       int     `json:"price" form:"price" `
	Capacity    int     `json:"capacity" form:"capacity" `
	Size        int     `json:"size" form:"size" `
	TypeId      string  `json:"type_id" form:"type_id" `
	Description string  `json:"description" form:"description" `
	IsOffer     *bool   `json:"is_offer" form:"is_offer"`
	OfferCode   *string `json:"offer_code" form:"offer_code"`
}
