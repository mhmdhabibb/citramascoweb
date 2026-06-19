package types

import "time"

type TypeRoom struct {
	Id          string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Image       string    `json:"image"`
	CategoryId  string    `gorm:"type:varchar(191)" json:"category_id"`
	Price       int       `json:"price"`
	Capacity    int       `json:"capacity"`
	Size        int       `json:"size"`
	TypeId      string    `gorm:"type:varchar(191)" json:"type_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (TypeRoom) TableName() string {
	return "rooms"
}

type Types struct {
	Id        string     `json:"id" gorm:"type:varchar(191);primaryKey"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Rooms     []TypeRoom `gorm:"foreignKey:TypeId" json:"rooms,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}
