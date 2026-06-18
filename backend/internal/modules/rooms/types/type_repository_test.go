package types

import (
	"testing"

	"gorm.io/gorm"
)

func TestNewTypeRepository(t *testing.T) {
	var db *gorm.DB
	repo := NewTypeRepository(db)
	if repo == nil {
		t.Error("expected NewTypeRepository to return a non-nil interface implementation")
	}
}
