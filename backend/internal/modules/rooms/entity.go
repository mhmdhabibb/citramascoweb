package rooms

import (
	"citramascoweb-backend/internal/modules/rooms/category"
	"citramascoweb-backend/internal/modules/rooms/types"
	"time"
)

type Room struct {
	Id          string            `gorm:"type:varchar(191);primaryKey" json:"id"`
	Code        string            `json:"code"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Image       string            `json:"image"`
	CategoryId  string            `gorm:"type:varchar(191)" json:"category_id"`
	Category    category.Category `gorm:"foreignKey:CategoryId;references:Id" json:"category"`
	Price       int               `json:"price"`
	Capacity    int               `json:"capacity"`
	Size        int               `json:"size"`
	TypeId      string            `gorm:"type:varchar(191)" json:"type_id"`
	Type        types.Types       `gorm:"foreignKey:TypeId;references:Id" json:"type"`
	Description string            `json:"description"`
	Status      string            `json:"status" gorm:"-"`
	AvailabilityStatus string      `json:"availability_status" gorm:"-"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
