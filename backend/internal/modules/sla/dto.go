package sla

type CreateSLAPolicyRequest struct {
	Name                    string `json:"name" binding:"required"`
	Category                string `json:"category" binding:"required"`
	Priority                string `json:"priority" binding:"required"`
	TargetResponseMinutes   int    `json:"target_response_minutes" binding:"required"`
	TargetResolutionMinutes int    `json:"target_resolution_minutes" binding:"required"`
	IsActive                bool   `json:"is_active"`
}

type UpdateSLAPolicyRequest struct {
	Name                    string `json:"name"`
	Category                string `json:"category"`
	Priority                string `json:"priority"`
	TargetResponseMinutes   int    `json:"target_response_minutes"`
	TargetResolutionMinutes int    `json:"target_resolution_minutes"`
	IsActive                bool   `json:"is_active"`
}

type CreateOperationalTaskRequest struct {
	SLAPolicyId      *string `json:"sla_policy_id"`
	RoomUnitId       *string `json:"room_unit_id"`
	ReservationId    *string `json:"reservation_id"`
	Title            string  `json:"title" binding:"required"`
	Description      string  `json:"description"`
	TaskType         string  `json:"task_type" binding:"required"`
	Priority         string  `json:"priority"`
	AssignedToUserId *string `json:"assigned_to_user_id"`
}

type UpdateTaskStatusRequest struct {
	Status       string  `json:"status" binding:"required"` // in_progress, resolved, inspected, cancelled
	BreachReason *string `json:"breach_reason"`
}

type AssignTaskRequest struct {
	AssignedToUserId string `json:"assigned_to_user_id" binding:"required"`
}
