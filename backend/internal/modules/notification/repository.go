package notification

import (
	"citramascoweb-backend/internal/modules/auth"

	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	GetByUser(userId string) ([]Notification, error)
	GetUnreadCount(userId string) (int64, error)
	MarkRead(id, userId string) error
	GetFinanceUserIds() ([]string, error)
	CreateMany(items []Notification) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepositoryInterface {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) GetByUser(userId string) ([]Notification, error) {
	var items []Notification
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&items).Error
	return items, err
}

func (r *notificationRepository) GetUnreadCount(userId string) (int64, error) {
	var count int64
	err := r.db.Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userId, false).
		Count(&count).Error
	return count, err
}

func (r *notificationRepository) MarkRead(id, userId string) error {
	return r.db.Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userId).
		Update("is_read", true).Error
}

func (r *notificationRepository) GetFinanceUserIds() ([]string, error) {
	var users []auth.User
	if err := r.db.Select("id").Where("role = ?", "finance").Find(&users).Error; err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(users))
	for _, u := range users {
		ids = append(ids, u.Id)
	}
	return ids, nil
}

func (r *notificationRepository) CreateMany(items []Notification) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}