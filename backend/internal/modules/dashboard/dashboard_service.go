package dashboard

import (
	"citramascoweb-backend/internal/dto"
	"sync"
	"time"
)

type DashboardService interface {
	GetDashboardData(param dto.DashboardParam) (*dto.DashboardResponse, error)
}

type dashboardService struct {
	repo DashboardRepository
}

func NewDashboardService(repo DashboardRepository) DashboardService {
	return &dashboardService{repo: repo}
}

func (s *dashboardService) GetDashboardData(param dto.DashboardParam) (*dto.DashboardResponse, error) {
	var startDate, endDate time.Time
	now := time.Now()

	// 1. LOGIKA FILTER RENTANG WAKTU (Sudah Sempurna)
	switch param.Range {
	case "today":
		startDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		endDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())

	case "weekly":
		daysSinceMonday := int(now.Weekday()) - 1
		if daysSinceMonday < 0 {
			daysSinceMonday = 6
		}
		monday := now.AddDate(0, 0, -daysSinceMonday)
		startDate = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())
		sunday := startDate.AddDate(0, 0, 6)
		endDate = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 999999999, now.Location())

	case "monthly":
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endDate = startDate.AddDate(0, 1, 0).Add(-time.Second)

	default:
		// Default dialihkan ke "weekly" agar saat load awal halaman, chart terisi data trend seminggu berjalan
		daysSinceMonday := int(now.Weekday()) - 1
		if daysSinceMonday < 0 {
			daysSinceMonday = 6
		}
		monday := now.AddDate(0, 0, -daysSinceMonday)
		startDate = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())
		sunday := startDate.AddDate(0, 0, 6)
		endDate = time.Date(sunday.Year(), sunday.Month(), sunday.Day(), 23, 59, 59, 999999999, now.Location())
	}

	var response dto.DashboardResponse
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	wg.Add(3)

	// Goroutine 1: Hitung Angka Kartu (Occupancy, Available, Revenue)
	go func() {
		defer wg.Done()
		total, occupied, revenue, err := s.repo.GetCardStats(startDate, endDate)

		mu.Lock()
		if err != nil {
			errs = append(errs, err)
		} else {
			response.TotalRooms = total
			response.TotalRevenue = revenue
			response.AvailableRooms = total - occupied
			if total > 0 {
				response.OccupancyRate = (float64(occupied) / float64(total)) * 100
			}
		}
		mu.Unlock()
	}()

	// Goroutine 2: Hitung Booking Breakdown (Donut Chart) -> SEKARANG MENGGUNAKAN RANGE WAKTU
	go func() {
		defer wg.Done()
		breakdown, err := s.repo.GetBookingBreakdown(startDate, endDate)

		mu.Lock()
		if err != nil {
			errs = append(errs, err)
		} else if breakdown != nil {
			response.BookingBreakdown = *breakdown
		}
		mu.Unlock()
	}()

	// Goroutine 3: Hitung Grafik Pendapatan Bulanan (Area Chart) -> SEKARANG MENGGUNAKAN RANGE WAKTU
	go func() {
		defer wg.Done()
		chart, err := s.repo.GetMonthlyRevenueChart(startDate, endDate)

		mu.Lock()
		if err != nil {
			errs = append(errs, err)
		} else {
			response.RevenueChart = chart
		}
		mu.Unlock()
	}()

	// Tunggu sampai ketiga query paralel GORM selesai diproses serentak
	wg.Wait()

	// Jika ada salah satu query yang gagal, kembalikan error agar frontend mendapatkan sinyal
	if len(errs) > 0 {
		return nil, errs[0]
	}

	return &response, nil
}
