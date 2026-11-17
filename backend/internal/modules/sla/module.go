package sla

import (
	"citramascoweb-backend/internal/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SLAModule struct {
	handler *SLAHandler
}

func NewSLAModule(db *gorm.DB) *SLAModule {
	repo := NewSLARepository(db)
	service := NewSLAService(repo)
	handler := NewSLAHandler(service)
	return &SLAModule{handler: handler}
}

func (m *SLAModule) RegisterRoutes(r *gin.RouterGroup) {
	slaGroup := r.Group("/sla")
	slaGroup.Use(middlewares.AuthMiddleware())

	// Policies
	slaGroup.GET("/policies", m.handler.GetPolicies)
	slaGroup.POST("/policies", middlewares.RoleMiddleware("admin", "manager"), m.handler.CreatePolicy)
	slaGroup.PUT("/policies/:id", middlewares.RoleMiddleware("admin", "manager"), m.handler.UpdatePolicy)
	slaGroup.DELETE("/policies/:id", middlewares.RoleMiddleware("admin", "manager"), m.handler.DeletePolicy)

	// Tasks
	slaGroup.GET("/tasks", m.handler.GetTasks)
	slaGroup.GET("/tasks/:id", m.handler.GetTaskById)
	slaGroup.POST("/tasks", m.handler.CreateTask)
	slaGroup.PATCH("/tasks/:id/assign", m.handler.AssignTask)
	slaGroup.PATCH("/tasks/:id/status", m.handler.UpdateTaskStatus)

	// Performance
	slaGroup.GET("/performance", m.handler.GetPerformance)
}

func SeedDefaultSLAPolicies(db *gorm.DB) {
	defaults := []struct {
		Name                    string
		Category                string
		Priority                string
		TargetResponseMinutes   int
		TargetResolutionMinutes int
	}{
		{Name: "Urgent Room Leak/Power Outage", Category: "maintenance", Priority: "urgent", TargetResponseMinutes: 5, TargetResolutionMinutes: 30},
		{Name: "AC / Electronics Repair", Category: "maintenance", Priority: "high", TargetResponseMinutes: 15, TargetResolutionMinutes: 60},
		{Name: "Extra Cleaning / Spill Incident", Category: "housekeeping", Priority: "high", TargetResponseMinutes: 10, TargetResolutionMinutes: 30},
		{Name: "Guest Amenities & Towels", Category: "housekeeping", Priority: "medium", TargetResponseMinutes: 15, TargetResolutionMinutes: 30},
		{Name: "Daily Room Turnover Cleaning", Category: "housekeeping", Priority: "medium", TargetResponseMinutes: 30, TargetResolutionMinutes: 60},
		{Name: "Minor Plumbing Issue", Category: "maintenance", Priority: "medium", TargetResponseMinutes: 30, TargetResolutionMinutes: 90},
	}

	for _, p := range defaults {
		var count int64
		db.Model(&SLAPolicy{}).Where("name = ?", p.Name).Count(&count)
		if count == 0 {
			db.Create(&SLAPolicy{
				Id:                      uuid.New().String(),
				Name:                    p.Name,
				Category:                p.Category,
				Priority:                p.Priority,
				TargetResponseMinutes:   p.TargetResponseMinutes,
				TargetResolutionMinutes: p.TargetResolutionMinutes,
				IsActive:                true,
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
			})
		}
	}
}
