package smartlock

import (
	"citramascoweb-backend/internal/middlewares"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SmartLockModule struct {
	handler *SmartLockHandler
}

func NewSmartLockModule(db *gorm.DB) *SmartLockModule {
	repo := NewSmartLockRepository(db)
	service := NewSmartLockService(repo)
	handler := NewSmartLockHandler(service)
	return &SmartLockModule{handler: handler}
}

func (m *SmartLockModule) RegisterRoutes(r *gin.RouterGroup) {
	lock := r.Group("/smartlock")
	lock.Use(middlewares.AuthMiddleware())

	lock.GET("/units/:unitId/credentials", m.handler.GetUnitCredentials)
	lock.GET("/reservations/:reservationId/credentials", m.handler.GetReservationCredentials)
	lock.POST("/credentials", m.handler.GenerateCredential)
	lock.PATCH("/credentials/:id/revoke", m.handler.RevokeCredential)

	lock.GET("/logs", m.handler.GetAccessLogs)
	lock.POST("/logs", m.handler.RecordAccessLog) // Can also be used by IoT gateway webhook
}
