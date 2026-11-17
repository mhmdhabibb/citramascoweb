package service_request

import (
	"citramascoweb-backend/config"
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
	CreateRequest(dto *CreateServiceRequestDto, userId string) (*ServiceRequest, error)
	GetAll(status string) ([]ServiceRequest, error)
	GetById(id string) (*ServiceRequest, error)
	AssignToHousekeeping(id string, dto *AssignTaskDto) error
	CompleteRequest(id string, dto *CompleteTaskDto, userId string) error
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

func (s *serviceRequestService) CreateRequest(dto *CreateServiceRequestDto, userId string) (*ServiceRequest, error) {
	room, err := s.roomRepo.GetById(dto.RoomId)
	if err != nil {
		return nil, errors.New("room not found")
	}

	priority := dto.Priority
	if priority == "" {
		if dto.Category == CategoryBrokenItem {
			priority = PriorityUrgent
		} else {
			priority = PriorityMedium
		}
	}

	var reportedBy *string
	if userId != "" {
		reportedBy = &userId
	}

	req := &ServiceRequest{
		Id:               uuid.New().String(),
		RoomId:           dto.RoomId,
		ReservationId:    dto.ReservationId,
		GuestName:        dto.GuestName,
		GuestPhone:       dto.GuestPhone,
		Category:         dto.Category,
		Title:            dto.Title,
		Description:      dto.Description,
		Priority:         priority,
		ReportedByUserId: reportedBy,
		Status:           StatusPendingReception,
		DamageCharge:     0,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.Create(req); err != nil {
		return nil, err
	}

	// Otomatis kirim notifikasi ke Frontdesk & Admin via Notification & FCM
	if s.notifService != nil {
		title := fmt.Sprintf("🚨 Guest Report (%s): %s", room.Name, req.Title)
		msg := fmt.Sprintf("Guest %s in room %s reported: %s. Please review and assign immediately.", req.GuestName, room.Name, req.Description)
		_ = s.notifService.NotifyRoomCleaned(room.Name, room.Code, req.Id)
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
		return errors.New("report not found")
	}

	if err := s.repo.AssignTask(id, dto.AssignedToUserId, dto.Notes, dto.Priority); err != nil {
		return err
	}

	// Otomatis kirim notifikasi ke seluruh staf HOUSEKEEPING (HP Android + Web)
	if s.notifService != nil {
		roomName := "Room"
		if item.Room.Name != "" {
			roomName = item.Room.Name
		}
		_ = s.notifService.NotifyDirtyRoom(roomName, item.Room.Code, item.RoomId)
	}

	return nil
}

func (s *serviceRequestService) CompleteRequest(id string, dto *CompleteTaskDto, userId string) error {
	item, err := s.repo.GetById(id)
	if err != nil {
		return errors.New("report not found")
	}

	charge := 0
	if dto != nil {
		charge = dto.DamageCharge
	}

	now := time.Now()
	resMinutes := int(now.Sub(item.CreatedAt).Minutes())

	targetMinutes := 60
	switch item.Priority {
	case PriorityUrgent:
		targetMinutes = 30
	case PriorityHigh:
		targetMinutes = 45
	case PriorityMedium:
		targetMinutes = 60
	case PriorityLow:
		targetMinutes = 120
	}

	isBreached := resMinutes > targetMinutes
	breachReason := ""
	if isBreached {
		breachReason = fmt.Sprintf("Waktu penanganan (%d menit) melebihi batas target SLA (%d menit)", resMinutes, targetMinutes)
	}

	if err := s.repo.CompleteTask(id, charge, resMinutes, isBreached, breachReason); err != nil {
		return err
	}

	// Update SLA Performance harian untuk staf yang mengerjakan
	staffId := ""
	if item.AssignedToUserId != nil && *item.AssignedToUserId != "" {
		staffId = *item.AssignedToUserId
	} else if userId != "" {
		staffId = userId
	}

	if staffId != "" {
		s.updateStaffDailySLA(staffId, isBreached, resMinutes)
	}

	// Otomatis kirim notifikasi ke Frontdesk bahwa tugas telah selesai
	if s.notifService != nil {
		roomName := "Room"
		if item.Room.Name != "" {
			roomName = item.Room.Name
		}
		_ = s.notifService.NotifyRoomCleaned(roomName, item.Room.Code, item.RoomId)
	}

	return nil
}

func (s *serviceRequestService) updateStaffDailySLA(staffId string, isBreached bool, resMinutes int) {
	if config.DB == nil {
		return
	}
	loc := time.FixedZone("WIB", 7*3600)
	today := time.Now().In(loc).Truncate(24 * time.Hour)

	type DailyPerf struct {
		Id                       string
		TotalTasksAssigned       int
		TotalTasksCompleted      int
		TotalBreachedTasks       int
		SLAComplianceRate        float64
		AvgResolutionTimeMinutes float64
	}

	var record DailyPerf
	err := config.DB.Table("staff_sla_performance_daily").
		Where("snapshot_date = ? AND staff_id = ?", today, staffId).
		First(&record).Error

	if err != nil {
		breachedCount := 0
		compliance := 100.0
		if isBreached {
			breachedCount = 1
			compliance = 0.0
		}
		config.DB.Table("staff_sla_performance_daily").Create(map[string]interface{}{
			"id":                          uuid.New().String(),
			"snapshot_date":               today,
			"staff_id":                    staffId,
			"department":                  "housekeeping",
			"total_tasks_assigned":        1,
			"total_tasks_completed":       1,
			"total_breached_tasks":        breachedCount,
			"sla_compliance_rate":         compliance,
			"avg_resolution_time_minutes": float64(resMinutes),
			"created_at":                  time.Now(),
		})
	} else {
		completed := record.TotalTasksCompleted + 1
		breached := record.TotalBreachedTasks
		if isBreached {
			breached++
		}
		compliance := float64(completed-breached) / float64(completed) * 100.0
		avgRes := (record.AvgResolutionTimeMinutes*float64(record.TotalTasksCompleted) + float64(resMinutes)) / float64(completed)

		config.DB.Table("staff_sla_performance_daily").
			Where("id = ?", record.Id).
			Updates(map[string]interface{}{
				"total_tasks_completed":       completed,
				"total_breached_tasks":        breached,
				"sla_compliance_rate":         compliance,
				"avg_resolution_time_minutes": avgRes,
			})
	}
}

func (s *serviceRequestService) CancelRequest(id string) error {
	return s.repo.Cancel(id)
}
