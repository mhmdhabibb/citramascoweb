package dto

type DashboardParam struct {
	// Menerima nilai: "today", "weekly", atau "monthly"
	Range string `form:"range"`
}

type BookingBreakdown struct {
	Booked              int64   `json:"booked"`
	Canceled            int64   `json:"canceled"`
	PendingConfirmation int64   `json:"pending_confirmation"`
	TotalBookings       int64   `json:"total_bookings"`
	BookedPercentage    float64 `json:"booked_percentage"`
	CanceledPercentage  float64 `json:"canceled_percentage"`
	PendingPercentage   float64 `json:"pending_percentage"`
}

type RevenueMonthly struct {
	Month  string `json:"month"` // e.g., "Jan", "Feb"
	Income int64  `json:"income"`
}

type DashboardResponse struct {
	OccupancyRate    float64          `json:"occupancy_rate"`
	AvailableRooms   int64            `json:"available_rooms"`
	TotalRooms       int64            `json:"total_rooms"`
	TotalRevenue     int64            `json:"total_revenue"`
	BookingBreakdown BookingBreakdown `json:"booking_breakdown"`
	RevenueChart     []RevenueMonthly `json:"revenue_chart"`
}
