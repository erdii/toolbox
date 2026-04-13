// Package timeutil provides utilities for parsing and calculating durations
// with human-friendly time formats. It supports flexible time parsing that
// accepts various formats (12/24 hour, with/without dates) and calculates
// durations to the next occurrence of a given time.
package timeutil

import (
	"fmt"
	"time"
)

// ParseHumanTime parses a human-friendly time string into a time.Time value.
//
// It tries multiple common time formats in order, returning the first successful parse.
// Supported formats include:
//   - 24-hour times: "15:30", "15:30:45", "15"
//   - 12-hour times: "3:30PM", "3PM"
//   - Dates: "2.1.2006", "2.1.06", "2.1", "2006-01-02"
//   - Date+time combinations: "2006-01-02T15:30", "2006-01-02 15:30:45"
//   - RFC3339 formats
//
// Note: Single-digit format specifiers (like "15:4") work with both single
// and double-digit inputs, so "15:4" matches both "15:4" and "15:30".
//
// Returns an error if the input doesn't match any supported format.
func ParseHumanTime(input string) (time.Time, error) {
	for _, format := range []string{
		"15:4",
		"15:4:5",
		"15",
		"3:4PM",
		"3:4pm",
		"3PM",
		"3pm",
		"2.1.2006",
		"2.1.06",
		"2.1",
		"2006-01-02",
		"2006-01-02T15:4:5",
		"2006-01-02 15:4:5",
		"2006-01-02T15:4",
		"2006-01-02 15:4",
		time.RFC3339,
		time.RFC3339Nano,
	} {
		t, err := time.Parse(format, input)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse time: %s", input)
}
