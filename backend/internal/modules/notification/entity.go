package notification

import "time"

// Notification is a per-user in-app notification (e.g. a deposit sent to finance users).
type Notification struct {
	Id          string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	UserId      string    `gorm:"type:varchar(191);index" json:"user_id"`
	Type        string    `gorm:"type:varchar(100)" json:"type"`
	Title       string    `json:"title"`
	Message     string    `json:"message"`
	ReferenceId string    `gorm:"type:varchar(191)" json:"reference_id"`
	IsRead      bool      `gorm:"default:false" json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}