package audit

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StaffAuditLog struct {
	Id          string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	UserId      string    `gorm:"type:varchar(191);index;not null" json:"user_id"`
	Action      string    `gorm:"type:varchar(50);not null" json:"action"` // create, update, delete, status_change, override
	TargetTable string    `gorm:"type:varchar(100);index;not null" json:"target_table"`
	TargetId    string    `gorm:"type:varchar(191);index" json:"target_id"`
	Payload     string    `gorm:"type:text" json:"payload"`
	CreatedAt   time.Time `json:"created_at"`
}

func (StaffAuditLog) TableName() string {
	return "staff_audit_logs"
}

func Record(db *gorm.DB, userId, action, table, targetId, payload string) {
	if db == nil {
		return
	}
	log := StaffAuditLog{
		Id:          uuid.New().String(),
		UserId:      userId,
		Action:      action,
		TargetTable: table,
		TargetId:    targetId,
		Payload:     payload,
		CreatedAt:   time.Now(),
	}
	db.Create(&log)
}

type AuditModule struct {
	db *gorm.DB
}

func NewAuditModule(db *gorm.DB) *AuditModule {
	return &AuditModule{db: db}
}

func (m *AuditModule) RegisterRoutes(r *gin.RouterGroup) {
	auditGroup := r.Group("/audit-logs")
	auditGroup.GET("", func(c *gin.Context) {
		var logs []StaffAuditLog
		q := m.db.Order("created_at DESC").Limit(100)
		if table := c.Query("table"); table != "" {
			q = q.Where("target_table = ?", table)
		}
		if userId := c.Query("user_id"); userId != "" {
			q = q.Where("user_id = ?", userId)
		}
		if err := q.Find(&logs).Error; err != nil {
			c.JSON(500, gin.H{"success": false, "message": err.Error()})
			return
		}
		c.JSON(200, gin.H{"success": true, "data": logs})
	})
}
