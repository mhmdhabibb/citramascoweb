package unit

import (
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
	router.GET("/room-units", m.Handler.GetAll)
	router.GET("/room-unit/:id", m.Handler.GetById)
	router.GET("/room-unit/next-available/:roomId", m.Handler.GetNextAvailable)
	router.POST("/room-unit", m.Handler.Create)
	router.PUT("/room-unit/:id", m.Handler.Update)
	router.PATCH("/room-unit/status/:id", m.Handler.UpdateStatus)
	router.DELETE("/room-unit/:id", m.Handler.Delete)
}
