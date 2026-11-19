package dto

import "mime/multipart"

type CreateTypesRequest struct {
	Name  string                `json:"name" binding:"required"`
	Image *multipart.FileHeader `form:"image"`
}

type UpdateTypesRequest struct {
	Name  string                `json:"name" binding:"required"`
	Image *multipart.FileHeader `form:"image"`
}

type CreateBedTypeRequest struct {
	Name     string `json:"name" binding:"required"`
	Slug     string `json:"slug"`
	IsActive *bool  `json:"is_active"`
}

type UpdateBedTypeRequest struct {
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	IsActive *bool  `json:"is_active"`
}
