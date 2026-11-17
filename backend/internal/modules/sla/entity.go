package sla

import (
	"citramascoweb-backend/internal/modules/rooms/unit"
	"time"

	"gorm.io/gorm"
)

type SLAPolicy struct {
	Id                      string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	Name                    string         `gorm:"type:varchar(100);not null" json:"name"`
	Category                string         `gorm:"type:varchar(50);not null" json:"category"` // housekeeping, maintenance, front_desk
	Priority                string         `gorm:"type:varchar(20);not null" json:"priority"` // urgent, high, medium, low
	TargetResponseMinutes   int            `gorm:"not null" json:"target_response_minutes"`
	TargetResolutionMinutes int            `gorm:"not null" json:"target_resolution_minutes"`
	IsActive                bool           `gorm:"default:true" json:"is_active"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
	DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SLAPolicy) TableName() string {
	return "sla_policies"
}

type OperationalTask struct {
	Id                      string         `gorm:"type:varchar(191);primaryKey" json:"id"`
	TaskNumber              string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"task_number"`
	SLAPolicyId             *string        `gorm:"type:varchar(191);index" json:"sla_policy_id"`
	SLAPolicy               *SLAPolicy     `gorm:"foreignKey:SLAPolicyId;references:Id" json:"sla_policy,omitempty"`
	RoomUnitId              *string        `gorm:"type:varchar(191);index" json:"room_unit_id"`
	RoomUnit                *unit.RoomUnit `gorm:"foreignKey:RoomUnitId;references:Id" json:"room_unit,omitempty"`
	ReservationId           *string        `gorm:"type:varchar(191);index" json:"reservation_id"`
	Title                   string         `gorm:"type:varchar(255);not null" json:"title"`
	Description             string         `gorm:"type:text" json:"description"`
	TaskType                string         `gorm:"type:varchar(50);not null" json:"task_type"` // cleaning, repair, inspection, amenities
	Priority                string         `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	ReportedByUserId        *string        `gorm:"type:varchar(191)" json:"reported_by_user_id"`
	AssignedToUserId        *string        `gorm:"type:varchar(191);index" json:"assigned_to_user_id"`
	Status                  string         `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, in_progress, resolved, inspected, cancelled
	AcknowledgedAt          *time.Time     `json:"acknowledged_at"`
	ResolvedAt              *time.Time     `json:"resolved_at"`
	InspectedAt             *time.Time     `json:"inspected_at"`
	ActualResponseMinutes   *int           `json:"actual_response_minutes"`
	ActualResolutionMinutes *int           `json:"actual_resolution_minutes"`
	IsResponseBreached      bool           `gorm:"default:false" json:"is_response_breached"`
	IsResolutionBreached    bool           `gorm:"default:false" json:"is_resolution_breached"`
	BreachReason            *string        `gorm:"type:text" json:"breach_reason"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`
}

func (OperationalTask) TableName() string {
	return "operational_tasks"
}

type StaffSLAPerformanceDaily struct {
	Id                       string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	SnapshotDate             time.Time `gorm:"type:date;index;not null" json:"snapshot_date"`
	StaffId                  string    `gorm:"type:varchar(191);index;not null" json:"staff_id"`
	Department               string    `gorm:"type:varchar(100)" json:"department"`
	TotalTasksAssigned       int       `gorm:"default:0" json:"total_tasks_assigned"`
	TotalTasksCompleted      int       `gorm:"default:0" json:"total_tasks_completed"`
	TotalBreachedTasks       int       `gorm:"default:0" json:"total_breached_tasks"`
	SLAComplianceRate        float64   `gorm:"type:decimal(5,2);default:100.00" json:"sla_compliance_rate"`
	AvgResolutionTimeMinutes float64   `gorm:"type:decimal(7,2);default:0" json:"avg_resolution_time_minutes"`
	CreatedAt                time.Time `json:"created_at"`
}

func (StaffSLAPerformanceDaily) TableName() string {
	return "staff_sla_performance_daily"
}
