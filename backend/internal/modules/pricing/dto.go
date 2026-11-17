package pricing

import "time"

type CreateSeasonRequest struct {
	Name             string    `json:"name" binding:"required"`
	SeasonType       string    `json:"season_type" binding:"required"`
	StartDate        time.Time `json:"start_date" binding:"required"`
	EndDate          time.Time `json:"end_date" binding:"required"`
	DemandMultiplier float64   `json:"demand_multiplier" binding:"required"`
}

type UpdateSeasonRequest struct {
	Name             string    `json:"name"`
	SeasonType       string    `json:"season_type"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	DemandMultiplier float64   `json:"demand_multiplier"`
}

type SetDailyRateRequest struct {
	RoomTypeId       string    `json:"room_type_id" binding:"required"`
	Date             time.Time `json:"date" binding:"required"`
	BasePrice        float64   `json:"base_price" binding:"required"`
	RecommendedPrice float64   `json:"recommended_price"`
	FinalPrice       float64   `json:"final_price" binding:"required"`
	IsLocked         bool      `json:"is_locked"`
}

type CreateMarketEventRequest struct {
	Name                 string    `json:"name" binding:"required"`
	EventDate            time.Time `json:"event_date" binding:"required"`
	ImpactLevel          string    `json:"impact_level"`
	EstimatedDemandBoost float64   `json:"estimated_demand_boost"`
	Notes                string    `json:"notes"`
}
