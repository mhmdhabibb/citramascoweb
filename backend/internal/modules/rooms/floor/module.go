package floor

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FloorModule struct {
	Handler *floorHandler
	Service FloorServiceInterface
	Repo    FloorRepositoryInterface
}

func InitModule(db *gorm.DB) *FloorModule {
	repo := NewFloorRepository(db)
	service := NewFloorService(repo)
	handler := NewFloorHandler(service)

	return &FloorModule{
		Handler: handler,
		Service: service,
		Repo:    repo,
	}
}

func (m *FloorModule) FloorRoutes(router *gin.RouterGroup) {
	router.GET("/floors", m.Handler.GetAll)
	router.GET("/floor/:id", m.Handler.GetById)
	router.POST("/floor", m.Handler.Create)
	router.PUT("/floor/:id", m.Handler.Update)
	router.DELETE("/floor/:id", m.Handler.Delete)
}
