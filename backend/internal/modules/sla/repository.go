package sla

import (
	"errors"

	"gorm.io/gorm"
)

type SLARepository interface {
	// Policies
	FindPolicies() ([]SLAPolicy, error)
	FindPolicyById(id string) (*SLAPolicy, error)
	CreatePolicy(p *SLAPolicy) error
	UpdatePolicy(p *SLAPolicy) error
	DeletePolicy(id string) error

	// Operational Tasks
	FindTasks(status, taskType, assignedTo string, limit, offset int) ([]OperationalTask, int64, error)
	FindTaskById(id string) (*OperationalTask, error)
	CreateTask(t *OperationalTask) error
	UpdateTask(t *OperationalTask) error
	DeleteTask(id string) error

	// Performance Snapshot
	FindPerformanceSnapshots(startDate, endDate string) ([]StaffSLAPerformanceDaily, error)
	SavePerformanceSnapshot(snap *StaffSLAPerformanceDaily) error
}

type slaRepository struct {
	db *gorm.DB
}

func NewSLARepository(db *gorm.DB) SLARepository {
	return &slaRepository{db: db}
}

func (r *slaRepository) FindPolicies() ([]SLAPolicy, error) {
	var list []SLAPolicy
	err := r.db.Order("category ASC, priority ASC").Find(&list).Error
	return list, err
}

func (r *slaRepository) FindPolicyById(id string) (*SLAPolicy, error) {
	var p SLAPolicy
	err := r.db.Where("id = ?", id).First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *slaRepository) CreatePolicy(p *SLAPolicy) error {
	return r.db.Create(p).Error
}

func (r *slaRepository) UpdatePolicy(p *SLAPolicy) error {
	return r.db.Save(p).Error
}

func (r *slaRepository) DeletePolicy(id string) error {
	return r.db.Where("id = ?", id).Delete(&SLAPolicy{}).Error
}

func (r *slaRepository) FindTasks(status, taskType, assignedTo string, limit, offset int) ([]OperationalTask, int64, error) {
	var tasks []OperationalTask
	var total int64

	q := r.db.Model(&OperationalTask{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if taskType != "" {
		q = q.Where("task_type = ?", taskType)
	}
	if assignedTo != "" {
		q = q.Where("assigned_to_user_id = ?", assignedTo)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit <= 0 {
		limit = 20
	}

	err := q.Preload("SLAPolicy").
		Preload("RoomUnit").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&tasks).Error

	return tasks, total, err
}

func (r *slaRepository) FindTaskById(id string) (*OperationalTask, error) {
	var task OperationalTask
	err := r.db.Preload("SLAPolicy").Preload("RoomUnit").Where("id = ?", id).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (r *slaRepository) CreateTask(t *OperationalTask) error {
	return r.db.Create(t).Error
}

func (r *slaRepository) UpdateTask(t *OperationalTask) error {
	return r.db.Save(t).Error
}

func (r *slaRepository) DeleteTask(id string) error {
	return r.db.Where("id = ?", id).Delete(&OperationalTask{}).Error
}

func (r *slaRepository) FindPerformanceSnapshots(startDate, endDate string) ([]StaffSLAPerformanceDaily, error) {
	var list []StaffSLAPerformanceDaily
	q := r.db.Model(&StaffSLAPerformanceDaily{})
	if startDate != "" && endDate != "" {
		q = q.Where("snapshot_date >= ? AND snapshot_date <= ?", startDate, endDate)
	}
	err := q.Order("snapshot_date DESC, department ASC").Find(&list).Error
	return list, err
}

func (r *slaRepository) SavePerformanceSnapshot(snap *StaffSLAPerformanceDaily) error {
	return r.db.Save(snap).Error
}
