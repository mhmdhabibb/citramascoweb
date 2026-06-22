package reservation

import (
	"citramascoweb-backend/internal/dto"
	"log"
	"net/http"

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
	log.Printf("[DEBUG][RESERVATION-HANDLER] Processing availability check request")

	var req dto.CheckAvailabilityRequest
	// Menggunakan ShouldBindQuery untuk membaca query parameters (?check_in_date=...&check_out_date=...)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("[WARN][RESERVATION-HANDLER] Query binding failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid query parameters. Please provide valid check-in and check-out dates.",
		})
		return
	}

	log.Printf("[DEBUG][RESERVATION-HANDLER] Checking dates range: %s to %s (Room ID: '%s')", req.CheckInDate, req.CheckOutDate, req.RoomId)

	// Memanggil fungsi service yang sudah disinkronkan dengan status maintenance
	result, err := h.reservationService.CheckAvailability(&req)
	if err != nil {
		log.Printf("[ERROR][RESERVATION-HANDLER] Service failed to check room availability: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(), // Mengembalikan pesan error validasi dalam bahasa Inggris dari service
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Room availability has been checked successfully.",
		"data":    result,
	})
}

func (h *reservationHandler) CheckIn(c *gin.Context) {
	id := c.Param("id")
	if err := h.reservationService.CheckIn(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Proses check-in berhasil, status kamar diubah menjadi Occupied"})
}

func (h *reservationHandler) CheckOut(c *gin.Context) {
	id := c.Param("id")
	if err := h.reservationService.CheckOut(id); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Proses check-out berhasil, status kamar kembali Available"})
}
