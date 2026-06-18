package dto

type CreateReservationRequest struct {
	RoomId       string `json:"room_id" form:"room_id" binding:"required"`
	UserId       string `json:"user_id" form:"user_id" binding:"required"`
	CheckInDate  string `json:"check_in_date" form:"check_in_date" binding:"required"`
	CheckOutDate string `json:"check_out_date" form:"check_out_date" binding:"required"`
}

type UpdateReservationRequest struct {
	RoomId       string `json:"room_id" form:"room_id"`
	UserId       string `json:"user_id" form:"user_id"`
	CheckInDate  string `json:"check_in_date" form:"check_in_date"`
	CheckOutDate string `json:"check_out_date" form:"check_out_date"`
	Status       string `json:"status" form:"status"`
}
