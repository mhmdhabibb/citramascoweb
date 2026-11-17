package facility

import (
	"time"

	"gorm.io/gorm"
)

type Facility struct {
	Id        string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null;unique" json:"name"`
	Icon      string         `gorm:"type:varchar(100)" json:"icon"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Facility) TableName() string {
	return "facilities"
}

type RoomTypeFacility struct {
	RoomTypeId string    `gorm:"type:varchar(191);primaryKey" json:"room_type_id"`
	FacilityId string    `gorm:"type:varchar(191);primaryKey" json:"facility_id"`
	Facility   *Facility `gorm:"foreignKey:FacilityId;references:Id" json:"facility,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

func (RoomTypeFacility) TableName() string {
	return "room_type_facilities"
}

type RoomTypeImage struct {
	Id         string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	RoomTypeId string    `gorm:"type:varchar(191);index;not null" json:"room_type_id"`
	ImageUrl   string    `gorm:"type:text;not null" json:"image_url"`
	IsPrimary  bool      `gorm:"default:false" json:"is_primary"`
	SortOrder  int       `gorm:"default:0" json:"sort_order"`
	CreatedAt  time.Time `json:"created_at"`
}

func (RoomTypeImage) TableName() string {
	return "room_type_images"
}
