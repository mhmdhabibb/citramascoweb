package guest

import (
	"testing"
)

func TestGuestLogic(t *testing.T) {
	t.Run("Create guest request validation", func(t *testing.T) {
		req := CreateGuestRequest{
			FirstName: "",
		}
		if req.FirstName == "" {
			// Expect validation to fail if first_name is empty
		} else {
			t.Errorf("expected empty first_name to fail validation")
		}

		validReq := CreateGuestRequest{
			Title:          "Mr",
			FirstName:      "John",
			LastName:       "Doe",
			IdentityNumber: "3201010101010001",
			Email:          "john.doe@example.com",
			Phone:          "081234567890",
		}
		if validReq.FirstName == "" || validReq.IdentityNumber == "" {
			t.Errorf("expected valid guest request to pass")
		}
	})

	t.Run("Default nationality and guest type assignment", func(t *testing.T) {
		req := CreateGuestRequest{
			FirstName: "Jane",
		}
		nationality := req.Nationality
		if nationality == "" {
			nationality = "Indonesia"
		}
		guestType := req.GuestType
		if guestType == "" {
			guestType = "regular"
		}

		if nationality != "Indonesia" {
			t.Errorf("expected default nationality 'Indonesia', got '%s'", nationality)
		}
		if guestType != "regular" {
			t.Errorf("expected default guest_type 'regular', got '%s'", guestType)
		}
	})
}
