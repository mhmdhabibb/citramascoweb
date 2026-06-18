package types

import "time"

type Types struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
