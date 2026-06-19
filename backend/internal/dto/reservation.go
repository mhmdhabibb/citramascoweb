package dto

type CreateReservationRequest struct {
	RoomId        string      `json:"room_id" form:"room_id" binding:"required"`
	FullName      string      `json:"full_name" form:"full_name" binding:"required"`
	Email         string      `json:"email" form:"email" binding:"required"`
	CheckInDate   *CustomDate `json:"check_in_date" form:"check_in_date" binding:"required"`
	CheckOutDate  *CustomDate `json:"check_out_date" form:"check_out_date" binding:"required"`
	NumberOfGuest int         `json:"number_of_guest" form:"number_of_guest" binding:"required"`
}

type UpdateReservationRequest struct {
	RoomId        string      `json:"room_id" form:"room_id"`
	UserId        string      `json:"user_id" form:"user_id"`
	CheckInDate   *CustomDate `json:"check_in_date" form:"check_in_date"`
	CheckOutDate  *CustomDate `json:"check_out_date" form:"check_out_date"`
	Status        string      `json:"status" form:"status"`
	NumberOfGuest int         `json:"number_of_guest" form:"number_of_guest"`
}

type CheckAvailabilityRequest struct {
	CheckInDate  *CustomDate `form:"check_in_date" json:"check_in_date" binding:"required"`
	CheckOutDate *CustomDate `form:"check_out_date" json:"check_out_date" binding:"required"`
	RoomId       string      `form:"room_id" json:"room_id"`
}
