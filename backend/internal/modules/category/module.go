package category

import (
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
	category.POST("/", m.Handler.Create)
	category.GET("/:slug", m.Handler.GetBySlug)
	category.GET("/detail/:id", m.Handler.GetById)
	category.PATCH("/:id", m.Handler.Update)
	category.DELETE("/:id", m.Handler.Delete)
}
