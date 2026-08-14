package types

import (
	"fmt"
	"strings"
	"time"
)

const dateOnlyLayout = "2006-01-02"

// DateOnly adalah wrapper time.Time yang hanya menerima format tanggal "YYYY-MM-DD"
// dari JSON request, tanpa komponen jam/menit/detik/timezone.
type DateOnly struct {
	time.Time
}

// UnmarshalJSON mem-parsing string "YYYY-MM-DD" dari JSON
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse(dateOnlyLayout, s)
	if err != nil {
		return fmt.Errorf("format tanggal harus YYYY-MM-DD: %w", err)
	}
	d.Time = t
	return nil
}

// MarshalJSON mengeluarkan string "YYYY-MM-DD"
func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d.Time.Format(dateOnlyLayout))), nil
}

// ToTimePtr mengonversi *DateOnly menjadi *time.Time, dipakai saat mapping DTO -> model
func (d *DateOnly) ToTimePtr() *time.Time {
	if d == nil || d.Time.IsZero() {
		return nil
	}
	t := d.Time
	return &t
}

// NewDateOnlyPtr mengonversi *time.Time menjadi *DateOnly, dipakai saat mapping model -> response
func NewDateOnlyPtr(t *time.Time) *DateOnly {
	if t == nil || t.IsZero() {
		return nil
	}
	return &DateOnly{Time: *t}
}
