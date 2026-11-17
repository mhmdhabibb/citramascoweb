package pricing

import (
	"math"
	"testing"
)

func calculateDynamicRate(basePrice float64, multiplier float64) float64 {
	return math.Round(basePrice * multiplier)
}

func TestPricingLogic(t *testing.T) {
	t.Run("Normal season pricing (1.00x multiplier)", func(t *testing.T) {
		base := 500000.0
		multiplier := 1.00
		final := calculateDynamicRate(base, multiplier)

		if final != 500000.0 {
			t.Errorf("expected 500000.0, got %.2f", final)
		}
	})

	t.Run("Peak season Lebaran pricing (+35% / 1.35x)", func(t *testing.T) {
		base := 500000.0
		multiplier := 1.35
		final := calculateDynamicRate(base, multiplier)

		expected := 675000.0
		if final != expected {
			t.Errorf("expected %.2f, got %.2f", expected, final)
		}
	})

	t.Run("Low season discount pricing (-15% / 0.85x)", func(t *testing.T) {
		base := 500000.0
		multiplier := 0.85
		final := calculateDynamicRate(base, multiplier)

		expected := 425000.0
		if final != expected {
			t.Errorf("expected %.2f, got %.2f", expected, final)
		}
	})

	t.Run("Occupancy rate and RevPAR calculation", func(t *testing.T) {
		totalRooms := 50
		soldRooms := 40
		roomRevenue := 20000000.0 // 20 juta

		occupancyRate := (float64(soldRooms) / float64(totalRooms)) * 100.0
		adr := roomRevenue / float64(soldRooms)     // 500.000
		revpar := roomRevenue / float64(totalRooms) // 400.000

		if occupancyRate != 80.0 {
			t.Errorf("expected occupancy 80%%, got %.2f%%", occupancyRate)
		}
		if adr != 500000.0 {
			t.Errorf("expected ADR 500000, got %.2f", adr)
		}
		if revpar != 400000.0 {
			t.Errorf("expected RevPAR 400000, got %.2f", revpar)
		}
	})
}
