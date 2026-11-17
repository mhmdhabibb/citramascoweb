package sla

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SLAHandler struct {
	service SLAService
}

func NewSLAHandler(service SLAService) *SLAHandler {
	return &SLAHandler{service: service}
}

// Policies
func (h *SLAHandler) GetPolicies(c *gin.Context) {
	data, err := h.service.GetPolicies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *SLAHandler) CreatePolicy(c *gin.Context) {
	var req CreateSLAPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.CreatePolicy(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "SLA Policy created", "data": data})
}

func (h *SLAHandler) UpdatePolicy(c *gin.Context) {
	id := c.Param("id")
	var req UpdateSLAPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.UpdatePolicy(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "SLA Policy updated", "data": data})
}

func (h *SLAHandler) DeletePolicy(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePolicy(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "SLA Policy deleted"})
}

// Tasks
func (h *SLAHandler) GetTasks(c *gin.Context) {
	status := c.Query("status")
	taskType := c.Query("task_type")
	assignedTo := c.Query("assigned_to")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	data, total, err := h.service.GetTasks(status, taskType, assignedTo, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *SLAHandler) GetTaskById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetTaskById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *SLAHandler) CreateTask(c *gin.Context) {
	userId := c.GetString("user_id")
	var req CreateOperationalTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.CreateTask(userId, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Task created", "data": data})
}

func (h *SLAHandler) AssignTask(c *gin.Context) {
	id := c.Param("id")
	var req AssignTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.AssignTask(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Task assigned", "data": data})
}

func (h *SLAHandler) UpdateTaskStatus(c *gin.Context) {
	id := c.Param("id")
	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.UpdateTaskStatus(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Task status updated", "data": data})
}

func (h *SLAHandler) GetPerformance(c *gin.Context) {
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	data, err := h.service.GetDailyPerformance(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
