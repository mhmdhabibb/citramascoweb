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

type CreateOfferRequest struct {
	Title       string      `json:"title" form:"title" binding:"required"`
	Image       *string     `json:"image" form:"-"`
	Price       *int        `json:"price" form:"price"`
	Discount    *int        `json:"discount" form:"discount"`
	Code        string      `json:"code" form:"code"`
	ValidStart  *CustomDate `json:"valid_start" form:"valid_start"`
	ValidEnd    *CustomDate `json:"valid_end" form:"valid_end"`
	MaxQuota    *int        `json:"max_quota" form:"max_quota"`
	Description *string     `json:"description" form:"description"`
}

type UpdateOfferRequest struct {
	Title string  `json:"title" form:"title"`
	Image *string `json:"image" form:"-"`
	Code  string  `json:"code" form:"code"`

	Price       *int        `json:"price" form:"price"`
	Discount    *int        `json:"discount" form:"discount"`
	ValidStart  *CustomDate `json:"valid_start" form:"valid_start"`
	ValidEnd    *CustomDate `json:"valid_end" form:"valid_end"`
	MaxQuota    *int        `json:"max_quota" form:"max_quota"`
	Description *string     `json:"description" form:"description"`
	Status      *string     `json:"status" form:"status"`
}
