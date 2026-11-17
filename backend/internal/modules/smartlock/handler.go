package smartlock

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SmartLockHandler struct {
	service SmartLockService
}

func NewSmartLockHandler(service SmartLockService) *SmartLockHandler {
	return &SmartLockHandler{service: service}
}

func (h *SmartLockHandler) GetUnitCredentials(c *gin.Context) {
	unitId := c.Param("unitId")
	data, err := h.service.GetCredentialsByUnit(unitId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *SmartLockHandler) GetReservationCredentials(c *gin.Context) {
	resId := c.Param("reservationId")
	data, err := h.service.GetCredentialsByReservation(resId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *SmartLockHandler) GenerateCredential(c *gin.Context) {
	var req CreateCredentialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.GenerateCredential(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Access credential created", "data": data})
}

func (h *SmartLockHandler) RevokeCredential(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.RevokeCredential(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Access credential revoked"})
}

func (h *SmartLockHandler) GetAccessLogs(c *gin.Context) {
	unitId := c.Query("unit_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	data, total, err := h.service.GetAccessLogs(unitId, page, limit)
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

func (h *SmartLockHandler) RecordAccessLog(c *gin.Context) {
	var req RecordAccessLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.RecordLog(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Access log recorded", "data": data})
}
