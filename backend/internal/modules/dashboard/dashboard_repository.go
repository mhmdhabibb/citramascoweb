package dashboard

import (
	"citramascoweb-backend/internal/dto"
	"time"

	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetCardStats(startDate, endDate time.Time) (total, occupied, revenue int64, err error)
	GetBookingBreakdown(startDate, endDate time.Time) (*dto.BookingBreakdown, error)
	GetMonthlyRevenueChart(startDate, endDate time.Time) ([]dto.RevenueMonthly, error)
}

type dashboardRepo struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepo{db: db}
}

// 1. Ambil data untuk 3 Kartu Utama atas (Occupancy, Available, Revenue)
func (r *dashboardRepo) GetCardStats(startDate, endDate time.Time) (total, occupied, revenue int64, err error) {
	// Hitung Total Kamar Fisik
	if err = r.db.Table("rooms").Count(&total).Error; err != nil {
		return
	}

	// Hitung Kamar Terisi (Occupied) dalam rentang waktu bersangkutan
	if total > 0 {
		err = r.db.Table("reservations").
			Where("status IN ?", []string{"approved", "checked-in"}). // checked-out sudah keluar kamar, tidak dihitung occupied
			Where("checkin_date < ? AND checkout_date > ?", endDate, startDate).
			Distinct("room_id").
			Count(&occupied).Error
		if err != nil {
			return
		}
	}

	// Hitung Total Revenue dari reservasi yang check-in pada periode tersebut
	err = r.db.Table("reservations").
		Where("status IN ?", []string{"approved", "checked-in", "checked-out"}).
		Where("checkin_date >= ? AND checkin_date <= ?", startDate, endDate).
		Select("COALESCE(SUM(total_price), 0)").
		Row().
		Scan(&revenue)

	return
}

// 2. Ambil data untuk Pie Chart (Booking Breakdown) sesuai rentang waktu terpilih
// 2. Ambil data untuk Pie Chart (Booking Breakdown) dengan 6 Status Lengkap
// 2. Ambil data untuk Pie Chart (Booking Breakdown) dengan Query yang Terisolasi
func (r *dashboardRepo) GetBookingBreakdown(startDate, endDate time.Time) (*dto.BookingBreakdown, error) {
	var breakdown dto.BookingBreakdown

	var bookedCount, canceledCount, pendingCount, rejectedCount int64

	// Gunakan r.db langsung untuk tiap status agar query Where tidak menumpuk/saling menjegal

	err := r.db.Table("reservations").
		Where("TRIM(LOWER(status)) IN ?", []string{"approve", "approved", "checked-in", "checked-out"}).
		Count(&bookedCount).Error
	if err != nil {
		return nil, err
	}

	// 2. Hitung Cancel & Rejected (Gunakan LOWER dan TRIM)
	err = r.db.Table("reservations").
		Where("TRIM(LOWER(status)) = ?", "cancel").
		Count(&canceledCount).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Table("reservations").
		Where("TRIM(LOWER(status)) = ?", "rejected").
		Count(&rejectedCount).Error
	if err != nil {
		return nil, err
	}

	// 3. Hitung Pending (Paksa TRIM dan LOWER)
	err = r.db.Table("reservations").
		Where("TRIM(LOWER(status)) = ?", "pending").
		Count(&pendingCount).Error
	if err != nil {
		return nil, err
	}

	// Masukkan ke properti struct DTO
	breakdown.Booked = bookedCount
	breakdown.Canceled = canceledCount + rejectedCount
	breakdown.PendingConfirmation = pendingCount

	// Total kalkulasi seluruh status
	breakdown.TotalBookings = breakdown.Booked + breakdown.Canceled + breakdown.PendingConfirmation

	// Hitung Persentase masing-masing untuk UI chart
	if breakdown.TotalBookings > 0 {
		totalFloat := float64(breakdown.TotalBookings)
		breakdown.BookedPercentage = (float64(breakdown.Booked) / totalFloat) * 100
		breakdown.CanceledPercentage = (float64(breakdown.Canceled) / totalFloat) * 100
		breakdown.PendingPercentage = (float64(breakdown.PendingConfirmation) / totalFloat) * 100
	}

	return &breakdown, nil
}

// 3. Ambil data untuk Line Chart (Pendapatan Bulanan berjalan)
func (r *dashboardRepo) GetMonthlyRevenueChart(startDate, endDate time.Time) ([]dto.RevenueMonthly, error) {
	var chart []dto.RevenueMonthly

	// Query fleksibel PostgreSQL yang dibatasi rentang waktu parameter dashboard
	err := r.db.Table("reservations").
		Select("TO_CHAR(checkin_date, 'Mon') as month, SUM(total_price) as income").
		Where("status IN ?", []string{"approved", "checked-in", "checked-out"}).
		Where("checkin_date >= ? AND checkin_date <= ?", startDate, endDate). // Filter range waktu chart
		Group("TO_CHAR(checkin_date, 'Mon'), EXTRACT(MONTH FROM checkin_date)").
		Order("EXTRACT(MONTH FROM checkin_date) ASC").
		Scan(&chart).Error

	return chart, err
}
