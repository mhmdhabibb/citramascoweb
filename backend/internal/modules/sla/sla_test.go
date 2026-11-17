package sla

import (
	"testing"
	"time"
)

func TestSLALogic(t *testing.T) {
	t.Run("Urgent task SLA compliance (target 30m)", func(t *testing.T) {
		createdAt := time.Now().Add(-20 * time.Minute)
		resolvedAt := time.Now()
		targetMinutes := 30

		actualMinutes := int(resolvedAt.Sub(createdAt).Minutes())
		isBreached := actualMinutes > targetMinutes

		if isBreached {
			t.Errorf("expected 20 minutes to not breach 30 minutes target")
		}
		if actualMinutes < 19 || actualMinutes > 21 {
			t.Errorf("expected around 20 minutes, got %d", actualMinutes)
		}
	})

	t.Run("High priority task SLA breached (target 45m, took 60m)", func(t *testing.T) {
		createdAt := time.Now().Add(-60 * time.Minute)
		resolvedAt := time.Now()
		targetMinutes := 45

		actualMinutes := int(resolvedAt.Sub(createdAt).Minutes())
		isBreached := actualMinutes > targetMinutes

		if !isBreached {
			t.Errorf("expected 60 minutes to breach 45 minutes target")
		}
	})

	t.Run("Staff SLA compliance rate calculation", func(t *testing.T) {
		totalCompleted := 10
		totalBreached := 2

		complianceRate := float64(totalCompleted-totalBreached) / float64(totalCompleted) * 100.0

		expectedRate := 80.0
		if complianceRate != expectedRate {
			t.Errorf("expected compliance rate %.2f%%, got %.2f%%", expectedRate, complianceRate)
		}
	})

	t.Run("Staff SLA average resolution time calculation", func(t *testing.T) {
		prevCompleted := 3
		prevAvg := 30.0     // minutes
		newDuration := 50.0 // minutes

		newCompleted := prevCompleted + 1
		newAvg := (prevAvg*float64(prevCompleted) + newDuration) / float64(newCompleted)

		expectedAvg := 35.0 // (90 + 50) / 4 = 140 / 4 = 35.0
		if newAvg != expectedAvg {
			t.Errorf("expected new average %.2f, got %.2f", expectedAvg, newAvg)
		}
	})
}
