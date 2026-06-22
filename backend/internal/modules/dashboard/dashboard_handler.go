package dashboard

import (
	"citramascoweb-backend/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service DashboardService
}

func NewDashboardHandler(service DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// Mengubah nama fungsi menjadi GetDashboardData agar selaras dengan Service Layer
func (h *DashboardHandler) GetDashboardData(c *gin.Context) {
	var param dto.DashboardParam

	// Bind query parameter (?date=2026-06-21)
	if err := c.ShouldBindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format query parameter salah"})
		return
	}

	// Panggil service yang sekarang sudah berjalan menggunakan Goroutine
	data, err := h.service.GetDashboardData(param)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Response sukses membawa data Card, Line Chart, dan Pie Chart sekaligus
	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data dashboard",
		"data":    data,
	})
}
