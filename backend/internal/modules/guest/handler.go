package guest

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GuestHandler struct {
	service GuestService
}

func NewGuestHandler(service GuestService) *GuestHandler {
	return &GuestHandler{service: service}
}

// Identity Types Handlers
func (h *GuestHandler) GetIdentityTypes(c *gin.Context) {
	data, err := h.service.GetAllIdentityTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *GuestHandler) CreateIdentityType(c *gin.Context) {
	var req CreateIdentityTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.CreateIdentityType(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Identity type created", "data": data})
}

func (h *GuestHandler) UpdateIdentityType(c *gin.Context) {
	id := c.Param("id")
	var req UpdateIdentityTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.UpdateIdentityType(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Identity type updated", "data": data})
}

func (h *GuestHandler) DeleteIdentityType(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteIdentityType(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Identity type deleted"})
}

// Guests Handlers
func (h *GuestHandler) GetGuests(c *gin.Context) {
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	data, total, err := h.service.GetAllGuests(search, page, limit)
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

func (h *GuestHandler) GetGuestById(c *gin.Context) {
	id := c.Param("id")
	data, err := h.service.GetGuestById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *GuestHandler) CreateGuest(c *gin.Context) {
	var req CreateGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.CreateGuest(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Guest created successfully", "data": data})
}

func (h *GuestHandler) UpdateGuest(c *gin.Context) {
	id := c.Param("id")
	var req UpdateGuestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.UpdateGuest(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Guest updated successfully", "data": data})
}

func (h *GuestHandler) DeleteGuest(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteGuest(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Guest deleted successfully"})
}
