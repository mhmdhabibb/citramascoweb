package reservation

import (
	"citramascoweb-backend/internal/middlewares"
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

	service := NewReservationService(repo, roomRepo)
	handler := NewReservationHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) ReservationRoutes(router *gin.RouterGroup) {
	router.GET("/reservations", m.Handler.GetAll)
	router.GET("/reservation/:id", m.Handler.GetById)

	res := router.Group("/reservation")
	res.POST("", middlewares.AuthMiddleware(), m.Handler.Store)
	res.PATCH("/:id", middlewares.AuthMiddleware(), m.Handler.Update)
	res.DELETE("/:id", middlewares.AuthMiddleware(), m.Handler.Delete)
	res.PATCH("/approve/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.ApproveReservation)
	res.PATCH("/cancel/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.CancelReservation)
	res.PATCH("/reject/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.RejectReservation)
}
