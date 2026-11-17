package pricing

import (
	"time"
)

type CalendarSeason struct {
	Id               string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	Name             string    `gorm:"type:varchar(100);not null" json:"name"`
	SeasonType       string    `gorm:"type:varchar(50);not null" json:"season_type"` // high_season, low_season, peak_season, normal
	StartDate        time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate          time.Time `gorm:"type:date;not null" json:"end_date"`
	DemandMultiplier float64   `gorm:"type:decimal(4,2);default:1.00" json:"demand_multiplier"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (CalendarSeason) TableName() string {
	return "calendar_seasons"
}

type DailyRoomRate struct {
	Id               string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	RoomTypeId       string    `gorm:"type:varchar(191);index:idx_room_rate_date;not null" json:"room_type_id"`
	Date             time.Time `gorm:"type:date;index:idx_room_rate_date;not null" json:"date"`
	BasePrice        float64   `gorm:"type:decimal(12,2);not null" json:"base_price"`
	RecommendedPrice float64   `gorm:"type:decimal(12,2);not null" json:"recommended_price"`
	FinalPrice       float64   `gorm:"type:decimal(12,2);not null" json:"final_price"`
	IsLocked         bool      `gorm:"default:false" json:"is_locked"`
	LockedByUserId   *string   `gorm:"type:varchar(191)" json:"locked_by_user_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (DailyRoomRate) TableName() string {
	return "daily_room_rates"
}

type DailyMarketSnapshot struct {
	Id                 string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	SnapshotDate       time.Time `gorm:"type:date;index:idx_snapshot_date_room;not null" json:"snapshot_date"`
	RoomTypeId         *string   `gorm:"type:varchar(191);index:idx_snapshot_date_room" json:"room_type_id"`
	TotalRooms         int       `gorm:"default:0" json:"total_rooms"`
	OooRooms           int       `gorm:"default:0" json:"ooo_rooms"`
	AvailableRooms     int       `gorm:"default:0" json:"available_rooms"`
	SoldRooms          int       `gorm:"default:0" json:"sold_rooms"`
	RoomRevenue        float64   `gorm:"type:decimal(14,2);default:0" json:"room_revenue"`
	OccupancyRate      float64   `gorm:"type:decimal(5,2);default:0" json:"occupancy_rate"`
	Adr                float64   `gorm:"type:decimal(12,2);default:0" json:"adr"`    // Average Daily Rate
	Revpar             float64   `gorm:"type:decimal(12,2);default:0" json:"revpar"` // Revenue Per Available Room
	AvgLeadTimeDays    float64   `gorm:"type:decimal(5,2);default:0" json:"avg_lead_time_days"`
	CancellationsCount int       `gorm:"default:0" json:"cancellations_count"`
	NoShowCount        int       `gorm:"default:0" json:"no_show_count"`
	CreatedAt          time.Time `json:"created_at"`
}

func (DailyMarketSnapshot) TableName() string {
	return "daily_market_snapshots"
}

type CompetitorRate struct {
	Id             string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	CompetitorName string    `gorm:"type:varchar(100);not null" json:"competitor_name"`
	RoomTypeName   string    `gorm:"type:varchar(100);not null" json:"room_type_name"`
	CheckInDate    time.Time `gorm:"type:date;not null" json:"check_in_date"`
	ScrapedAt      time.Time `json:"scraped_at"`
	Price          float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	CreatedAt      time.Time `json:"created_at"`
}

func (CompetitorRate) TableName() string {
	return "competitor_rates"
}

type MarketEvent struct {
	Id                   string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	Name                 string    `gorm:"type:varchar(100);not null" json:"name"`
	EventDate            time.Time `gorm:"type:date;not null" json:"event_date"`
	ImpactLevel          string    `gorm:"type:varchar(20);default:'medium'" json:"impact_level"` // high, medium, low
	EstimatedDemandBoost float64   `gorm:"type:decimal(5,2);default:0" json:"estimated_demand_boost"`
	Notes                string    `gorm:"type:text" json:"notes"`
	CreatedAt            time.Time `json:"created_at"`
}

func (MarketEvent) TableName() string {
	return "market_events"
}

type MonthlyKPISummary struct {
	Id                     string    `gorm:"type:varchar(191);primaryKey" json:"id"`
	Year                   int       `gorm:"index:idx_kpi_year_month;not null" json:"year"`
	Month                  int       `gorm:"index:idx_kpi_year_month;not null" json:"month"`
	RoomTypeId             *string   `gorm:"type:varchar(191)" json:"room_type_id"`
	TargetRevenue          float64   `gorm:"type:decimal(14,2);default:0" json:"target_revenue"`
	ActualRevenue          float64   `gorm:"type:decimal(14,2);default:0" json:"actual_revenue"`
	RevenueAchievementRate float64   `gorm:"type:decimal(5,2);default:0" json:"revenue_achievement_rate"`
	AvgOccupancyRate       float64   `gorm:"type:decimal(5,2);default:0" json:"avg_occupancy_rate"`
	AvgAdr                 float64   `gorm:"type:decimal(12,2);default:0" json:"avg_adr"`
	AvgRevpar              float64   `gorm:"type:decimal(12,2);default:0" json:"avg_revpar"`
	DirectBookingRatio     float64   `gorm:"type:decimal(5,2);default:0" json:"direct_booking_ratio"`
	AvgCommissionRate      float64   `gorm:"type:decimal(5,2);default:0" json:"avg_commission_rate"`
	TotalCancellations     int       `gorm:"default:0" json:"total_cancellations"`
	CalculatedAt           time.Time `json:"calculated_at"`
}

func (MonthlyKPISummary) TableName() string {
	return "monthly_kpi_summaries"
}
