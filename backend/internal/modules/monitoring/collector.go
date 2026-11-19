package monitoring

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type CategoryStats struct {
	Name             string     `json:"name"`
	Category         string     `json:"category"`
	TotalRequests    int64      `json:"total_requests"`
	SuccessCount     int64      `json:"success_count"`
	ClientErrorCount int64      `json:"client_error_count"` // 4xx
	ServerErrorCount int64      `json:"server_error_count"` // 5xx
	TotalDurationMs  float64    `json:"-"`
	AvgLatencyMs     float64    `json:"avg_latency_ms"`
	MinLatencyMs     float64    `json:"min_latency_ms"`
	MaxLatencyMs     float64    `json:"max_latency_ms"`
	Status           string     `json:"status"` // healthy, degraded, down
	ErrorRate        float64    `json:"error_rate_percent"`
	LastStatusCode   int        `json:"last_status_code"`
	LastRequestAt    *time.Time `json:"last_request_at"`
}

type OverallStats struct {
	Status            string    `json:"status"` // healthy, degraded, down
	UptimeSeconds     int64     `json:"uptime_seconds"`
	UptimeFormatted   string    `json:"uptime_formatted"`
	StartTime         time.Time `json:"start_time"`
	TotalRequests     int64     `json:"total_requests"`
	Total2xx          int64     `json:"total_2xx"`
	Total3xx          int64     `json:"total_3xx"`
	Total4xx          int64     `json:"total_4xx"`
	Total5xx          int64     `json:"total_5xx"`
	GlobalErrorRate   float64   `json:"global_error_rate_percent"`
	AvgLatencyMs      float64   `json:"avg_latency_ms"`
	RequestsPerMinute float64   `json:"requests_per_minute"`
}

type SystemStats struct {
	GoVersion         string  `json:"go_version"`
	Goroutines        int     `json:"goroutines"`
	MemoryAllocatedMB float64 `json:"memory_allocated_mb"`
	TotalAllocatedMB  float64 `json:"total_allocated_mb"`
	SystemMemoryMB    float64 `json:"system_memory_mb"`
	NumGC             uint32  `json:"num_gc"`
}

type DatabaseStats struct {
	Status             string  `json:"status"` // healthy, down
	PingMs             float64 `json:"ping_ms"`
	OpenConnections    int     `json:"open_connections"`
	InUse              int     `json:"in_use"`
	Idle               int     `json:"idle"`
	WaitCount          int64   `json:"wait_count"`
	MaxOpenConnections int     `json:"max_open_connections"`
}

type MonitoringReport struct {
	Overall    OverallStats    `json:"overall"`
	Categories []CategoryStats `json:"categories"`
	Database   DatabaseStats   `json:"database"`
	System     SystemStats     `json:"system"`
	Timestamp  time.Time       `json:"timestamp"`
}

type MetricsCollector struct {
	mu         sync.RWMutex
	startTime  time.Time
	totalReq   int64
	total2xx   int64
	total3xx   int64
	total4xx   int64
	total5xx   int64
	totalDurMs float64
	categories map[string]*CategoryStats
}

var (
	defaultCollector *MetricsCollector
	onceCollector    sync.Once
)

var categoryDefinitions = []struct {
	Key         string
	DisplayName string
}{
	{Key: "auth", DisplayName: "Authentication & Identity Management"},
	{Key: "rooms", DisplayName: "Rooms, Physical Units & Master Setup"},
	{Key: "reservations", DisplayName: "Reservations & Front Desk Stays"},
	{Key: "guests", DisplayName: "Guest Book & Identities"},
	{Key: "inventory", DisplayName: "Inventory & Supply Tracking"},
	{Key: "finance", DisplayName: "Accounting, Journals & Financial Statements"},
	{Key: "channels", DisplayName: "OTA Channel Manager Integration"},
	{Key: "housekeeping", DisplayName: "Housekeeping & Service Requests"},
	{Key: "offers", DisplayName: "Promotions, Discounts & Vouchers"},
	{Key: "pricing", DisplayName: "Dynamic Pricing & Revenue Management"},
	{Key: "sla", DisplayName: "SLA & Operational Performance"},
	{Key: "smartlock", DisplayName: "Smart Lock & Room Keycard Access"},
	{Key: "notifications", DisplayName: "Live Notification Service"},
	{Key: "system", DisplayName: "Core System & Health Checks"},
}

