package channel

import "time"

type Channel struct {
	Id        string    `json:"id" gorm:"primaryKey"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
