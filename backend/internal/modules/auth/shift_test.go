package auth

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// Helper function to test the shift window logic used in authService.checkShiftAccess
func isWithinShiftWindow(now time.Time, startTime, endTime string, isOvernight bool, toleranceMinutes int) bool {
	loc := now.Location()
	startParts := strings.Split(startTime, ":")
	endParts := strings.Split(endTime, ":")
	if len(startParts) < 2 || len(endParts) < 2 {
		return false
	}

	startHour, _ := strconv.Atoi(startParts[0])
	startMin, _ := strconv.Atoi(startParts[1])
	endHour, _ := strconv.Atoi(endParts[0])
	endMin, _ := strconv.Atoi(endParts[1])

	shiftStart := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMin, 0, 0, loc)
	shiftEnd := time.Date(now.Year(), now.Month(), now.Day(), endHour, endMin, 0, 0, loc)

	if isOvernight || shiftEnd.Before(shiftStart) {
		if now.Hour() < endHour || (now.Hour() == endHour && now.Minute() <= endMin) {
			shiftStart = shiftStart.AddDate(0, 0, -1)
		} else {
			shiftEnd = shiftEnd.AddDate(0, 0, 1)
		}
	}

	allowedStart := shiftStart.Add(-time.Duration(toleranceMinutes) * time.Minute)
	allowedEnd := shiftEnd.Add(time.Duration(toleranceMinutes) * time.Minute)

	return !now.Before(allowedStart) && !now.After(allowedEnd)
}

func TestShiftAccessLogic(t *testing.T) {
	loc := time.FixedZone("WIB", 7*3600)

	t.Run("Normal shift (07:00 - 15:00): exactly inside shift", func(t *testing.T) {
		now := time.Date(2026, 11, 17, 10, 0, 0, 0, loc)
		if !isWithinShiftWindow(now, "07:00", "15:00", false, 30) {
			t.Errorf("expected 10:00 to be within 07:00 - 15:00 shift")
		}
	})

	t.Run("Normal shift: 20 minutes before shift (within 30m tolerance)", func(t *testing.T) {
		now := time.Date(2026, 11, 17, 6, 40, 0, 0, loc)
		if !isWithinShiftWindow(now, "07:00", "15:00", false, 30) {
			t.Errorf("expected 06:40 to be within 30m tolerance of 07:00 shift")
		}
	})

	t.Run("Normal shift: 45 minutes before shift (outside tolerance)", func(t *testing.T) {
		now := time.Date(2026, 11, 17, 6, 15, 0, 0, loc)
		if isWithinShiftWindow(now, "07:00", "15:00", false, 30) {
			t.Errorf("expected 06:15 to be OUTSIDE tolerance of 07:00 shift")
		}
	})

	t.Run("Normal shift: 20 minutes after shift (within 30m handover tolerance)", func(t *testing.T) {
		now := time.Date(2026, 11, 17, 15, 20, 0, 0, loc)
		if !isWithinShiftWindow(now, "07:00", "15:00", false, 30) {
			t.Errorf("expected 15:20 to be within 30m closing tolerance of 15:00 shift")
		}
	})

	t.Run("Normal shift: 2 hours after shift ended", func(t *testing.T) {
		now := time.Date(2026, 11, 17, 17, 0, 0, 0, loc)
		if isWithinShiftWindow(now, "07:00", "15:00", false, 30) {
			t.Errorf("expected 17:00 to be LOCKED OUT of 07:00-15:00 shift")
		}
	})

	t.Run("Overnight shift (23:00 - 07:00): at 02:00 AM", func(t *testing.T) {
		now := time.Date(2026, 11, 17, 2, 0, 0, 0, loc)
		if !isWithinShiftWindow(now, "23:00", "07:00", true, 30) {
			t.Errorf("expected 02:00 AM to be within overnight 23:00 - 07:00 shift")
		}
	})

	t.Run("Overnight shift (23:00 - 07:00): at 23:30 PM", func(t *testing.T) {
		now := time.Date(2026, 11, 17, 23, 30, 0, 0, loc)
		if !isWithinShiftWindow(now, "23:00", "07:00", true, 30) {
			t.Errorf("expected 23:30 PM to be within overnight 23:00 - 07:00 shift")
		}
	})

	t.Run("Overnight shift (23:00 - 07:00): at 12:00 PM noon (locked out)", func(t *testing.T) {
		now := time.Date(2026, 11, 17, 12, 0, 0, 0, loc)
		if isWithinShiftWindow(now, "23:00", "07:00", true, 30) {
			t.Errorf("expected 12:00 PM to be locked out of overnight shift")
		}
	})
}
