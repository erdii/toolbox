/*
Copyright (C) 2026 erdii <me@erdii.engineering>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	input := strings.Join(os.Args[1:], " ")

	given, err := parseTime(input)
	if err != nil {
		log.Fatalf("failed to parse given time: %s", err)
	}

	now := time.Now()

	diff, err := calcSleepDuration(given, now)
	if err != nil {
		log.Fatalf("failed to calculate sleep duration: %s", err)
	}

	fmt.Printf("given:\t%s\n", given.Format("15:04:05"))
	fmt.Printf("now:\t%s\n", now.Format("15:04:05"))
	fmt.Printf("diff:\t%s\n", diff.Round(time.Second))

	time.Sleep(diff)
}

func parseTime(input string) (time.Time, error) {
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

	return time.Time{}, fmt.Errorf("could not parse time: %w", input)
}

func calcSleepDuration(given, now time.Time) (time.Duration, error) {
	// If given is already past now, return difference immediately.
	if given.After(now) {
		return given.Sub(now), nil
	}

	// Otherwise, assume that given's date is 1970-01-01 and backdate now.
	nowBackdated, err := time.Parse("15:04:05", now.Format("15:04:05"))
	if err != nil {
		return 0, fmt.Errorf("failed to parse backdated now: %err", err)
	}

	// Add 1 day if today's clock is already past the given time.
	if nowBackdated.After(given) {
		given = given.AddDate(0, 0, 1)
	}

	return given.Sub(nowBackdated), nil
}
