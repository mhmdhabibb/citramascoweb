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
