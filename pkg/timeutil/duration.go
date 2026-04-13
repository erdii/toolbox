package timeutil

import (
	"fmt"
	"time"
)

// DurationUntil calculates the duration from now until the next occurrence of the given time.
//
// If given is in the future (after now), it returns the difference directly.
// If given appears to be a time-only value (year is 0 or 1970) or is in the past,
// it calculates the duration until the next occurrence of that time:
//   - If the time has already passed today, it assumes tomorrow
//   - If the time hasn't passed today, it assumes today
//
// If given is a specific past date (year > 1970), it returns an error since
// "next occurrence" doesn't make sense for a specific date in the past.
//
// This is useful for calculating sleep durations until a specific clock time,
// regardless of whether the user specified a date or just a time.
//
// Example:
//   - now is 14:30, given is 16:00 (today) → returns 1h30m
//   - now is 14:30, given is 09:00 (time-only) → returns ~19h (tomorrow at 09:00)
//   - now is 14:30, given is future date/time → returns exact difference
//   - now is 2026-04-12, given is 2026-04-11 → returns error (past date)
func DurationUntil(given, now time.Time) (time.Duration, error) {
	// If given is already past now, return difference immediately.
	if given.After(now) {
		return given.Sub(now), nil
	}

	// If given is in the past and not a time-only value,
	// it's an error - we can't calculate "next occurrence" for a specific past date.
	// Time-only values have year 0 or 1970 (depending on the parse format used).
	if given.Year() > 1970 {
		return 0, fmt.Errorf("given time %s is in the past and not a time-only value", given.Format(time.RFC3339))
	}

	// Otherwise, assume that given's date is 1970-01-01 and backdate now.
	nowBackdated, err := time.Parse("15:04:05", now.Format("15:04:05"))
	if err != nil {
		return 0, fmt.Errorf("failed to parse backdated now: %w", err)
	}

	// Add 1 day if today's clock is already past the given time.
	if nowBackdated.After(given) {
		given = given.AddDate(0, 0, 1)
	}

	return given.Sub(nowBackdated), nil
}
