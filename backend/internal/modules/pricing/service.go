package pricing

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type PricingService interface {
	GetSeasons() ([]CalendarSeason, error)
	GetSeasonById(id string) (*CalendarSeason, error)
	CreateSeason(req *CreateSeasonRequest) (*CalendarSeason, error)
	UpdateSeason(id string, req *UpdateSeasonRequest) (*CalendarSeason, error)
	DeleteSeason(id string) error

	GetRates(roomTypeId string, startDate, endDate time.Time) ([]DailyRoomRate, error)
	SetDailyRate(userId string, req *SetDailyRateRequest) (*DailyRoomRate, error)

	GetSnapshots(startDate, endDate time.Time) ([]DailyMarketSnapshot, error)
	GetCompetitorRates(checkInDate time.Time) ([]CompetitorRate, error)

	GetMarketEvents(startDate, endDate time.Time) ([]MarketEvent, error)
	CreateMarketEvent(req *CreateMarketEventRequest) (*MarketEvent, error)
	DeleteMarketEvent(id string) error

	GetKPISummaries(year int) ([]MonthlyKPISummary, error)
}

type pricingService struct {
	repo PricingRepository
}

func NewPricingService(repo PricingRepository) PricingService {
	return &pricingService{repo: repo}
}

func (s *pricingService) GetSeasons() ([]CalendarSeason, error) {
	return s.repo.FindSeasons()
}

func (s *pricingService) GetSeasonById(id string) (*CalendarSeason, error) {
	return s.repo.FindSeasonById(id)
}

func (s *pricingService) CreateSeason(req *CreateSeasonRequest) (*CalendarSeason, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	season := &CalendarSeason{
		Id:               uuid.New().String(),
		Name:             req.Name,
		SeasonType:       req.SeasonType,
		StartDate:        req.StartDate,
		EndDate:          req.EndDate,
		DemandMultiplier: req.DemandMultiplier,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if err := s.repo.CreateSeason(season); err != nil {
		return nil, err
	}
	return season, nil
}

func (s *pricingService) UpdateSeason(id string, req *UpdateSeasonRequest) (*CalendarSeason, error) {
	season, err := s.repo.FindSeasonById(id)
	if err != nil {
		return nil, err
	}
	if req.Name != "" {
		season.Name = req.Name
	}
	if req.SeasonType != "" {
		season.SeasonType = req.SeasonType
	}
	if !req.StartDate.IsZero() {
		season.StartDate = req.StartDate
	}
	if !req.EndDate.IsZero() {
		season.EndDate = req.EndDate
	}
	if req.DemandMultiplier > 0 {
		season.DemandMultiplier = req.DemandMultiplier
	}
	season.UpdatedAt = time.Now()
	if err := s.repo.UpdateSeason(season); err != nil {
		return nil, err
	}
	return season, nil
}

func (s *pricingService) DeleteSeason(id string) error {
	return s.repo.DeleteSeason(id)
}

func (s *pricingService) GetRates(roomTypeId string, startDate, endDate time.Time) ([]DailyRoomRate, error) {
	return s.repo.FindRatesByRange(roomTypeId, startDate, endDate)
}

func (s *pricingService) SetDailyRate(userId string, req *SetDailyRateRequest) (*DailyRoomRate, error) {
	rate := &DailyRoomRate{
		Id:               uuid.New().String(),
		RoomTypeId:       req.RoomTypeId,
		Date:             req.Date,
		BasePrice:        req.BasePrice,
		RecommendedPrice: req.RecommendedPrice,
		FinalPrice:       req.FinalPrice,
		IsLocked:         req.IsLocked,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
	if req.IsLocked {
		rate.LockedByUserId = &userId
	}
	if err := s.repo.SaveDailyRate(rate); err != nil {
		return nil, err
	}
	return rate, nil
}

func (s *pricingService) GetSnapshots(startDate, endDate time.Time) ([]DailyMarketSnapshot, error) {
	return s.repo.FindSnapshots(startDate, endDate)
}

func (s *pricingService) GetCompetitorRates(checkInDate time.Time) ([]CompetitorRate, error) {
	return s.repo.FindCompetitorRates(checkInDate)
}

func (s *pricingService) GetMarketEvents(startDate, endDate time.Time) ([]MarketEvent, error) {
	return s.repo.FindMarketEvents(startDate, endDate)
}

func (s *pricingService) CreateMarketEvent(req *CreateMarketEventRequest) (*MarketEvent, error) {
	me := &MarketEvent{
		Id:                   uuid.New().String(),
		Name:                 req.Name,
		EventDate:            req.EventDate,
		ImpactLevel:          req.ImpactLevel,
		EstimatedDemandBoost: req.EstimatedDemandBoost,
		Notes:                req.Notes,
		CreatedAt:            time.Now(),
	}
	if err := s.repo.CreateMarketEvent(me); err != nil {
		return nil, err
	}
	return me, nil
}

func (s *pricingService) DeleteMarketEvent(id string) error {
	return s.repo.DeleteMarketEvent(id)
}

func (s *pricingService) GetKPISummaries(year int) ([]MonthlyKPISummary, error) {
	if year == 0 {
		year = time.Now().Year()
	}
	return s.repo.FindKPISummaries(year)
}
