package notification

import (
	"github.com/gin-gonic/gin"
)

type notificationHandler struct {
	service *NotificationService
}

func NewNotificationHandler(service *NotificationService) *notificationHandler {
	return &notificationHandler{service: service}
}

func (h *notificationHandler) List(c *gin.Context) {
	userId := c.GetString("user_id")
	items, err := h.service.GetMine(userId)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": items})
}

func (h *notificationHandler) UnreadCount(c *gin.Context) {
	userId := c.GetString("user_id")
	count, err := h.service.GetUnreadCount(userId)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": count})
}

func (h *notificationHandler) MarkRead(c *gin.Context) {
	userId := c.GetString("user_id")
	id := c.Param("id")
	if err := h.service.MarkRead(id, userId); err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Notification marked as read"})
}