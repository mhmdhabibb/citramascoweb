package inventory

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *inventoryHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewInventoryRepository(db)
	service := NewInventoryService(repo)
	handler := NewInventoryHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) InventoryRoutes(router *gin.RouterGroup) {
	inv := router.Group("/inventory")

	inv.GET("/items", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "inventory"), m.Handler.ListItems)
	inv.POST("/items", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.CreateItem)
	inv.PATCH("/items/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.UpdateItem)
	inv.DELETE("/items/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.DeleteItem)

	inv.POST("/items/:id/stock-in", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "inventory"), m.Handler.StockIn)
	inv.POST("/items/:id/usage", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "inventory"), m.Handler.Usage)
	inv.POST("/stock-take", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "inventory"), m.Handler.StockTake)

	inv.GET("/transactions", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "inventory"), m.Handler.Transactions)
	inv.GET("/report", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager", "inventory"), m.Handler.Report)
}