package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DateOnly time.Time

func (d *DateOnly) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}
	*d = DateOnly(t)
	return nil
}

func (d DateOnly) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	return json.Marshal(t.Format("2006-01-02"))
}

func (d DateOnly) ToTime() time.Time {
	return time.Time(d)
}

func NewDateOnly(t time.Time) DateOnly {
	return DateOnly(t.Truncate(24 * time.Hour))
}
