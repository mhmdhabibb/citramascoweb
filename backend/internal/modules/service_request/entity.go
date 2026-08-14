package service_request

import (
	"citramascoweb-backend/internal/modules/auth"
	"citramascoweb-backend/internal/modules/rooms"
	"time"
)

type ServiceRequestCategory string

const (
	CategoryBrokenItem      ServiceRequestCategory = "incident_broken_item" // Gelas pecah, piring pecah
	CategoryExtraCleaning   ServiceRequestCategory = "extra_cleaning"       // Tumpahan air, ganti sprei
	CategoryAmenities       ServiceRequestCategory = "amenities_request"    // Tambahan handuk, sabun
	CategoryMaintenance     ServiceRequestCategory = "maintenance_repair"   // AC bocor, lampu mati
	CategoryOther           ServiceRequestCategory = "other"
)

type ServiceRequestStatus string

const (
	StatusPendingReception      ServiceRequestStatus = "pending_reception"        // Baru dilaporkan tamu, menunggu frontdesk
	StatusAssignedToHousekeeper ServiceRequestStatus = "assigned_to_housekeeping" // Sudah ditugaskan frontdesk ke housekeeping
	StatusInProgress            ServiceRequestStatus = "in_progress"              // Sedang dikerjakan
	StatusCompleted             ServiceRequestStatus = "completed"                // Selesai
	StatusCancelled             ServiceRequestStatus = "cancelled"                // Dibatalkan
)

type ServiceRequestPriority string

const (
	PriorityLow    ServiceRequestPriority = "low"
	PriorityMedium ServiceRequestPriority = "medium"
	PriorityHigh   ServiceRequestPriority = "high"
	PriorityUrgent ServiceRequestPriority = "urgent"
)

type ServiceRequest struct {
	Id                 string                 `gorm:"type:varchar(191);primaryKey" json:"id"`
	RoomId             string                 `gorm:"type:varchar(191);index" json:"room_id"`
	Room               rooms.Room             `gorm:"foreignKey:RoomId;references:Id" json:"room"`
	ReservationId      *string                `gorm:"type:varchar(191)" json:"reservation_id,omitempty"`
	GuestName          string                 `gorm:"type:varchar(255)" json:"guest_name"`
	GuestPhone         string                 `gorm:"type:varchar(50)" json:"guest_phone,omitempty"`
	Category           ServiceRequestCategory `gorm:"type:varchar(50)" json:"category"`
	Title              string                 `gorm:"type:varchar(255)" json:"title"`
	Description        string                 `gorm:"type:text" json:"description"`
	Priority           ServiceRequestPriority `gorm:"type:varchar(20);default:'medium'" json:"priority"`
	Status             ServiceRequestStatus   `gorm:"type:varchar(50);default:'pending_reception'" json:"status"`
	AssignedToUserId   *string                `gorm:"type:varchar(191)" json:"assigned_to_user_id,omitempty"`
	AssignedToUser     *auth.User             `gorm:"foreignKey:AssignedToUserId;references:Id" json:"assigned_to_user,omitempty"`
	NotesFromReception *string                `gorm:"type:text" json:"notes_from_reception,omitempty"`
	DamageCharge       int                    `gorm:"default:0" json:"damage_charge"` // Biaya ganti rugi barang jika ada
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	CompletedAt        *time.Time             `json:"completed_at,omitempty"`
}
