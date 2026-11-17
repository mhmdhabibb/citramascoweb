package facility

import (
	"citramascoweb-backend/internal/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FacilityModule struct {
	handler *FacilityHandler
}

func NewFacilityModule(db *gorm.DB) *FacilityModule {
	repo := NewFacilityRepository(db)
	service := NewFacilityService(repo)
	handler := NewFacilityHandler(service)
	return &FacilityModule{handler: handler}
}

func (m *FacilityModule) RegisterRoutes(r *gin.RouterGroup) {
	facilities := r.Group("/facilities")
	facilities.GET("", m.handler.GetAllFacilities)
	facilities.GET("/:id", m.handler.GetFacilityById)
	facilities.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.CreateFacility)
	facilities.PUT("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.UpdateFacility)
	facilities.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.DeleteFacility)

	// Room Type Relations
	roomTypes := r.Group("/room-types")
	roomTypes.GET("/:roomTypeId/facilities", m.handler.GetRoomFacilities)
	roomTypes.POST("/:roomTypeId/facilities", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.SyncRoomFacilities)

	roomTypes.GET("/:roomTypeId/images", m.handler.GetRoomImages)
	roomTypes.POST("/:roomTypeId/images", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.AddRoomImage)
	roomTypes.DELETE("/images/:imageId", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.DeleteRoomImage)
	roomTypes.PATCH("/:roomTypeId/images/:imageId/primary", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.SetPrimaryRoomImage)
}

func SeedDefaultFacilities(db *gorm.DB) {
	defaults := []struct {
		Name string
		Icon string
	}{
		{Name: "High-Speed WiFi", Icon: "wifi"},
		{Name: "Air Conditioning", Icon: "wind"},
		{Name: "Smart TV 43 inch", Icon: "tv"},
		{Name: "Water Heater", Icon: "droplets"},
		{Name: "Free Breakfast", Icon: "utensils"},
		{Name: "Bathtub", Icon: "bath"},
		{Name: "Safe Deposit Box", Icon: "shield"},
		{Name: "Coffee & Tea Maker", Icon: "coffee"},
		{Name: "Working Desk & Chair", Icon: "briefcase"},
		{Name: "Mini Refrigerator", Icon: "refrigerator"},
	}

	for _, item := range defaults {
		var count int64
		db.Model(&Facility{}).Where("name = ?", item.Name).Count(&count)
		if count == 0 {
			db.Create(&Facility{
				Id:        uuid.New().String(),
				Name:      item.Name,
				Icon:      item.Icon,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}
}
