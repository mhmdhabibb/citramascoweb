package monitoring

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Module struct {
	Handler   *MonitoringHandler
	collector *MetricsCollector
}

func InitModule(db *gorm.DB) *Module {
	collector := GetCollector()
	handler := NewMonitoringHandler(db)
	return &Module{
		Handler:   handler,
		collector: collector,
	}
}

// Middleware returns a Gin middleware that records traffic and latencies for each request
func (m *Module) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		path := c.Request.URL.Path

		m.collector.Record(path, c.Request.Method, statusCode, duration)
	}
}

// RegisterRoutes registers monitoring and metrics endpoints on root and api groups
func (m *Module) RegisterRoutes(rootEngine *gin.Engine, apiGroup *gin.RouterGroup) {
	// Standard Prometheus / Metrics endpoint on root
	rootEngine.GET("/metrics", m.Handler.GetMetrics)

	// API grouped routes
	apiGroup.GET("/metrics", m.Handler.GetMetrics)
	apiGroup.GET("/monitoring/services", m.Handler.GetServicesStatus)
	apiGroup.GET("/monitoring/health", m.Handler.GetHealth)
}
