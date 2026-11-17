package unit

import (
	"citramascoweb-backend/internal/modules/rooms"
	"citramascoweb-backend/internal/modules/rooms/floor"
	"time"

	"gorm.io/gorm"
)

type RoomStatus struct {
	Id          string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(50);not null;unique" json:"code"` // VC, VD, OC, OD, OOO, OOS
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`       // Vacant Clean, Vacant Dirty, etc.
	Description string         `gorm:"type:text" json:"description"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (RoomStatus) TableName() string {
	return "room_status"
}

type RoomUnit struct {
	Id               string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	RoomNumber       string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"room_number"`
	RoomId           string         `gorm:"type:varchar(191);not null" json:"room_id"`
	Room             rooms.Room     `gorm:"foreignKey:RoomId;references:Id" json:"room"`
	FloorId          string         `gorm:"type:varchar(191);not null" json:"floor_id"`
	Floor            floor.Floor    `gorm:"foreignKey:FloorId;references:Id" json:"floor"`
	RoomStatusId     *string        `gorm:"type:varchar(191);index" json:"room_status_id"`
	RoomStatus       *RoomStatus    `gorm:"foreignKey:RoomStatusId;references:Id" json:"room_status,omitempty"`
	Status           string         `gorm:"type:varchar(50);default:'available'" json:"status"` // backward compatibility
	LockProvider     *string        `gorm:"type:varchar(100)" json:"lock_provider"`             // ttlock, tuya, yale
	LockName         *string        `gorm:"type:varchar(100)" json:"lock_name"`
	LockMacAddress   *string        `gorm:"type:varchar(100)" json:"lock_mac_address"`
	GatewayId        *string        `gorm:"type:varchar(100)" json:"gateway_id"`
	BatteryLevel     *int           `gorm:"default:100" json:"battery_level"`
	BatteryUpdatedAt *time.Time     `json:"battery_updated_at"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
