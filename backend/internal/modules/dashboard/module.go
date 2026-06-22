package dashboard

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardModule struct {
	Handler *DashboardHandler
}

func InitModule(db *gorm.DB) *DashboardModule {
	// Inisialisasi dependency secara berurutan (Repo -> Service -> Handler)
	repo := NewDashboardRepository(db)
	svc := NewDashboardService(repo)
	hdl := NewDashboardHandler(svc)

	return &DashboardModule{
		Handler: hdl,
	}
}

// DashboardRoutes untuk mendaftarkan endpoint khusus dashboard
func (m *DashboardModule) DashboardRoutes(router *gin.RouterGroup) {
	// Hilangkan slash di akhir grup agar bersih
	api := router.Group("/dashboard")
	{

		api.GET("/summary", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.GetDashboardData)
	}
}
