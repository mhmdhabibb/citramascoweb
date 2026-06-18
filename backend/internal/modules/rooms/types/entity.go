package types

import "time"

type Types struct {
	Id        string    `json:"id" gorm:"type:varchar(191);primaryKey"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
