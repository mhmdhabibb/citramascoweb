package types

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *typeHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewTypeRepository(db)
	service := NewTypeService(repo)
	handler := NewTypeHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) TypeRoutes(router *gin.RouterGroup) {
	router.GET("/types", m.Handler.GetAll)

	category := router.Group("/type")
	category.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Store)
	category.GET("/:id", middlewares.AuthMiddleware(), m.Handler.GetById)
	category.PATCH("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Update)
	category.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Delete)
}
