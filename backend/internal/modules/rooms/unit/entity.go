package unit

import (
	"citramascoweb-backend/internal/modules/rooms"
	"citramascoweb-backend/internal/modules/rooms/floor"
	"time"
)

type RoomUnit struct {
	Id         string      `gorm:"type:varchar(191);primaryKey" json:"id"`
	RoomNumber string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"room_number"`
	RoomId     string      `gorm:"type:varchar(191);not null" json:"room_id"`
	Room       rooms.Room  `gorm:"foreignKey:RoomId;references:Id" json:"room"`
	FloorId    string      `gorm:"type:varchar(191);not null" json:"floor_id"`
	Floor      floor.Floor `gorm:"foreignKey:FloorId;references:Id" json:"floor"`
	Status     string      `gorm:"type:varchar(50);default:'available'" json:"status"` // 'available', 'occupied', 'dirty', 'maintenance'
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}
