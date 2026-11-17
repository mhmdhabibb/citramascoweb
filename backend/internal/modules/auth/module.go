package auth

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler *authHandler
}

func InitModule(db *gorm.DB) *Module {
	service := NewAuthService(db)
	handler := NewAuthHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) AuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", m.Handler.SignIn)
		auth.POST("/login", m.Handler.SignIn)
		auth.GET("/me", middlewares.AuthMiddleware(), m.Handler.Profile)
	}
}

