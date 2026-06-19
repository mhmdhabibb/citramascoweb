package category

import "time"

type CategoryRoom struct {
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

func (CategoryRoom) TableName() string {
	return "rooms"
}

type Category struct {
	Id        string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	Name      string         `json:"name"`
	Slug      string         `json:"slug"`
	IsDeleted bool           `json:"is_deleted"`
	Rooms     []CategoryRoom `gorm:"foreignKey:CategoryId" json:"rooms,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at"`
}
