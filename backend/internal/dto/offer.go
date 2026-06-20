package dto

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type CustomDate time.Time

// UnmarshalJSON parses JSON strings to CustomDate
func (cd *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse("02-01-2006", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	*cd = CustomDate(t)
	return nil
}

// MarshalJSON marshals CustomDate to JSON string
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Format("02-01-2006"))), nil
}

// UnmarshalText implements encoding.TextUnmarshaler for form binding
func (cd *CustomDate) UnmarshalText(text []byte) error {
	s := string(text)
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse("02-01-2006", s)
	if err != nil {
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}
	*cd = CustomDate(t)
	return nil
}

// Scan implements sql.Scanner to read from database
func (cd *CustomDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("failed to scan CustomDate")
	}
	*cd = CustomDate(t)
	return nil
}

// Value implements driver.Valuer to write to database
func (cd CustomDate) Value() (driver.Value, error) {
	t := time.Time(cd)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

// ParseCustomDate parses a string representation into a CustomDate pointer
func ParseCustomDate(s string) (*CustomDate, error) {
	if s == "" || s == "null" {
		return nil, nil
	}
	// Try DD-MM-YYYY (e.g. 20-06-2026)
	t, err := time.Parse("02-01-2006", s)
	if err != nil {
		// Try YYYY-MM-DD (e.g. 2026-06-20)
		t, err = time.Parse("2006-01-02", s)
		if err != nil {
			// Try RFC3339
			t, err = time.Parse(time.RFC3339, s)
			if err != nil {
				// Try D-M-YYYY (e.g. 2-1-2006)
				t, err = time.Parse("2-1-2006", s)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	cd := CustomDate(t)
	return &cd, nil
}

type CreateOfferRequest struct {
	Title       string  `json:"title" form:"title" binding:"required"`
	SubTitle    *string `json:"sub_title" form:"sub_title"`
	Image       *string `json:"image" form:"-"`
	Price       *int    `json:"price" form:"price"`
	Discount    *int    `json:"discount" form:"discount"`
	Code        string  `json:"code" form:"code"`
	ValidStart  string  `json:"valid_start" form:"valid_start"`
	ValidEnd    string  `json:"valid_end" form:"valid_end"`
	MaxQuota    *int    `json:"max_quota" form:"max_quota"`
	Description *string `json:"description" form:"description"`
	Status      *string `json:"status" form:"status"`
}

type UpdateOfferRequest struct {
	Title    string  `json:"title" form:"title"`
	SubTitle *string `json:"sub_title" form:"sub_title"`
	Image    *string `json:"image" form:"-"`
	Code     string  `json:"code" form:"code"`

	Price       *int    `json:"price" form:"price"`
	Discount    *int    `json:"discount" form:"discount"`
	ValidStart  string  `json:"valid_start" form:"valid_start"`
	ValidEnd    string  `json:"valid_end" form:"valid_end"`
	MaxQuota    *int    `json:"max_quota" form:"max_quota"`
	Description *string `json:"description" form:"description"`
	Status      *string `json:"status" form:"status"`
}
