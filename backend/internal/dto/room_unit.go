package dto

type CreateRoomUnitRequest struct {
	RoomNumber   string  `json:"room_number" binding:"required"`
	RoomId       string  `json:"room_id" binding:"required"`
	FloorId      string  `json:"floor_id" binding:"required"`
	RoomStatusId *string `json:"room_status_id"`
	Status       string  `json:"status"`
}

type UpdateRoomUnitRequest struct {
	RoomNumber   string  `json:"room_number"`
	RoomId       string  `json:"room_id"`
	FloorId      string  `json:"floor_id"`
	RoomStatusId *string `json:"room_status_id"`
	Status       string  `json:"status"`
}

type UpdateRoomUnitStatusRequest struct {
	Status       string  `json:"status"`
	RoomStatusId *string `json:"room_status_id"`
}

type CreateRoomStatusRequest struct {
	Code        string `json:"code" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateRoomStatusRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsActive    *bool  `json:"is_active"`
}
