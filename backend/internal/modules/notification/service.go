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

func (s *NotificationService) GetMine(userId string) ([]Notification, error) {
	return s.repo.GetByUser(userId)
}

func (s *NotificationService) GetUnreadCount(userId string) (int64, error) {
	return s.repo.GetUnreadCount(userId)
}

func (s *NotificationService) MarkRead(id, userId string) error {
	return s.repo.MarkRead(id, userId)
}