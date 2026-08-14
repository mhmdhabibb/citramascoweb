package service_request

import (
	"time"

	"gorm.io/gorm"
)

type ServiceRequestRepositoryInterface interface {
	Create(req *ServiceRequest) error
	GetAll(status string) ([]ServiceRequest, error)
	GetById(id string) (*ServiceRequest, error)
	AssignTask(id string, assignedUserId string, notes string, priority string) error
	CompleteTask(id string, damageCharge int) error
	Cancel(id string) error
}

type serviceRequestRepository struct {
	db *gorm.DB
}

func NewServiceRequestRepository(db *gorm.DB) ServiceRequestRepositoryInterface {
	_ = db.AutoMigrate(&ServiceRequest{})
	return &serviceRequestRepository{db: db}
}

func (r *serviceRequestRepository) Create(req *ServiceRequest) error {
	return r.db.Create(req).Error
}

func (r *serviceRequestRepository) GetAll(status string) ([]ServiceRequest, error) {
	var items []ServiceRequest
	query := r.db.Preload("Room").Preload("AssignedToUser").Order("created_at DESC")
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&items).Error
	return items, err
}

func (r *serviceRequestRepository) GetById(id string) (*ServiceRequest, error) {
	var item ServiceRequest
	err := r.db.Preload("Room").Preload("AssignedToUser").Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *serviceRequestRepository) AssignTask(id string, assignedUserId string, notes string, priority string) error {
	updates := map[string]interface{}{
		"status":                StatusAssignedToHousekeeper,
		"notes_from_reception":  notes,
		"updated_at":            time.Now(),
	}
	if assignedUserId != "" {
		updates["assigned_to_user_id"] = assignedUserId
	}
	if priority != "" {
		updates["priority"] = priority
	}
	return r.db.Model(&ServiceRequest{}).Where("id = ?", id).Updates(updates).Error
}

func (r *serviceRequestRepository) CompleteTask(id string, damageCharge int) error {
	now := time.Now()
	return r.db.Model(&ServiceRequest{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        StatusCompleted,
		"damage_charge": damageCharge,
		"completed_at":  &now,
		"updated_at":    now,
	}).Error
}

func (r *serviceRequestRepository) Cancel(id string) error {
	return r.db.Model(&ServiceRequest{}).Where("id = ?", id).Update("status", StatusCancelled).Error
}
