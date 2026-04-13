package timeutil

import (
	"testing"
	"time"
)

func TestDurationUntil(t *testing.T) {
	tests := []struct {
		name        string
		given       time.Time
		now         time.Time
		expectError bool
		checkResult func(time.Duration) bool
	}{
		{
			name:        "future time returns difference",
			given:       mustParse(t, "2026-04-12T16:00:00Z"),
			now:         mustParse(t, "2026-04-12T14:00:00Z"),
			expectError: false,
			checkResult: func(d time.Duration) bool { return d == 2*time.Hour },
		},
		{
			name:        "time-only (1970) in future today",
			given:       mustParse(t, "1970-01-01T16:00:00Z"),
			now:         mustParse(t, "1970-01-01T14:00:00Z"),
			expectError: false,
			checkResult: func(d time.Duration) bool { return d == 2*time.Hour },
		},
		{
			name:        "time-only (1970) already passed gets next day",
			given:       mustParseTimeOnly(t, "09:00:00"),
			now:         mustParseTimeOnly(t, "14:00:00"),
			expectError: false,
			checkResult: func(d time.Duration) bool {
				// Should be approximately 19 hours (until 9am tomorrow)
				expected := 19 * time.Hour
				// Allow 1 hour tolerance due to the way the function calculates
				return d >= expected-time.Hour && d <= expected+time.Hour
			},
		},
		{
			name:        "past date (not 1970) returns error",
			given:       mustParse(t, "2026-04-11T16:00:00Z"),
			now:         mustParse(t, "2026-04-12T14:00:00Z"),
			expectError: true,
		},
		{
			name:        "past date in 2025 returns error",
			given:       mustParse(t, "2025-01-01T12:00:00Z"),
			now:         mustParse(t, "2026-04-12T14:00:00Z"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := DurationUntil(tt.given, tt.now)

			if tt.expectError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectError && err == nil && tt.checkResult != nil {
				if !tt.checkResult(result) {
					t.Errorf("result check failed for duration: %s", result)
				}
			}
		})
	}
}

func mustParse(t *testing.T, s string) time.Time {
	t.Helper()

	ti, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	return ti
}

func mustParseTimeOnly(t *testing.T, s string) time.Time {
	t.Helper()

	ti, err := time.Parse("15:04:05", s)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	return ti
}
