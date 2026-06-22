package rooms

import (
	"citramascoweb-backend/internal/modules/rooms/category"
	"citramascoweb-backend/internal/modules/rooms/types"
	"time"
)

type RoomStatus string

const (
	RoomStatusActive      RoomStatus = "active"
	RoomStatusMaintenance RoomStatus = "maintenance"
	RoomStatusDirty       RoomStatus = "dirty"
)

// Define Enum untuk RoomAvailabilityStatus (Keterisian)
type RoomAvailabilityStatus string

const (
	RoomAvailabilityAvailable RoomAvailabilityStatus = "available"
	RoomAvailabilityOccupied  RoomAvailabilityStatus = "occupied"
)

type Room struct {
	Id                 string                 `gorm:"type:varchar(191);primaryKey" json:"id"`
	Code               string                 `json:"code"`
	Name               string                 `json:"name"`
	Slug               string                 `json:"slug"`
	Image              string                 `json:"image"`
	CategoryId         string                 `gorm:"type:varchar(191)" json:"category_id"`
	Category           category.Category      `gorm:"foreignKey:CategoryId;references:Id" json:"category"`
	Price              int                    `json:"price"`
	Capacity           int                    `json:"capacity"`
	Size               int                    `json:"size"`
	TypeId             string                 `gorm:"type:varchar(191)" json:"type_id"`
	Type               types.Types            `gorm:"foreignKey:TypeId;references:Id" json:"type"`
	Description        string                 `json:"description"`
	Status             RoomStatus             `json:"status"`
	AvailabilityStatus RoomAvailabilityStatus `json:"availability_status"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}
