package sla

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type SLAService interface {
	GetPolicies() ([]SLAPolicy, error)
	GetPolicyById(id string) (*SLAPolicy, error)
	CreatePolicy(req *CreateSLAPolicyRequest) (*SLAPolicy, error)
	UpdatePolicy(id string, req *UpdateSLAPolicyRequest) (*SLAPolicy, error)
	DeletePolicy(id string) error

	GetTasks(status, taskType, assignedTo string, page, limit int) ([]OperationalTask, int64, error)
	GetTaskById(id string) (*OperationalTask, error)
	CreateTask(userId string, req *CreateOperationalTaskRequest) (*OperationalTask, error)
	AssignTask(taskId string, req *AssignTaskRequest) (*OperationalTask, error)
	UpdateTaskStatus(taskId string, req *UpdateTaskStatusRequest) (*OperationalTask, error)

	GetDailyPerformance(startDate, endDate string) ([]StaffSLAPerformanceDaily, error)
}

type slaService struct {
	repo SLARepository
}

func NewSLAService(repo SLARepository) SLAService {
	return &slaService{repo: repo}
}

func (s *slaService) GetPolicies() ([]SLAPolicy, error) {
	return s.repo.FindPolicies()
}

func (s *slaService) GetPolicyById(id string) (*SLAPolicy, error) {
	return s.repo.FindPolicyById(id)
}

func (s *slaService) CreatePolicy(req *CreateSLAPolicyRequest) (*SLAPolicy, error) {
	p := &SLAPolicy{
		Id:                      uuid.New().String(),
		Name:                    req.Name,
		Category:                req.Category,
		Priority:                req.Priority,
		TargetResponseMinutes:   req.TargetResponseMinutes,
		TargetResolutionMinutes: req.TargetResolutionMinutes,
		IsActive:                req.IsActive,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}
	if err := s.repo.CreatePolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *slaService) UpdatePolicy(id string, req *UpdateSLAPolicyRequest) (*SLAPolicy, error) {
	p, err := s.repo.FindPolicyById(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		p.Name = req.Name
	}
	if req.Category != "" {
		p.Category = req.Category
	}
	if req.Priority != "" {
		p.Priority = req.Priority
	}
	if req.TargetResponseMinutes > 0 {
		p.TargetResponseMinutes = req.TargetResponseMinutes
	}
	if req.TargetResolutionMinutes > 0 {
		p.TargetResolutionMinutes = req.TargetResolutionMinutes
	}
	p.IsActive = req.IsActive
	p.UpdatedAt = time.Now()
	if err := s.repo.UpdatePolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *slaService) DeletePolicy(id string) error {
	return s.repo.DeletePolicy(id)
}

func (s *slaService) GetTasks(status, taskType, assignedTo string, page, limit int) ([]OperationalTask, int64, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit
	return s.repo.FindTasks(status, taskType, assignedTo, limit, offset)
}

func (s *slaService) GetTaskById(id string) (*OperationalTask, error) {
	return s.repo.FindTaskById(id)
}

func (s *slaService) CreateTask(userId string, req *CreateOperationalTaskRequest) (*OperationalTask, error) {
	now := time.Now()
	taskNo := fmt.Sprintf("TSK-%s-%04d", now.Format("060102"), rand.Intn(10000))

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	task := &OperationalTask{
		Id:               uuid.New().String(),
		TaskNumber:       taskNo,
		SLAPolicyId:      req.SLAPolicyId,
		RoomUnitId:       req.RoomUnitId,
		ReservationId:    req.ReservationId,
		Title:            req.Title,
		Description:      req.Description,
		TaskType:         req.TaskType,
		Priority:         priority,
		ReportedByUserId: &userId,
		AssignedToUserId: req.AssignedToUserId,
		Status:           "pending",
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	if err := s.repo.CreateTask(task); err != nil {
		return nil, err
	}
	return s.repo.FindTaskById(task.Id)
}

func (s *slaService) AssignTask(taskId string, req *AssignTaskRequest) (*OperationalTask, error) {
	task, err := s.repo.FindTaskById(taskId)
	if err != nil {
		return nil, err
	}
	task.AssignedToUserId = &req.AssignedToUserId
	task.Status = "in_progress"
	now := time.Now()
	if task.AcknowledgedAt == nil {
		task.AcknowledgedAt = &now
		respMin := int(now.Sub(task.CreatedAt).Minutes())
		task.ActualResponseMinutes = &respMin

		if task.SLAPolicy != nil && respMin > task.SLAPolicy.TargetResponseMinutes {
			task.IsResponseBreached = true
		}
	}
	task.UpdatedAt = now
	if err := s.repo.UpdateTask(task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *slaService) UpdateTaskStatus(taskId string, req *UpdateTaskStatusRequest) (*OperationalTask, error) {
	task, err := s.repo.FindTaskById(taskId)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	task.Status = req.Status
	task.UpdatedAt = now

	if req.Status == "in_progress" && task.AcknowledgedAt == nil {
		task.AcknowledgedAt = &now
		respMin := int(now.Sub(task.CreatedAt).Minutes())
		task.ActualResponseMinutes = &respMin
		if task.SLAPolicy != nil && respMin > task.SLAPolicy.TargetResponseMinutes {
			task.IsResponseBreached = true
		}
	}

	if req.Status == "resolved" && task.ResolvedAt == nil {
		task.ResolvedAt = &now
		resMin := int(now.Sub(task.CreatedAt).Minutes())
		task.ActualResolutionMinutes = &resMin
		if task.SLAPolicy != nil && resMin > task.SLAPolicy.TargetResolutionMinutes {
			task.IsResolutionBreached = true
			if req.BreachReason != nil {
				task.BreachReason = req.BreachReason
			} else {
				defaultReason := "Waktu penanganan melebihi batas SLA yang ditentukan"
				task.BreachReason = &defaultReason
			}
		}
	}

	if req.Status == "inspected" && task.InspectedAt == nil {
		task.InspectedAt = &now
	}

	if err := s.repo.UpdateTask(task); err != nil {
		return nil, err
	}
	return s.repo.FindTaskById(taskId)
}

func (s *slaService) GetDailyPerformance(startDate, endDate string) ([]StaffSLAPerformanceDaily, error) {
	return s.repo.FindPerformanceSnapshots(startDate, endDate)
}
