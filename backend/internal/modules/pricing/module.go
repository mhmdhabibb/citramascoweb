package pricing

import (
	"citramascoweb-backend/internal/middlewares"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PricingModule struct {
	handler *PricingHandler
}

func NewPricingModule(db *gorm.DB) *PricingModule {
	repo := NewPricingRepository(db)
	service := NewPricingService(repo)
	handler := NewPricingHandler(service)
	return &PricingModule{handler: handler}
}

func (m *PricingModule) RegisterRoutes(r *gin.RouterGroup) {
	pricing := r.Group("/pricing")
	pricing.Use(middlewares.AuthMiddleware())

	// Seasons
	pricing.GET("/seasons", m.handler.GetSeasons)
	pricing.POST("/seasons", middlewares.RoleMiddleware("admin", "manager"), m.handler.CreateSeason)
	pricing.PUT("/seasons/:id", middlewares.RoleMiddleware("admin", "manager"), m.handler.UpdateSeason)
	pricing.DELETE("/seasons/:id", middlewares.RoleMiddleware("admin", "manager"), m.handler.DeleteSeason)

	// Daily Rates
	pricing.GET("/daily-rates", m.handler.GetRates)
	pricing.POST("/daily-rates", middlewares.RoleMiddleware("admin", "manager"), m.handler.SetRate)

	// Market Snapshots & Competitors
	pricing.GET("/snapshots", m.handler.GetSnapshots)
	pricing.GET("/competitor-rates", m.handler.GetCompetitorRates)

	// Events
	pricing.GET("/events", m.handler.GetMarketEvents)
	pricing.POST("/events", middlewares.RoleMiddleware("admin", "manager"), m.handler.CreateMarketEvent)
	pricing.DELETE("/events/:id", middlewares.RoleMiddleware("admin", "manager"), m.handler.DeleteMarketEvent)

	// KPI
	pricing.GET("/kpi-summaries", m.handler.GetKPI)
}

func SeedDefaultSeasons(db *gorm.DB) {
	year := time.Now().Year()
	defaults := []struct {
		Name             string
		SeasonType       string
		StartDate        time.Time
		EndDate          time.Time
		DemandMultiplier float64
	}{
		{Name: "High Season Libur Lebaran", SeasonType: "peak_season", StartDate: time.Date(year, 3, 25, 0, 0, 0, 0, time.UTC), EndDate: time.Date(year, 4, 15, 0, 0, 0, 0, time.UTC), DemandMultiplier: 1.35},
		{Name: "High Season Libur Sekolah", SeasonType: "high_season", StartDate: time.Date(year, 6, 15, 0, 0, 0, 0, time.UTC), EndDate: time.Date(year, 7, 20, 0, 0, 0, 0, time.UTC), DemandMultiplier: 1.20},
		{Name: "Peak Season Akhir Tahun & Tahun Baru", SeasonType: "peak_season", StartDate: time.Date(year, 12, 20, 0, 0, 0, 0, time.UTC), EndDate: time.Date(year+1, 1, 5, 0, 0, 0, 0, time.UTC), DemandMultiplier: 1.50},
	}

	for _, s := range defaults {
		var count int64
		db.Model(&CalendarSeason{}).Where("name = ?", s.Name).Count(&count)
		if count == 0 {
			db.Create(&CalendarSeason{
				Id:               uuid.New().String(),
				Name:             s.Name,
				SeasonType:       s.SeasonType,
				StartDate:        s.StartDate,
				EndDate:          s.EndDate,
				DemandMultiplier: s.DemandMultiplier,
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			})
		}
	}
}
