package unit

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoomUnitModule struct {
	Handler *roomUnitHandler
	Service RoomUnitServiceInterface
	Repo    RoomUnitRepositoryInterface
}

func InitModule(db *gorm.DB) *RoomUnitModule {
	repo := NewRoomUnitRepository(db)
	service := NewRoomUnitService(repo)
	handler := NewRoomUnitHandler(service)

	return &RoomUnitModule{
		Handler: handler,
		Service: service,
		Repo:    repo,
	}
}

func (m *RoomUnitModule) RoomUnitRoutes(router *gin.RouterGroup) {
	// Room Units
	router.GET("/room-units", m.Handler.GetAll)
	router.GET("/room-unit/:id", m.Handler.GetById)
	router.GET("/room-unit/next-available/:roomId", m.Handler.GetNextAvailable)
	router.POST("/room-unit", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.Create)
	router.PUT("/room-unit/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.Update)
	router.PATCH("/room-unit/status/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "housekeeping", "reception"), m.Handler.UpdateStatus)
	router.DELETE("/room-unit/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Delete)

	// Room Statuses (Master Data: VC, VD, OC, OD, OOO, OOS)
	router.GET("/room-statuses", m.Handler.GetAllStatuses)
	router.GET("/room-status", m.Handler.GetAllStatuses)
	router.GET("/room-status/:id", m.Handler.GetStatusById)
	router.POST("/room-status", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.CreateStatus)
	router.PUT("/room-status/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.UpdateStatusModel)
	router.DELETE("/room-status/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.DeleteStatus)
}
