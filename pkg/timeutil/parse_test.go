package timeutil

import (
	"testing"
)

func intPtr(i int) *int {
	return &i
}

func TestParseHumanTime(t *testing.T) {

	tests := []struct {
		name    string
		input   string
		wantErr bool
		wanted  wantedTime
	}{
		// 24-hour time formats
		{
			name:   "24-hour time HH:MM single digits",
			input:  "9:5",
			wanted: newWantedTime().Hour(9).Minute(5),
		},
		{
			name:   "24-hour time HH:MM double digits",
			input:  "15:30",
			wanted: newWantedTime().Hour(15).Minute(30),
		},
		{
			name:   "24-hour time HH:MM:SS single digits",
			input:  "9:5:7",
			wanted: newWantedTime().Hour(9).Minute(5).Second(7),
		},
		{
			name:   "24-hour time HH:MM:SS double digits",
			input:  "15:30:45",
			wanted: newWantedTime().Hour(15).Minute(30).Second(45),
		},
		{
			name:   "24-hour time hour only single digit",
			input:  "9",
			wanted: newWantedTime().Hour(9),
		},
		{
			name:   "24-hour time hour only double digit",
			input:  "15",
			wanted: newWantedTime().Hour(15),
		},

		// 12-hour time formats with PM/AM
		{
			name:   "12-hour time with PM uppercase",
			input:  "3:30PM",
			wanted: newWantedTime().Hour(15).Minute(30),
		},
		{
			name:   "12-hour time with pm lowercase",
			input:  "3:30pm",
			wanted: newWantedTime().Hour(15).Minute(30),
		},
		{
			name:   "12-hour time hour only PM uppercase",
			input:  "3PM",
			wanted: newWantedTime().Hour(15),
		},
		{
			name:   "12-hour time hour only pm lowercase",
			input:  "11pm",
			wanted: newWantedTime().Hour(23),
		},

		// Date formats
		{
			name:   "Date format D.M.YYYY",
			input:  "2.1.2006",
			wanted: newWantedTime().Year(2006).Month(1).Day(2),
		},
		{
			name:   "Date format D.M.YY",
			input:  "15.3.24",
			wanted: newWantedTime().Year(2024).Month(3).Day(15),
		},
		{
			name:   "Date format D.M (current year)",
			input:  "25.12",
			wanted: newWantedTime().Month(12).Day(25),
		},
		{
			name:   "Date format YYYY-MM-DD",
			input:  "2026-04-13",
			wanted: newWantedTime().Year(2026).Month(4).Day(13),
		},

		// Date+time combinations
		{
			name:   "ISO8601 with time HH:MM:SS",
			input:  "2026-04-13T15:30:45",
			wanted: newWantedTime().Year(2026).Month(4).Day(13).Hour(15).Minute(30).Second(45),
		},
		{
			name:   "Date with space and time HH:MM:SS",
			input:  "2026-04-13 15:30:45",
			wanted: newWantedTime().Year(2026).Month(4).Day(13).Hour(15).Minute(30).Second(45),
		},
		{
			name:   "ISO8601 with time HH:MM",
			input:  "2026-04-13T09:05",
			wanted: newWantedTime().Year(2026).Month(4).Day(13).Hour(9).Minute(5),
		},
		{
			name:   "Date with space and time HH:MM",
			input:  "2026-04-13 09:05",
			wanted: newWantedTime().Year(2026).Month(4).Day(13).Hour(9).Minute(5),
		},

		// RFC3339 formats
		{
			name:   "RFC3339 format",
			input:  "2026-04-13T15:30:45Z",
			wanted: newWantedTime().Year(2026).Month(4).Day(13).Hour(15).Minute(30).Second(45),
		},

		// Error cases
		{
			name:    "Invalid format",
			input:   "not a time",
			wantErr: true,
		},
		{
			name:    "Invalid date",
			input:   "99/99/9999",
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseHumanTime(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if tt.wanted.year != nil && result.Year() != *tt.wanted.year {
				t.Errorf("Year: expected %d, got %d", *tt.wanted.year, result.Year())
			}
			if tt.wanted.month != nil && int(result.Month()) != *tt.wanted.month {
				t.Errorf("Month: expected %d, got %d", *tt.wanted.month, result.Month())
			}
			if tt.wanted.day != nil && result.Day() != *tt.wanted.day {
				t.Errorf("Day: expected %d, got %d", *tt.wanted.day, result.Day())
			}
			if tt.wanted.hour != nil && result.Hour() != *tt.wanted.hour {
				t.Errorf("Hour: expected %d, got %d", *tt.wanted.hour, result.Hour())
			}
			if tt.wanted.minute != nil && result.Minute() != *tt.wanted.minute {
				t.Errorf("Minute: expected %d, got %d", *tt.wanted.minute, result.Minute())
			}
			if tt.wanted.second != nil && result.Second() != *tt.wanted.second {
				t.Errorf("Second: expected %d, got %d", *tt.wanted.second, result.Second())
			}
		})
	}
}

type wantedTime struct {
	year, month, day     *int
	hour, minute, second *int
}

func newWantedTime() wantedTime {
	return wantedTime{}
}

func (wt wantedTime) Year(year int) wantedTime {
	wt.year = &year
	return wt
}

func (wt wantedTime) Month(month int) wantedTime {
	wt.month = &month
	return wt
}

func (wt wantedTime) Day(day int) wantedTime {
	wt.day = &day
	return wt
}

func (wt wantedTime) Hour(hour int) wantedTime {
	wt.hour = &hour
	return wt
}

func (wt wantedTime) Minute(minute int) wantedTime {
	wt.minute = &minute
	return wt
}

func (wt wantedTime) Second(second int) wantedTime {
	wt.second = &second
	return wt
}
