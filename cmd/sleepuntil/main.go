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
	"os"
	"strings"
	"time"

	"github.com/erdii/toolbox/pkg/stderr"
	"github.com/erdii/toolbox/pkg/timeutil"
)

func main() {
	input := strings.Join(os.Args[1:], " ")

	given, err := timeutil.ParseHumanTime(input)
	if err != nil {
		stderr.Fatalf("failed to parse given time: %s", err)
	}

	now := time.Now()

	diff, err := timeutil.DurationUntil(given, now)
	if err != nil {
		stderr.Fatalf("failed to calculate sleep duration: %s", err)
	}

	stderr.Printf("given:\t%s\n", given.Format("15:04:05"))
	stderr.Printf("now:\t%s\n", now.Format("15:04:05"))
	stderr.Printf("diff:\t%s\n", diff.Round(time.Second))

	time.Sleep(diff)
}
