package dto

type CreateReservationRequest struct {
	RoomId            string  `json:"room_id" form:"room_id" binding:"required"`
	FullName          string  `json:"full_name" form:"full_name" binding:"required"`
	Email             string  `json:"email" form:"email" binding:"required"`
	CheckInDate       string  `json:"check_in_date" form:"check_in_date" binding:"required"`
	CheckOutDate      string  `json:"check_out_date" form:"check_out_date" binding:"required"`
	NumberOfAdult     int     `json:"number_of_adult" form:"number_of_adult" binding:"required"`
	NumberOfChildren  *int    `json:"number_of_children" form:"number_of_children"`
	IsOffer           *bool   `json:"is_offer" form:"is_offer"`
	OfferCode         *string `json:"offer_code" form:"offer_code"`
	Deposit           int     `json:"deposit" form:"deposit"`
	TransactionStatus string  `json:"transaction_status" form:"transaction_status"`
	PaymentMethod     string  `json:"payment_method" form:"payment_method"`
	ChannelId         *string `json:"channel_id" form:"channel_id"`
}

type UpdateReservationRequest struct {
	RoomId            string  `json:"room_id" form:"room_id"`
	UserId            string  `json:"user_id" form:"user_id"`
	CheckInDate       string  `json:"check_in_date" form:"check_in_date"`
	CheckOutDate      string  `json:"check_out_date" form:"check_out_date"`
	Status            string  `json:"status" form:"status"`
	TransactionStatus string  `json:"transaction_status" form:"transaction_status"`
	PaymentMethod     string  `json:"payment_method" form:"payment_method"`
	NumberOfAdult     int     `json:"number_of_adult" form:"number_of_adult" `
	NumberOfChildren  *int    `json:"number_of_children" form:"number_of_children"`
	IsOffer           *bool   `json:"is_offer" form:"is_offer"`
	OfferCode         *string `json:"offer_code" form:"offer_code"`
	ChannelId         *string `json:"channel_id" form:"channel_id"`
}

type CheckAvailabilityRequest struct {
	CheckInDate  string `form:"check_in_date" json:"check_in_date" binding:"required"`
	CheckOutDate string `form:"check_out_date" json:"check_out_date" binding:"required"`
	RoomId       string `form:"room_id" json:"room_id"`
}
