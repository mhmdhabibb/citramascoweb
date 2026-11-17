package smartlock

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
	"time"
)

func generatePIN() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()+100000), nil
}

func TestSmartLockLogic(t *testing.T) {
	t.Run("Generate 6-digit PIN code", func(t *testing.T) {
		pin, err := generatePIN()
		if err != nil {
			t.Fatalf("failed to generate PIN: %v", err)
		}
		if len(pin) != 6 {
			t.Errorf("expected PIN length 6, got %d (%s)", len(pin), pin)
		}
		// ensure it contains digits only
		for _, ch := range pin {
			if ch < '0' || ch > '9' {
				t.Errorf("PIN contains non-digit character: %c", ch)
			}
		}
	})

	t.Run("Check validity period for guest room access", func(t *testing.T) {
		validFrom := time.Now().Add(-1 * time.Hour)
		validUntil := time.Now().Add(24 * time.Hour)

		now := time.Now()
		isAccessValid := now.After(validFrom) && now.Before(validUntil)

		if !isAccessValid {
			t.Errorf("expected access to be valid within booking window")
		}
	})

	t.Run("Expired credential check", func(t *testing.T) {
		validFrom := time.Now().Add(-48 * time.Hour)
		validUntil := time.Now().Add(-2 * time.Hour)

		now := time.Now()
		isAccessValid := now.After(validFrom) && now.Before(validUntil)

		if isAccessValid {
			t.Errorf("expected expired access to be invalid")
		}
	})
}
