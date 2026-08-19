package floor

import (
	"citramascoweb-backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type floorHandler struct {
	service FloorServiceInterface
}

func NewFloorHandler(service FloorServiceInterface) *floorHandler {
	return &floorHandler{service: service}
}

func (h *floorHandler) GetAll(c *gin.Context) {
	floors, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": floors})
}

func (h *floorHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	fl, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": fl})
}

func (h *floorHandler) Create(c *gin.Context) {
	var req dto.CreateFloorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Create(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Lantai berhasil dibuat"})
}

func (h *floorHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateFloorRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Update(id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Lantai berhasil diperbarui"})
}

func (h *floorHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Lantai berhasil dihapus"})
}
