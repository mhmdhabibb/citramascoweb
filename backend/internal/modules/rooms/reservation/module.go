package reservation

import (
	"citramascoweb-backend/internal/middlewares"
	"citramascoweb-backend/internal/modules/notification"
	"citramascoweb-backend/internal/modules/rooms"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *reservationHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewReservationRepository(db)
	roomRepo := rooms.NewRoomRepository(db)
	notifier := notification.NewNotificationService(notification.NewNotificationRepository(db))

	service := NewReservationService(repo, roomRepo, notifier)
	handler := NewReservationHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) ReservationRoutes(router *gin.RouterGroup) {
	router.GET("/reservations", m.Handler.GetAll)
	router.GET("/reservations/availability", m.Handler.CheckAvailability)
	router.GET("/reservation/:id", m.Handler.GetById)

	res := router.Group("/reservation")
	res.POST("", m.Handler.Store)
	res.PATCH("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "reception"), m.Handler.Update)
	res.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.Delete)
	res.PATCH("/approve/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "reception"), m.Handler.ApproveReservation)
	res.PATCH("/cancel/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "reception"), m.Handler.CancelReservation)
	res.PATCH("/reject/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "reception"), m.Handler.RejectReservation)
	res.PATCH("/check-in/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "reception"), m.Handler.CheckIn)
	res.PATCH("/check-out/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "reception"), m.Handler.CheckOut)

}
