package analytics

import (
	"fmt"
	"time"
)

// ParseDateRange parses optional UTC dates in YYYY-MM-DD format. startDate is
// inclusive and endDate is inclusive; an empty string leaves that side
// unbounded. It returns nil for unbounded sides and an exclusive end so the
// full end date is included. Invalid or reversed dates return an error.
func ParseDateRange(startDate, endDate string) (start, end *time.Time, err error) {
	if startDate != "" {
		value, parseErr := time.Parse(time.DateOnly, startDate)
		if parseErr != nil {
			return nil, nil, fmt.Errorf("start_date must be YYYY-MM-DD: %w", parseErr)
		}
		start = &value
	}
	if endDate != "" {
		value, parseErr := time.Parse(time.DateOnly, endDate)
		if parseErr != nil {
			return nil, nil, fmt.Errorf("end_date must be YYYY-MM-DD: %w", parseErr)
		}
		value = value.AddDate(0, 0, 1)
		end = &value
	}
	if err := validatePeriod(start, end); err != nil {
		return nil, nil, fmt.Errorf("end_date must not be before start_date")
	}
	return start, end, nil
}
