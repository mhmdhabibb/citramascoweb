package notification

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepositoryInterface interface {
	GetByUser(userId string) ([]Notification, error)
	GetUnreadCount(userId string) (int64, error)
	MarkRead(id, userId string) error
	GetFinanceUserIds() ([]string, error)
	GetHousekeepingUserIds() ([]string, error)
	GetReceptionAndAdminUserIds() ([]string, error)
	CreateMany(items []Notification) error
	SaveDeviceToken(userId, token, deviceType string) error
	GetDeviceTokensByUserIds(userIds []string) ([]string, error)
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepositoryInterface {
	_ = db.AutoMigrate(&Notification{}, &UserDeviceToken{})
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
	var userIds []string
	err := r.db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.slug = ? AND roles.deleted_at IS NULL", "finance").
		Pluck("user_roles.user_id", &userIds).Error
	return userIds, err
}

func (r *notificationRepository) GetHousekeepingUserIds() ([]string, error) {
	var userIds []string
	err := r.db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.slug = ? AND roles.deleted_at IS NULL", "housekeeping").
		Pluck("user_roles.user_id", &userIds).Error
	if err != nil {
		return nil, err
	}
	if len(userIds) == 0 {
		err = r.db.Table("user_roles").
			Joins("JOIN roles ON roles.id = user_roles.role_id").
			Where("roles.slug IN ? AND roles.deleted_at IS NULL", []string{"admin", "manager", "reception"}).
			Pluck("user_roles.user_id", &userIds).Error
	}
	return userIds, err
}

func (r *notificationRepository) GetReceptionAndAdminUserIds() ([]string, error) {
	var userIds []string
	err := r.db.Table("user_roles").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("roles.slug IN ? AND roles.deleted_at IS NULL", []string{"admin", "manager", "reception"}).
		Pluck("user_roles.user_id", &userIds).Error
	return userIds, err
}

func (r *notificationRepository) CreateMany(items []Notification) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

func (r *notificationRepository) SaveDeviceToken(userId, token, deviceType string) error {
	if token == "" {
		return nil
	}
	var existing UserDeviceToken
	if err := r.db.Where("token = ?", token).First(&existing).Error; err == nil {
		return r.db.Model(&existing).Updates(map[string]interface{}{
			"user_id":     userId,
			"device_type": deviceType,
		}).Error
	}
	newToken := UserDeviceToken{
		Id:         uuid.New().String(),
		UserId:     userId,
		Token:      token,
		DeviceType: deviceType,
	}
	return r.db.Create(&newToken).Error
}

func (r *notificationRepository) GetDeviceTokensByUserIds(userIds []string) ([]string, error) {
	if len(userIds) == 0 {
		return nil, nil
	}
	var tokens []UserDeviceToken
	err := r.db.Where("user_id IN ?", userIds).Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	res := make([]string, 0, len(tokens))
	for _, t := range tokens {
		if t.Token != "" {
			res = append(res, t.Token)
		}
	}
	return res, nil
}