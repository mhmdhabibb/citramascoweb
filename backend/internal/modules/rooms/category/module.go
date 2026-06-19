package category

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *categoryHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewCategoryRepository(db)
	service := NewCategoryService(repo)
	handler := NewCategoryHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) CategoryRoutes(router *gin.RouterGroup) {
	router.GET("/categories", m.Handler.GetAll)

	category := router.Group("/category")
	category.POST("/", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Create)
	category.GET("/:slug", m.Handler.GetBySlug)
	category.GET("/detail/:id", m.Handler.GetById)
	category.PATCH("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Update)
	category.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Delete)
}
