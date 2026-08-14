package dto

type CreateChannelRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateChannelRequest struct {
	Name string `json:"name"`
}
