package user

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *userHandler
}

func InitModule(db *gorm.DB) *Module {
	repo := NewUserRepository(db)
	service := NewUserService(repo)
	handler := NewUserHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) UserRoutes(router *gin.RouterGroup) {
	router.GET("/customers", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin", "manager"), m.Handler.GetAllCustomer)

	users := router.Group("/users")

	users.GET("/", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.GetAllByRole)
	users.DELETE("/:id", middlewares.AuthMiddleware(), middlewares.RoleMiddleware("admin"), m.Handler.Delete)
}
