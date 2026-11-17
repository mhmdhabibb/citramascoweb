package facility

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type FacilityHandler struct {
	service FacilityService
}

func NewFacilityHandler(service FacilityService) *FacilityHandler {
	return &FacilityHandler{service: service}
}

func (h *FacilityHandler) GetAllFacilities(c *gin.Context) {
	data, err := h.service.GetAllFacilities()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *FacilityHandler) GetFacilityById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetFacilityById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *FacilityHandler) CreateFacility(c *gin.Context) {
	var req CreateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.CreateFacility(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Facility created", "data": data})
}

func (h *FacilityHandler) UpdateFacility(c *gin.Context) {
	id := c.Param("id")
	var req UpdateFacilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.UpdateFacility(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Facility updated", "data": data})
}

func (h *FacilityHandler) DeleteFacility(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteFacility(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Facility deleted"})
}

// Room Type Facilities
func (h *FacilityHandler) GetRoomFacilities(c *gin.Context) {
	roomTypeId := c.Param("roomTypeId")
	data, err := h.service.GetFacilitiesByRoomType(roomTypeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *FacilityHandler) SyncRoomFacilities(c *gin.Context) {
	roomTypeId := c.Param("roomTypeId")
	var req AssignFacilitiesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	if err := h.service.SyncRoomTypeFacilities(roomTypeId, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Facilities synced for room type"})
}

// Room Type Images
func (h *FacilityHandler) GetRoomImages(c *gin.Context) {
	roomTypeId := c.Param("roomTypeId")
	data, err := h.service.GetImagesByRoomType(roomTypeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *FacilityHandler) AddRoomImage(c *gin.Context) {
	roomTypeId := c.Param("roomTypeId")
	var req CreateRoomTypeImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.CreateRoomTypeImage(roomTypeId, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Image added to room type", "data": data})
}

func (h *FacilityHandler) DeleteRoomImage(c *gin.Context) {
	imageId := c.Param("imageId")
	if err := h.service.DeleteRoomTypeImage(imageId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Room image deleted"})
}

func (h *FacilityHandler) SetPrimaryRoomImage(c *gin.Context) {
	roomTypeId := c.Param("roomTypeId")
	imageId := c.Param("imageId")
	if err := h.service.SetPrimaryImage(roomTypeId, imageId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Primary image updated"})
}
