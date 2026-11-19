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
	// Room Types
	router.GET("/types", m.Handler.GetAll)

	category := router.Group("/type")
	category.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Store)
	category.GET("/:id", m.Handler.GetById)
	category.PATCH("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Update)
	category.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Delete)

	// Bed Types (Master Data: Single, Twin, Double, Queen, King, Super King)
	router.GET("/bed-types", m.Handler.GetAllBedTypes)
	bed := router.Group("/bed-type")
	bed.GET("/:id", m.Handler.GetBedTypeById)
	bed.POST("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.CreateBedType)
	bed.PATCH("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.UpdateBedType)
	bed.PUT("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.UpdateBedType)
	bed.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.DeleteBedType)
}
