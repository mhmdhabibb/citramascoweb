package monitoring

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MonitoringHandler struct {
	db        *gorm.DB
	collector *MetricsCollector
}

func NewMonitoringHandler(db *gorm.DB) *MonitoringHandler {
	return &MonitoringHandler{
		db:        db,
		collector: GetCollector(),
	}
}

// GetMetrics handles GET /metrics
// By default returns Prometheus plain text metrics.
// If "?format=json" is provided or "Accept: application/json" is requested, returns JSON report.
func (h *MonitoringHandler) GetMetrics(c *gin.Context) {
	format := strings.ToLower(c.Query("format"))
	accept := strings.ToLower(c.GetHeader("Accept"))

	wantsJSON := format == "json" || (strings.Contains(accept, "application/json") && !strings.Contains(accept, "text/plain"))

	if wantsJSON {
		report := h.collector.GetReport(h.db)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    report,
		})
		return
	}

	promText := h.collector.ToPrometheusFormat(h.db)
	c.Data(http.StatusOK, "text/plain; version=0.0.4; charset=utf-8", []byte(promText))
}

// GetServicesStatus handles GET /api/monitoring/services
// Returns complete categorized service breakdown, health status, and metrics in structured JSON
func (h *MonitoringHandler) GetServicesStatus(c *gin.Context) {
	report := h.collector.GetReport(h.db)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Service monitoring metrics retrieved successfully",
		"data":    report,
	})
}

// GetHealth handles GET /api/monitoring/health
// Quick health probe for load balancers or uptime monitors
func (h *MonitoringHandler) GetHealth(c *gin.Context) {
	report := h.collector.GetReport(h.db)
	statusCode := http.StatusOK
	if report.Database.Status != "healthy" {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":         report.Overall.Status,
		"database":       report.Database.Status,
		"uptime_seconds": report.Overall.UptimeSeconds,
		"uptime":         report.Overall.UptimeFormatted,
		"total_requests": report.Overall.TotalRequests,
		"avg_latency_ms": report.Overall.AvgLatencyMs,
		"error_rate":     report.Overall.GlobalErrorRate,
		"timestamp":      report.Timestamp,
	})
}
