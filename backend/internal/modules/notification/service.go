package notification

import (
	"fmt"

	"github.com/google/uuid"
)

// NotificationService is the exported service used to notify finance users
// about newly recorded deposits and to serve a user's own notification feed.
type NotificationService struct {
	repo NotificationRepositoryInterface
}

func NewNotificationService(repo NotificationRepositoryInterface) *NotificationService {
	return &NotificationService{repo: repo}
}

// NotifyFinanceDeposit creates an unread notification for every user with the
// "finance" role when a deposit is recorded against a reservation.
func (s *NotificationService) NotifyFinanceDeposit(deposit int, guestName, guestCode, resId string) error {
	ids, err := s.repo.GetFinanceUserIds()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	items := make([]Notification, 0, len(ids))
	for _, uid := range ids {
		items = append(items, Notification{
			Id:          uuid.New().String(),
			UserId:      uid,
			Type:        "deposit",
			Title:       "New Deposit Received",
			Message:     fmt.Sprintf("Deposit of Rp%d received for guest %s (booking %s).", deposit, guestName, guestCode),
			ReferenceId: resId,
			IsRead:      false,
		})
	}
	return s.repo.CreateMany(items)
}

// NotifyDirtyRoom sends notifications to housekeeping and staff when a room is marked dirty
func (s *NotificationService) NotifyDirtyRoom(roomName, roomCode, roomId string) error {
	ids, err := s.repo.GetHousekeepingUserIds()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	title := "🧹 Permintaan Pembersihan Kamar"
	msg := fmt.Sprintf("Kamar %s (%s) diset Dirty dan memerlukan pembersihan segera dari tim Housekeeping.", roomName, roomCode)

	items := make([]Notification, 0, len(ids))
	for _, uid := range ids {
		items = append(items, Notification{
			Id:          uuid.New().String(),
			UserId:      uid,
			Type:        "cleaning_task",
			Title:       title,
			Message:     msg,
			ReferenceId: roomId,
			IsRead:      false,
		})
	}
	_ = s.repo.CreateMany(items)

	// Kirim Push Notification ke HP Android & Browser Housekeeping
	go func() {
		tokens, err := s.repo.GetDeviceTokensByUserIds(ids)
		if err == nil && len(tokens) > 0 {
			_ = SendFCMPush(tokens, title, msg, map[string]interface{}{
				"type":    "cleaning_task",
				"room_id": roomId,
				"status":  "dirty",
			})
		}
	}()

	return nil
}

// NotifyRoomCleaned sends notifications to frontdesk and admin when a room is cleaned
func (s *NotificationService) NotifyRoomCleaned(roomName, roomCode, roomId string) error {
	ids, err := s.repo.GetReceptionAndAdminUserIds()
	if err != nil {
		return err
	}
	if len(ids) == 0 {
		return nil
	}

	title := "✨ Kamar Selesai Dibersihkan"
	msg := fmt.Sprintf("Kamar %s (%s) telah selesai dibersihkan oleh Housekeeping dan status kembali Available.", roomName, roomCode)

	items := make([]Notification, 0, len(ids))
	for _, uid := range ids {
		items = append(items, Notification{
			Id:          uuid.New().String(),
			UserId:      uid,
			Type:        "room_cleaned",
			Title:       title,
			Message:     msg,
			ReferenceId: roomId,
			IsRead:      false,
		})
	}
	_ = s.repo.CreateMany(items)

	// Kirim Push Notification ke HP Android & Browser Frontdesk/Admin
	go func() {
		tokens, err := s.repo.GetDeviceTokensByUserIds(ids)
		if err == nil && len(tokens) > 0 {
			_ = SendFCMPush(tokens, title, msg, map[string]interface{}{
				"type":    "room_cleaned",
				"room_id": roomId,
				"status":  "available",
			})
		}
	}()

	return nil
}

func (s *NotificationService) SaveDeviceToken(userId, token, deviceType string) error {
	return s.repo.SaveDeviceToken(userId, token, deviceType)
}

func (s *NotificationService) GetMine(userId string) ([]Notification, error) {
	return s.repo.GetByUser(userId)
}

func (s *NotificationService) GetUnreadCount(userId string) (int64, error) {
	return s.repo.GetUnreadCount(userId)
}

func (s *NotificationService) MarkRead(id, userId string) error {
	return s.repo.MarkRead(id, userId)
}