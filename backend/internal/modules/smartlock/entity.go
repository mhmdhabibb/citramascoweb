package smartlock

import (
	"citramascoweb-backend/internal/modules/guest"
	"citramascoweb-backend/internal/modules/rooms/unit"
	"time"
)

type RoomAccessCredential struct {
	Id                   string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	ReservationId        *string        `gorm:"type:varchar(191);index" json:"reservation_id"`
	RoomUnitId           string         `gorm:"type:varchar(191);index;not null" json:"room_unit_id"`
	RoomUnit             *unit.RoomUnit `gorm:"foreignKey:RoomUnitId;references:Id" json:"room_unit,omitempty"`
	GuestId              *string        `gorm:"type:varchar(191);index" json:"guest_id"`
	Guest                *guest.Guest   `gorm:"foreignKey:GuestId;references:Id" json:"guest,omitempty"`
	CredentialType       string         `gorm:"type:varchar(50);not null" json:"credential_type"` // passcode, card, ekey
	ExternalCredentialId *string        `gorm:"type:varchar(191)" json:"external_credential_id"`
	CredentialValue      string         `gorm:"type:varchar(191);not null" json:"credential_value"` // PIN / passcode / card UID
	ValidFrom            time.Time      `gorm:"not null" json:"valid_from"`
	ValidUntil           time.Time      `gorm:"not null" json:"valid_until"`
	Status               string         `gorm:"type:varchar(50);default:'active'" json:"status"`      // active, expired, revoked
	SyncStatus           string         `gorm:"type:varchar(50);default:'synced'" json:"sync_status"` // synced, pending, failed
	ErrorMessage         *string        `gorm:"type:text" json:"error_message"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
}

func (RoomAccessCredential) TableName() string {
	return "room_access_credentials"
}

type RoomAccessLog struct {
	Id             string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	RoomUnitId     string         `gorm:"type:varchar(191);index;not null" json:"room_unit_id"`
	RoomUnit       *unit.RoomUnit `gorm:"foreignKey:RoomUnitId;references:Id" json:"room_unit,omitempty"`
	ReservationId  *string        `gorm:"type:varchar(191);index" json:"reservation_id"`
	CredentialType string         `gorm:"type:varchar(50)" json:"credential_type"`
	UnlockedByName string         `gorm:"type:varchar(100)" json:"unlocked_by_name"`
	UnlockedAt     time.Time      `json:"unlocked_at"`
	Success        bool           `gorm:"default:true" json:"success"`
	RawPayload     string         `gorm:"type:text" json:"raw_payload"`
	CreatedAt      time.Time      `json:"created_at"`
}

func (RoomAccessLog) TableName() string {
	return "room_access_logs"
}
