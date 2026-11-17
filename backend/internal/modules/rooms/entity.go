package rooms

import (
	"citramascoweb-backend/internal/modules/rooms/category"
	"citramascoweb-backend/internal/modules/rooms/types"
	"time"

	"gorm.io/gorm"
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
	Code               string                 `gorm:"type:varchar(50);index" json:"code"`
	Name               string                 `gorm:"type:varchar(100);not null" json:"name"`
	Slug               string                 `gorm:"type:varchar(100);index" json:"slug"`
	Image              string                 `gorm:"type:text" json:"image"`
	CategoryId         string                 `gorm:"type:varchar(191)" json:"category_id"`
	Category           category.Category      `gorm:"foreignKey:CategoryId;references:Id" json:"category"`
	Price              int                    `json:"price"`
	BaseRate           float64                `gorm:"type:decimal(12,2);default:0" json:"base_rate"`
	Capacity           int                    `json:"capacity"` // adult
	AdultCapacity      int                    `gorm:"default:2" json:"adult_capacity"`
	ChildCapacity      *int                   `gorm:"default:0" json:"child_capacity"` // child ( 0 - 17 years old )
	Size               int                    `json:"size"`
	TypeId             string                 `gorm:"type:varchar(191)" json:"type_id"`
	Type               types.Types            `gorm:"foreignKey:TypeId;references:Id" json:"type"`
	BedTypeId          *string                `gorm:"type:varchar(191)" json:"bed_type_id"`
	Description        string                 `gorm:"type:text" json:"description"`
	IsActive           bool                   `gorm:"default:true" json:"is_active"`
	Status             RoomStatus             `json:"status"`
	AvailabilityStatus RoomAvailabilityStatus `json:"availability_status"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	DeletedAt          gorm.DeletedAt         `gorm:"index" json:"-"`
}
