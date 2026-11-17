package floor

import (
	"time"

	"gorm.io/gorm"
)

type Floor struct {
	Id          string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	FloorNumber int            `gorm:"not null" json:"floor_number"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
