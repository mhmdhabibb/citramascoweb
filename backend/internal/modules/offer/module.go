package offer

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *offerHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewOfferRepository(db)
	service := NewOfferService(repo)
	handler := NewOfferHandler(service)

	return &Module{Handler: handler}
}

func (m *Module) OfferRoutes(rg *gin.RouterGroup) {
	rg.GET("/offers", m.Handler.GetAll)
	offers := rg.Group("/offer")
	offers.GET("", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.GetAlls)
	offers.GET("/:id", m.Handler.GetById)
	offers.POST("/store", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Create)
	offers.PATCH("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Update)
	offers.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Delete)
}