func GetCollector() *MetricsCollector {
	onceCollector.Do(func() {
		catMap := make(map[string]*CategoryStats)
		for _, def := range categoryDefinitions {
			catMap[def.Key] = &CategoryStats{
				Category: def.Key,
				Name:     def.DisplayName,
				Status:   "healthy",
			}
		}
		defaultCollector = &MetricsCollector{
			startTime:  time.Now(),
			categories: catMap,
		}
	})
	return defaultCollector
}

func CategorizePath(path string) string {
	p := strings.ToLower(path)
	// Remove /api prefix for matching
	if strings.HasPrefix(p, "/api") {
		p = strings.TrimPrefix(p, "/api")
	}

	switch {
	case strings.HasPrefix(p, "/auth") || strings.HasPrefix(p, "/users") || strings.HasPrefix(p, "/user"):
		return "auth"
	case strings.HasPrefix(p, "/room") || strings.HasPrefix(p, "/floor") || strings.HasPrefix(p, "/type") || strings.HasPrefix(p, "/bed-type") || strings.HasPrefix(p, "/category") || strings.HasPrefix(p, "/categories"):
		return "rooms"
	case strings.HasPrefix(p, "/reservation"):
		return "reservations"
	case strings.HasPrefix(p, "/guest") || strings.HasPrefix(p, "/identit"):
		return "guests"
	case strings.HasPrefix(p, "/inventory"):
		return "inventory"
	case strings.HasPrefix(p, "/finance") || strings.HasPrefix(p, "/coa") || strings.HasPrefix(p, "/journal") || strings.HasPrefix(p, "/general-ledger") || strings.HasPrefix(p, "/cash-bank"):
		return "finance"
	case strings.HasPrefix(p, "/channel"):
		return "channels"
	case strings.HasPrefix(p, "/service-request"):
		return "housekeeping"
	case strings.HasPrefix(p, "/offer"):
		return "offers"
	case strings.HasPrefix(p, "/pricing") || strings.HasPrefix(p, "/season"):
		return "pricing"
	case strings.HasPrefix(p, "/sla"):
		return "sla"
	case strings.HasPrefix(p, "/smartlock"):
		return "smartlock"
	case strings.HasPrefix(p, "/notification"):
		return "notifications"
	case strings.HasPrefix(p, "/health") || strings.HasPrefix(p, "/metrics") || strings.HasPrefix(p, "/monitoring"):
		return "system"
	default:
		return "system"
	}
}

func (c *MetricsCollector) Record(path string, method string, statusCode int, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	durMs := float64(duration.Microseconds()) / 1000.0

	c.totalReq++
	c.totalDurMs += durMs

	switch {
	case statusCode >= 200 && statusCode < 300:
		c.total2xx++
	case statusCode >= 300 && statusCode < 400:
		c.total3xx++
	case statusCode >= 400 && statusCode < 500:
		c.total4xx++
	case statusCode >= 500:
		c.total5xx++
	}

	catKey := CategorizePath(path)
	cat, exists := c.categories[catKey]
	if !exists {
		cat = &CategoryStats{
			Category: catKey,
			Name:     strings.Title(catKey) + " Service",
			Status:   "healthy",
		}
		c.categories[catKey] = cat
	}

	cat.TotalRequests++
	cat.TotalDurationMs += durMs
	cat.AvgLatencyMs = cat.TotalDurationMs / float64(cat.TotalRequests)

	if cat.MinLatencyMs == 0 || durMs < cat.MinLatencyMs {
		cat.MinLatencyMs = durMs
	}
	if durMs > cat.MaxLatencyMs {
		cat.MaxLatencyMs = durMs
	}

	if statusCode >= 200 && statusCode < 400 {
		cat.SuccessCount++
	} else if statusCode >= 400 && statusCode < 500 {
		cat.ClientErrorCount++
	} else if statusCode >= 500 {
		cat.ServerErrorCount++
	}

	totalErrors := cat.ClientErrorCount + cat.ServerErrorCount
	if cat.TotalRequests > 0 {
		cat.ErrorRate = (float64(totalErrors) / float64(cat.TotalRequests)) * 100.0
	}

	// Status determination
	if cat.ServerErrorCount > 0 && cat.ErrorRate > 30.0 {
		cat.Status = "down"
	} else if cat.ErrorRate > 10.0 || cat.AvgLatencyMs > 1500.0 {
		cat.Status = "degraded"
	} else {
		cat.Status = "healthy"
	}

	cat.LastStatusCode = statusCode
	now := time.Now()
	cat.LastRequestAt = &now
}

