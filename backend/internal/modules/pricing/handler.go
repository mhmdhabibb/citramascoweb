package pricing

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type PricingHandler struct {
	service PricingService
}

func NewPricingHandler(service PricingService) *PricingHandler {
	return &PricingHandler{service: service}
}

// Seasons
func (h *PricingHandler) GetSeasons(c *gin.Context) {
	data, err := h.service.GetSeasons()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *PricingHandler) CreateSeason(c *gin.Context) {
	var req CreateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.CreateSeason(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Season created", "data": data})
}

func (h *PricingHandler) UpdateSeason(c *gin.Context) {
	id := c.Param("id")
	var req UpdateSeasonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.UpdateSeason(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Season updated", "data": data})
}

func (h *PricingHandler) DeleteSeason(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteSeason(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Season deleted"})
}

// Rates
func (h *PricingHandler) GetRates(c *gin.Context) {
	roomTypeId := c.Query("room_type_id")
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		start = time.Now()
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		end = start.AddDate(0, 1, 0)
	}

	data, err := h.service.GetRates(roomTypeId, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *PricingHandler) SetRate(c *gin.Context) {
	userId := c.GetString("user_id")
	var req SetDailyRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.SetDailyRate(userId, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Daily rate set", "data": data})
}

// Market Snapshots & Competitor Rates
func (h *PricingHandler) GetSnapshots(c *gin.Context) {
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")

	start, err := time.Parse("2006-01-02", startStr)
	if err != nil {
		start = time.Now().AddDate(0, -1, 0)
	}
	end, err := time.Parse("2006-01-02", endStr)
	if err != nil {
		end = time.Now()
	}

	data, err := h.service.GetSnapshots(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *PricingHandler) GetCompetitorRates(c *gin.Context) {
	dateStr := c.Query("date")
	var checkDate time.Time
	if dateStr != "" {
		checkDate, _ = time.Parse("2006-01-02", dateStr)
	}
	data, err := h.service.GetCompetitorRates(checkDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

// Events
func (h *PricingHandler) GetMarketEvents(c *gin.Context) {
	startStr := c.Query("start_date")
	endStr := c.Query("end_date")
	var start, end time.Time
	if startStr != "" {
		start, _ = time.Parse("2006-01-02", startStr)
	}
	if endStr != "" {
		end, _ = time.Parse("2006-01-02", endStr)
	}
	data, err := h.service.GetMarketEvents(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}

func (h *PricingHandler) CreateMarketEvent(c *gin.Context) {
	var req CreateMarketEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	data, err := h.service.CreateMarketEvent(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "Market event created", "data": data})
}

func (h *PricingHandler) DeleteMarketEvent(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteMarketEvent(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Market event deleted"})
}

// KPI
func (h *PricingHandler) GetKPI(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", strconv.Itoa(time.Now().Year())))
	data, err := h.service.GetKPISummaries(year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": data})
}
