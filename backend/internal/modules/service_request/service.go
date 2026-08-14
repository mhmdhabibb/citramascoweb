package service_request

import (
	"citramascoweb-backend/internal/modules/notification"
	"citramascoweb-backend/internal/modules/rooms"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateServiceRequestDto struct {
	RoomId        string                 `json:"room_id" binding:"required"`
	ReservationId *string                `json:"reservation_id,omitempty"`
	GuestName     string                 `json:"guest_name" binding:"required"`
	GuestPhone    string                 `json:"guest_phone,omitempty"`
	Category      ServiceRequestCategory `json:"category" binding:"required"`
	Title         string                 `json:"title" binding:"required"`
	Description   string                 `json:"description" binding:"required"`
	Priority      ServiceRequestPriority `json:"priority,omitempty"`
}

type AssignTaskDto struct {
	AssignedToUserId string `json:"assigned_to_user_id,omitempty"`
	Notes            string `json:"notes,omitempty"`
	Priority         string `json:"priority,omitempty"`
}

type CompleteTaskDto struct {
	DamageCharge int `json:"damage_charge"`
}

type ServiceRequestServiceInterface interface {
	CreateRequest(dto *CreateServiceRequestDto) (*ServiceRequest, error)
	GetAll(status string) ([]ServiceRequest, error)
	GetById(id string) (*ServiceRequest, error)
	AssignToHousekeeping(id string, dto *AssignTaskDto) error
	CompleteRequest(id string, dto *CompleteTaskDto) error
	CancelRequest(id string) error
}

type serviceRequestService struct {
	repo         ServiceRequestRepositoryInterface
	roomRepo     rooms.RoomRepositoryInterface
	notifService *notification.NotificationService
}

func NewServiceRequestService(
	repo ServiceRequestRepositoryInterface,
	roomRepo rooms.RoomRepositoryInterface,
	notifService *notification.NotificationService,
) ServiceRequestServiceInterface {
	return &serviceRequestService{
		repo:         repo,
		roomRepo:     roomRepo,
		notifService: notifService,
	}
}

func (s *serviceRequestService) CreateRequest(dto *CreateServiceRequestDto) (*ServiceRequest, error) {
	room, err := s.roomRepo.GetById(dto.RoomId)
	if err != nil {
		return nil, errors.New("kamar tidak ditemukan")
	}

	priority := dto.Priority
	if priority == "" {
		if dto.Category == CategoryBrokenItem {
			priority = PriorityUrgent
		} else {
			priority = PriorityMedium
		}
	}

	req := &ServiceRequest{
		Id:            uuid.New().String(),
		RoomId:        dto.RoomId,
		ReservationId: dto.ReservationId,
		GuestName:     dto.GuestName,
		GuestPhone:    dto.GuestPhone,
		Category:      dto.Category,
		Title:         dto.Title,
		Description:   dto.Description,
		Priority:      priority,
		Status:        StatusPendingReception,
		DamageCharge:  0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.repo.Create(req); err != nil {
		return nil, err
	}

	// Otomatis kirim notifikasi ke Frontdesk & Admin via Notification & FCM
	if s.notifService != nil {
		title := fmt.Sprintf("🚨 Laporan Tamu (%s): %s", room.Name, req.Title)
		msg := fmt.Sprintf("Tamu %s di kamar %s melaporkan: %s. Mohon segera ditinjau dan ditugaskan.", req.GuestName, room.Name, req.Description)
		_ = s.notifService.NotifyRoomCleaned(room.Name, room.Code, req.Id) // Trigger push ke Frontdesk
		go func() {
			_ = notification.SendFCMPush(nil, title, msg, map[string]interface{}{
				"type":       "guest_service_request",
				"request_id": req.Id,
				"room_id":    room.Id,
			})
		}()
	}

	return req, nil
}

func (s *serviceRequestService) GetAll(status string) ([]ServiceRequest, error) {
	return s.repo.GetAll(status)
}

func (s *serviceRequestService) GetById(id string) (*ServiceRequest, error) {
	return s.repo.GetById(id)
}

func (s *serviceRequestService) AssignToHousekeeping(id string, dto *AssignTaskDto) error {
	item, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("laporan tidak ditemukan")
	}

	if err := s.repo.AssignTask(id, dto.AssignedToUserId, dto.Notes, dto.Priority); err != nil {
		return err
	}

	// Otomatis kirim notifikasi ke seluruh staf HOUSEKEEPING (HP Android + Web)
	if s.notifService != nil {
		roomName := "Kamar"
		if item.Room.Name != "" {
			roomName = item.Room.Name
		}
		_ = s.notifService.NotifyDirtyRoom(roomName, item.Room.Code, item.RoomId)
	}

	return nil
}

func (s *serviceRequestService) CompleteRequest(id string, dto *CompleteTaskDto) error {
	item, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("laporan tidak ditemukan")
	}

	charge := 0
	if dto != nil {
		charge = dto.DamageCharge
	}

	if err := s.repo.CompleteTask(id, charge); err != nil {
		return err
	}

	// Otomatis kirim notifikasi ke Frontdesk bahwa tugas telah selesai
	if s.notifService != nil {
		roomName := "Kamar"
		if item.Room.Name != "" {
			roomName = item.Room.Name
		}
		_ = s.notifService.NotifyRoomCleaned(roomName, item.Room.Code, item.RoomId)
	}

	return nil
}

func (s *serviceRequestService) CancelRequest(id string) error {
	return s.repo.Cancel(id)
}
