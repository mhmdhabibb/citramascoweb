package dto

type CreateFloorRequest struct {
	Name        string `json:"name" binding:"required"`
	FloorNumber int    `json:"floor_number" binding:"required"`
	Description string `json:"description"`
}

type UpdateFloorRequest struct {
	Name        string `json:"name"`
	FloorNumber int    `json:"floor_number"`
	Description string `json:"description"`
}
