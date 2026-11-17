package guest

import (
	"citramascoweb-backend/internal/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GuestModule struct {
	handler *GuestHandler
}

func NewGuestModule(db *gorm.DB) *GuestModule {
	repo := NewGuestRepository(db)
	service := NewGuestService(repo)
	handler := NewGuestHandler(service)
	return &GuestModule{handler: handler}
}

func (m *GuestModule) RegisterRoutes(r *gin.RouterGroup) {
	// Identity Types
	idTypes := r.Group("/identity-types")
	idTypes.GET("", m.handler.GetIdentityTypes)
	idTypes.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.CreateIdentityType)
	idTypes.PUT("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.UpdateIdentityType)
	idTypes.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.handler.DeleteIdentityType)

	// Guests
	guests := r.Group("/guests")
	guests.Use(middlewares.AuthMiddleware())
	guests.GET("", m.handler.GetGuests)
	guests.GET("/:id", m.handler.GetGuestById)
	guests.POST("", m.handler.CreateGuest)
	guests.PUT("/:id", m.handler.UpdateGuest)
	guests.DELETE("/:id", middlewares.RoleMiddleware("admin", "manager"), m.handler.DeleteGuest)
}

func SeedDefaultIdentityTypes(db *gorm.DB) {
	defaults := []string{"KTP", "Paspor", "SIM"}
	for _, name := range defaults {
		var count int64
		db.Model(&IdentityType{}).Where("name = ?", name).Count(&count)
		if count == 0 {
			db.Create(&IdentityType{
				Id:        uuid.New().String(),
				Name:      name,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			})
		}
	}
}
