package notification

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *notificationHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewNotificationRepository(db)
	service := NewNotificationService(repo)
	handler := NewNotificationHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) NotificationRoutes(router *gin.RouterGroup) {
	g := router.Group("/notifications", middlewares.AuthMiddleware())
	g.GET("/", m.Handler.List)
	g.GET("/unread-count", m.Handler.UnreadCount)
	g.PATCH("/:id/read", m.Handler.MarkRead)
}