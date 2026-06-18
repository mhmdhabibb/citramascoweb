package reservation

import (
	"citramascoweb-backend/internal/dto"

	"github.com/gin-gonic/gin"
)

type reservationHandler struct {
	reservationService *reservationService
}

func NewReservationHandler(reservationService *reservationService) *reservationHandler {
	return &reservationHandler{reservationService: reservationService}
}

func (h *reservationHandler) GetAll(c *gin.Context) {
	reservations, err := h.reservationService.GetAll()
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Reservations retrieved successfully", "data": reservations})
}

func (h *reservationHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	reservation, err := h.reservationService.GetById(id)
	if err != nil {
		c.JSON(404, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Detail of reservation ", "data": reservation})
}

func (h *reservationHandler) Store(c *gin.Context) {

	var req dto.CreateReservationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	err = h.reservationService.Store(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Reservation created successfully!"})
}

func (h *reservationHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateReservationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	err = h.reservationService.Update(id, &req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Reservation updated successfully!"})
}

func (h *reservationHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.reservationService.Delete(id)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Reservation deleted successfully!"})
}

func (h *reservationHandler) ApproveReservation(c *gin.Context) {
	id := c.Param("id")

	err := h.reservationService.ApproveReservation(id)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Reservation approved successfully!"})
}

func (h *reservationHandler) CancelReservation(c *gin.Context) {
	id := c.Param("id")
	err := h.reservationService.CancelReservation(id)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Reservation cancelled successfully!"})
}

func (h *reservationHandler) RejectReservation(c *gin.Context) {
	id := c.Param("id")
	err := h.reservationService.RejectReservation(id)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Reservation rejected successfully!"})
}

func (h *reservationHandler) CheckAvailability(c *gin.Context) {
	var req dto.CheckAvailabilityRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	result, err := h.reservationService.CheckAvailability(&req)
	if err != nil {
		c.JSON(400, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": true, "message": "Availability checked successfully", "data": result})
}
