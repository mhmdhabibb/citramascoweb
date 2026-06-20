package dto

type CreateReservationRequest struct {
	RoomId        string  `json:"room_id" form:"room_id" binding:"required"`
	FullName      string  `json:"full_name" form:"full_name" binding:"required"`
	Email         string  `json:"email" form:"email" binding:"required"`
	CheckInDate   string  `json:"check_in_date" form:"check_in_date" binding:"required"`
	CheckOutDate  string  `json:"check_out_date" form:"check_out_date" binding:"required"`
	NumberOfGuest int     `json:"number_of_guest" form:"number_of_guest" binding:"required"`
	IsOffer       *bool   `json:"is_offer" form:"is_offer"`
	OfferCode     *string `json:"offer_code" form:"offer_code"`
}

type UpdateReservationRequest struct {
	RoomId        string  `json:"room_id" form:"room_id"`
	UserId        string  `json:"user_id" form:"user_id"`
	CheckInDate   string  `json:"check_in_date" form:"check_in_date"`
	CheckOutDate  string  `json:"check_out_date" form:"check_out_date"`
	Status        string  `json:"status" form:"status"`
	NumberOfGuest int     `json:"number_of_guest" form:"number_of_guest"`
	IsOffer       *bool   `json:"is_offer" form:"is_offer"`
	OfferCode     *string `json:"offer_code" form:"offer_code"`
}

type CheckAvailabilityRequest struct {
	CheckInDate  string `form:"check_in_date" json:"check_in_date" binding:"required"`
	CheckOutDate string `form:"check_out_date" json:"check_out_date" binding:"required"`
	RoomId       string `form:"room_id" json:"room_id"`
}
