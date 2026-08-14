package channel

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *channelHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewChannelRepository(db)
	service := NewChannelService(repo)
	handler := NewChannelHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) ChannelRoutes(router *gin.RouterGroup) {
	router.GET("/channels", m.Handler.GetAll)
	channelRoutes := router.Group("/channel", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"))
	{
		channelRoutes.GET("", m.Handler.Get)
		channelRoutes.GET("/last", m.Handler.GetLastChannel)
		channelRoutes.GET("/:id", m.Handler.GetById)
		channelRoutes.GET("/by-name", m.Handler.GetByName)
		channelRoutes.POST("", m.Handler.Store)
		channelRoutes.PUT("/:id", m.Handler.Update)
		channelRoutes.DELETE("/:id", m.Handler.Delete)
	}
}
