package unit

import (
	"citramascoweb-backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type roomUnitHandler struct {
	service RoomUnitServiceInterface
}

func NewRoomUnitHandler(service RoomUnitServiceInterface) *roomUnitHandler {
	return &roomUnitHandler{service: service}
}

func (h *roomUnitHandler) GetAll(c *gin.Context) {
	roomId := c.Query("room_id")
	floorId := c.Query("floor_id")

	var units []RoomUnit
	var err error

	if roomId != "" {
		units, err = h.service.GetByRoomId(roomId)
	} else if floorId != "" {
		units, err = h.service.GetByFloorId(floorId)
	} else {
		units, err = h.service.GetAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": units})
}

func (h *roomUnitHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	unit, err := h.service.GetById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": unit})
}

func (h *roomUnitHandler) GetNextAvailable(c *gin.Context) {
	roomId := c.Param("roomId")
	unit, err := h.service.GetNextAvailableUnit(roomId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Tidak ada unit kamar yang tersedia untuk tipe ini"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": unit})
}

func (h *roomUnitHandler) Create(c *gin.Context) {
	var req dto.CreateRoomUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Create(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Room unit added successfully"})
}

func (h *roomUnitHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateRoomUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.Update(id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Room unit updated successfully"})
}

func (h *roomUnitHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateRoomUnitStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Room status changed successfully"})
}

func (h *roomUnitHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Room unit deleted successfully"})
}

// ---------------------------------------------------------------------------
// RoomStatus Handlers
// ---------------------------------------------------------------------------

func (h *roomUnitHandler) GetAllStatuses(c *gin.Context) {
	statuses, err := h.service.GetAllStatuses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": statuses})
}

func (h *roomUnitHandler) GetStatusById(c *gin.Context) {
	id := c.Param("id")
	status, err := h.service.GetStatusById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": status})
}

func (h *roomUnitHandler) CreateStatus(c *gin.Context) {
	var req dto.CreateRoomStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.CreateStatus(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Room status created successfully"})
}

func (h *roomUnitHandler) UpdateStatusModel(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateRoomStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if err := h.service.UpdateStatusModel(id, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Room status updated successfully"})
}

func (h *roomUnitHandler) DeleteStatus(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteStatus(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Room status deleted successfully"})
}
