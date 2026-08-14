package service_request

import (
	"citramascoweb-backend/internal/middlewares"
	"citramascoweb-backend/internal/modules/notification"
	"citramascoweb-backend/internal/modules/rooms"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *serviceRequestHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewServiceRequestRepository(db)
	roomRepo := rooms.NewRoomRepository(db)
	notifier := notification.NewNotificationService(notification.NewNotificationRepository(db))
	service := NewServiceRequestService(repo, roomRepo, notifier)
	handler := NewServiceRequestHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) ServiceRequestRoutes(router *gin.RouterGroup) {
	// Guest endpoint to submit requests (e.g., via in-room QR code / public web)
	router.POST("/service-requests", m.Handler.Create)

	g := router.Group("/service-requests", middlewares.AuthMiddleware())
	g.GET("", m.Handler.GetAll)
	g.GET("/:id", m.Handler.GetById)
	g.PATCH("/:id/assign", middlewares.RoleMiddleware("admin", "manager", "reception"), m.Handler.AssignToHousekeeping)
	g.PATCH("/:id/complete", middlewares.RoleMiddleware("admin", "manager", "reception", "housekeeping"), m.Handler.CompleteRequest)
	g.PATCH("/:id/cancel", middlewares.RoleMiddleware("admin", "manager", "reception"), m.Handler.CancelRequest)
}
