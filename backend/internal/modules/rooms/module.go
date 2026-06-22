package rooms

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *roomHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewRoomRepository(db)
	service := NewRoomService(repo)
	handler := NewRoomHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) RoomRoutes(router *gin.RouterGroup) {
	router.GET("/rooms", m.Handler.GetAll)
	router.GET("/room/:id", m.Handler.GetById)

	room := router.Group("/room")
	room.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Store)
	room.PATCH("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Update)
	room.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Delete)
	room.GET("/filter/category/:id", m.Handler.FilerByCategory)
	room.GET("/filter/type/:id", m.Handler.FilterByType)
	room.GET("/filter", m.Handler.Filter)
	room.PATCH("/status/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.UpdateStatus)
}
