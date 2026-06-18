package category

import "time"

type Category struct {
	Id        string     `gorm:"primaryKey" json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	IsDeleted bool       `json:"is_deleted"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
