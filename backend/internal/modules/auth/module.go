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
	repo := NewAuthRepository(db)
	service := NewAuthService(repo)
	handler := NewAuthHandler(service)

	return &Module{
		Handler: handler,
	}
}

func (m *Module) AuthRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	auth.POST("/sign-up", m.Handler.Register)
	auth.POST("/sign-in", m.Handler.Login)
	auth.GET("/me", middlewares.AuthMiddleware(), m.Handler.Profile)
	auth.PATCH("/update-profile", middlewares.AuthMiddleware(), m.Handler.UpdateProfile)
	auth.POST("/sign-out", middlewares.AuthMiddleware(), m.Handler.Logout)
}
