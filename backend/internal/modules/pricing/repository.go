package pricing

import (
	"time"

	"gorm.io/gorm"
)

type PricingRepository interface {
	// Seasons
	FindSeasons() ([]CalendarSeason, error)
	FindSeasonById(id string) (*CalendarSeason, error)
	CreateSeason(s *CalendarSeason) error
	UpdateSeason(s *CalendarSeason) error
	DeleteSeason(id string) error

	// Daily Rates
	FindRatesByRange(roomTypeId string, startDate, endDate time.Time) ([]DailyRoomRate, error)
	SaveDailyRate(rate *DailyRoomRate) error

	// Market Snapshots
	FindSnapshots(startDate, endDate time.Time) ([]DailyMarketSnapshot, error)
	SaveSnapshot(snap *DailyMarketSnapshot) error

	// Competitor Rates
	FindCompetitorRates(checkInDate time.Time) ([]CompetitorRate, error)
	CreateCompetitorRate(cr *CompetitorRate) error

	// Market Events
	FindMarketEvents(startDate, endDate time.Time) ([]MarketEvent, error)
	CreateMarketEvent(me *MarketEvent) error
	DeleteMarketEvent(id string) error

	// KPI
	FindKPISummaries(year int) ([]MonthlyKPISummary, error)
	SaveKPISummary(kpi *MonthlyKPISummary) error
}

type pricingRepository struct {
	db *gorm.DB
}

func NewPricingRepository(db *gorm.DB) PricingRepository {
	return &pricingRepository{db: db}
}

func (r *pricingRepository) FindSeasons() ([]CalendarSeason, error) {
	var list []CalendarSeason
	err := r.db.Order("start_date ASC").Find(&list).Error
	return list, err
}

func (r *pricingRepository) FindSeasonById(id string) (*CalendarSeason, error) {
	var s CalendarSeason
	err := r.db.Where("id = ?", id).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *pricingRepository) CreateSeason(s *CalendarSeason) error {
	return r.db.Create(s).Error
}

func (r *pricingRepository) UpdateSeason(s *CalendarSeason) error {
	return r.db.Save(s).Error
}

func (r *pricingRepository) DeleteSeason(id string) error {
	return r.db.Where("id = ?", id).Delete(&CalendarSeason{}).Error
}

func (r *pricingRepository) FindRatesByRange(roomTypeId string, startDate, endDate time.Time) ([]DailyRoomRate, error) {
	var rates []DailyRoomRate
	q := r.db.Where("date >= ? AND date <= ?", startDate, endDate)
	if roomTypeId != "" {
		q = q.Where("room_type_id = ?", roomTypeId)
	}
	err := q.Order("date ASC").Find(&rates).Error
	return rates, err
}

func (r *pricingRepository) SaveDailyRate(rate *DailyRoomRate) error {
	var existing DailyRoomRate
	err := r.db.Where("room_type_id = ? AND date = ?", rate.RoomTypeId, rate.Date).First(&existing).Error
	if err == nil {
		rate.Id = existing.Id
		return r.db.Save(rate).Error
	}
	return r.db.Create(rate).Error
}

func (r *pricingRepository) FindSnapshots(startDate, endDate time.Time) ([]DailyMarketSnapshot, error) {
	var list []DailyMarketSnapshot
	err := r.db.Where("snapshot_date >= ? AND snapshot_date <= ?", startDate, endDate).Order("snapshot_date ASC").Find(&list).Error
	return list, err
}

func (r *pricingRepository) SaveSnapshot(snap *DailyMarketSnapshot) error {
	return r.db.Save(snap).Error
}

func (r *pricingRepository) FindCompetitorRates(checkInDate time.Time) ([]CompetitorRate, error) {
	var list []CompetitorRate
	q := r.db.Model(&CompetitorRate{})
	if !checkInDate.IsZero() {
		q = q.Where("check_in_date = ?", checkInDate)
	}
	err := q.Order("check_in_date ASC, competitor_name ASC").Find(&list).Error
	return list, err
}

func (r *pricingRepository) CreateCompetitorRate(cr *CompetitorRate) error {
	return r.db.Create(cr).Error
}

func (r *pricingRepository) FindMarketEvents(startDate, endDate time.Time) ([]MarketEvent, error) {
	var list []MarketEvent
	q := r.db.Model(&MarketEvent{})
	if !startDate.IsZero() && !endDate.IsZero() {
		q = q.Where("event_date >= ? AND event_date <= ?", startDate, endDate)
	}
	err := q.Order("event_date ASC").Find(&list).Error
	return list, err
}

func (r *pricingRepository) CreateMarketEvent(me *MarketEvent) error {
	return r.db.Create(me).Error
}

func (r *pricingRepository) DeleteMarketEvent(id string) error {
	return r.db.Where("id = ?", id).Delete(&MarketEvent{}).Error
}

func (r *pricingRepository) FindKPISummaries(year int) ([]MonthlyKPISummary, error) {
	var list []MonthlyKPISummary
	q := r.db.Model(&MonthlyKPISummary{})
	if year > 0 {
		q = q.Where("year = ?", year)
	}
	err := q.Order("year DESC, month ASC").Find(&list).Error
	return list, err
}

func (r *pricingRepository) SaveKPISummary(kpi *MonthlyKPISummary) error {
	return r.db.Save(kpi).Error
}
