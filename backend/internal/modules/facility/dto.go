package facility

type CreateFacilityRequest struct {
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon"`
}

type UpdateFacilityRequest struct {
	Name string `json:"name" binding:"required"`
	Icon string `json:"icon"`
}

type AssignFacilitiesRequest struct {
	FacilityIds []string `json:"facility_ids" binding:"required"`
}

type CreateRoomTypeImageRequest struct {
	ImageUrl  string `json:"image_url" binding:"required"`
	IsPrimary bool   `json:"is_primary"`
	SortOrder int    `json:"sort_order"`
}