func (c *MetricsCollector) GetReport(db *gorm.DB) *MonitoringReport {
	c.mu.RLock()
	defer c.mu.RUnlock()

	now := time.Now()
	uptimeSec := int64(now.Sub(c.startTime).Seconds())
	uptimeMinutes := float64(uptimeSec) / 60.0
	if uptimeMinutes <= 0 {
		uptimeMinutes = 1.0
	}

	totalErr := c.total4xx + c.total5xx
	var globalErrorRate float64
	var avgLatency float64
	if c.totalReq > 0 {
		globalErrorRate = (float64(totalErr) / float64(c.totalReq)) * 100.0
		avgLatency = c.totalDurMs / float64(c.totalReq)
	}

	reqPerMin := float64(c.totalReq) / uptimeMinutes

	overallStatus := "healthy"
	if c.total5xx > 0 && globalErrorRate > 25.0 {
		overallStatus = "down"
	} else if globalErrorRate > 10.0 || avgLatency > 1500.0 {
		overallStatus = "degraded"
	}

	// Format uptime string
	days := uptimeSec / 86400
	hours := (uptimeSec % 86400) / 3600
	mins := (uptimeSec % 3600) / 60
	secs := uptimeSec % 60
	uptimeStr := fmt.Sprintf("%dd %dh %dm %ds", days, hours, mins, secs)
	if days == 0 && hours == 0 {
		uptimeStr = fmt.Sprintf("%dm %ds", mins, secs)
	}

	overall := OverallStats{
		Status:            overallStatus,
		UptimeSeconds:     uptimeSec,
		UptimeFormatted:   uptimeStr,
		StartTime:         c.startTime,
		TotalRequests:     c.totalReq,
		Total2xx:          c.total2xx,
		Total3xx:          c.total3xx,
		Total4xx:          c.total4xx,
		Total5xx:          c.total5xx,
		GlobalErrorRate:   globalErrorRate,
		AvgLatencyMs:      avgLatency,
		RequestsPerMinute: reqPerMin,
	}

	// Categories clone
	catList := make([]CategoryStats, 0, len(c.categories))
	for _, def := range categoryDefinitions {
		if cat, ok := c.categories[def.Key]; ok {
			catCopy := *cat
			catList = append(catList, catCopy)
		}
	}

	// Go Runtime & Memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	system := SystemStats{
		GoVersion:         runtime.Version(),
		Goroutines:        runtime.NumGoroutine(),
		MemoryAllocatedMB: float64(m.Alloc) / 1024.0 / 1024.0,
		TotalAllocatedMB:  float64(m.TotalAlloc) / 1024.0 / 1024.0,
		SystemMemoryMB:    float64(m.Sys) / 1024.0 / 1024.0,
		NumGC:             m.NumGC,
	}

	// Database Connection Pool stats
	dbStats := DatabaseStats{
		Status: "down",
	}
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			startPing := time.Now()
			if pingErr := sqlDB.Ping(); pingErr == nil {
				dbStats.Status = "healthy"
				dbStats.PingMs = float64(time.Since(startPing).Microseconds()) / 1000.0
			}
			poolStats := sqlDB.Stats()
			dbStats.OpenConnections = poolStats.OpenConnections
			dbStats.InUse = poolStats.InUse
			dbStats.Idle = poolStats.Idle
			dbStats.WaitCount = poolStats.WaitCount
			dbStats.MaxOpenConnections = poolStats.MaxOpenConnections
		}
	}

	return &MonitoringReport{
		Overall:    overall,
		Categories: catList,
		Database:   dbStats,
		System:     system,
		Timestamp:  now,
	}
}

