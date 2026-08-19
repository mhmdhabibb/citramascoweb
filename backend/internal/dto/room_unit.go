package dto

type CreateRoomUnitRequest struct {
	RoomNumber string `json:"room_number" binding:"required"`
	RoomId     string `json:"room_id" binding:"required"`
	FloorId    string `json:"floor_id" binding:"required"`
	Status     string `json:"status"`
}

type UpdateRoomUnitRequest struct {
	RoomNumber string `json:"room_number"`
	RoomId     string `json:"room_id"`
	FloorId    string `json:"floor_id"`
	Status     string `json:"status"`
}

type UpdateRoomUnitStatusRequest struct {
	Status string `json:"status" binding:"required"`
}
