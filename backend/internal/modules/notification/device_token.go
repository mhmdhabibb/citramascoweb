package notification

import (
	"time"
)

type UserDeviceToken struct {
	Id         string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	UserId     string    `gorm:"type:varchar(191);index" json:"user_id"`
	Token      string    `gorm:"type:text;uniqueIndex" json:"token"`
	DeviceType string    `gorm:"type:varchar(50);default:'web'" json:"device_type"` // 'android', 'web', 'ios'
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