// ToPrometheusFormat produces standardized Prometheus / OpenMetrics plain text
func (c *MetricsCollector) ToPrometheusFormat(db *gorm.DB) string {
	report := c.GetReport(db)
	var sb strings.Builder

	sb.WriteString("# HELP http_requests_total Total number of HTTP requests processed.\n")
	sb.WriteString("# TYPE http_requests_total counter\n")
	sb.WriteString(fmt.Sprintf("http_requests_total %d\n\n", report.Overall.TotalRequests))

	sb.WriteString("# HELP http_requests_by_status_total Total HTTP requests partitioned by status code family.\n")
	sb.WriteString("# TYPE http_requests_by_status_total counter\n")
	sb.WriteString(fmt.Sprintf("http_requests_by_status_total{code=\"2xx\"} %d\n", report.Overall.Total2xx))
	sb.WriteString(fmt.Sprintf("http_requests_by_status_total{code=\"3xx\"} %d\n", report.Overall.Total3xx))
	sb.WriteString(fmt.Sprintf("http_requests_by_status_total{code=\"4xx\"} %d\n", report.Overall.Total4xx))
	sb.WriteString(fmt.Sprintf("http_requests_by_status_total{code=\"5xx\"} %d\n\n", report.Overall.Total5xx))

	sb.WriteString("# HELP http_category_requests_total Total requests processed per API category.\n")
	sb.WriteString("# TYPE http_category_requests_total counter\n")
	for _, cat := range report.Categories {
		sb.WriteString(fmt.Sprintf("http_category_requests_total{category=\"%s\"} %d\n", cat.Category, cat.TotalRequests))
	}
	sb.WriteString("\n")

	sb.WriteString("# HELP http_category_errors_total Total error responses (4xx and 5xx) per API category.\n")
	sb.WriteString("# TYPE http_category_errors_total counter\n")
	for _, cat := range report.Categories {
		errors := cat.ClientErrorCount + cat.ServerErrorCount
		sb.WriteString(fmt.Sprintf("http_category_errors_total{category=\"%s\"} %d\n", cat.Category, errors))
	}
	sb.WriteString("\n")

	sb.WriteString("# HELP http_category_latency_ms Average latency in milliseconds per API category.\n")
	sb.WriteString("# TYPE http_category_latency_ms gauge\n")
	for _, cat := range report.Categories {
		sb.WriteString(fmt.Sprintf("http_category_latency_ms{category=\"%s\"} %.2f\n", cat.Category, cat.AvgLatencyMs))
	}
	sb.WriteString("\n")

	sb.WriteString("# HELP http_category_status Health status indicator per API category (1=healthy, 0.5=degraded, 0=down).\n")
	sb.WriteString("# TYPE http_category_status gauge\n")
	for _, cat := range report.Categories {
		val := 1.0
		if cat.Status == "degraded" {
			val = 0.5
		} else if cat.Status == "down" {
			val = 0.0
		}
		sb.WriteString(fmt.Sprintf("http_category_status{category=\"%s\"} %.1f\n", cat.Category, val))
	}
	sb.WriteString("\n")

	sb.WriteString("# HELP go_goroutines Number of active goroutines.\n")
	sb.WriteString("# TYPE go_goroutines gauge\n")
	sb.WriteString(fmt.Sprintf("go_goroutines %d\n\n", report.System.Goroutines))

	sb.WriteString("# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.\n")
	sb.WriteString("# TYPE go_memstats_alloc_bytes gauge\n")
	sb.WriteString(fmt.Sprintf("go_memstats_alloc_bytes %.0f\n\n", report.System.MemoryAllocatedMB*1024*1024))

	sb.WriteString("# HELP db_connection_status Database connection status (1=healthy, 0=down).\n")
	sb.WriteString("# TYPE db_connection_status gauge\n")
	dbVal := 0
	if report.Database.Status == "healthy" {
		dbVal = 1
	}
	sb.WriteString(fmt.Sprintf("db_connection_status %d\n\n", dbVal))

	sb.WriteString("# HELP db_open_connections Current open database connections.\n")
	sb.WriteString("# TYPE db_open_connections gauge\n")
	sb.WriteString(fmt.Sprintf("db_open_connections %d\n\n", report.Database.OpenConnections))

	sb.WriteString("# HELP db_ping_latency_ms Database ping latency in milliseconds.\n")
	sb.WriteString("# TYPE db_ping_latency_ms gauge\n")
	sb.WriteString(fmt.Sprintf("db_ping_latency_ms %.2f\n", report.Database.PingMs))

	return sb.String()
}
