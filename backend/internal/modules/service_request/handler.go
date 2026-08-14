package service_request

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type serviceRequestHandler struct {
	service ServiceRequestServiceInterface
}

func NewServiceRequestHandler(service ServiceRequestServiceInterface) *serviceRequestHandler {
	return &serviceRequestHandler{service: service}
}

func (h *serviceRequestHandler) Create(c *gin.Context) {
	var dto CreateServiceRequestDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	result, err := h.service.CreateRequest(&dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Laporan layanan kamar berhasil dikirim ke Resepsionis.",
		"data":    result,
	})
}

func (h *serviceRequestHandler) GetAll(c *gin.Context) {
	status := c.Query("status")
	items, err := h.service.GetAll(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

func (h *serviceRequestHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	item, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": item})
}

func (h *serviceRequestHandler) AssignToHousekeeping(c *gin.Context) {
	id := c.Param("id")
	var dto AssignTaskDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.AssignToHousekeeping(id, &dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tugas berhasil diteruskan ke tim Housekeeping.",
	})
}

func (h *serviceRequestHandler) CompleteRequest(c *gin.Context) {
	id := c.Param("id")
	var dto CompleteTaskDto
	_ = c.ShouldBindJSON(&dto)

	if err := h.service.CompleteRequest(id, &dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Tugas berhasil diselesaikan oleh Housekeeping.",
	})
}

func (h *serviceRequestHandler) CancelRequest(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.CancelRequest(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Laporan layanan kamar berhasil dibatalkan.",
	})
}
