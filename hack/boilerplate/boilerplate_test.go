/*
Copyright (C) 2021 erdii <me@erdii.engineering>

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
	"bytes"
	"strings"
	"testing"
)

func Test_fileExtension(t *testing.T) {
	tests := []struct {
		name, file, extension string
	}{
		{
			name:      "Simple extension",
			file:      "test.go",
			extension: "go",
		},
		{
			name:      "dot in between",
			file:      "script.test.sh",
			extension: "sh",
		},
		{
			name:      "no extension",
			file:      "script",
			extension: "",
		},
		{
			name:      "dot in folder structure",
			file:      "k8s.io/script",
			extension: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ext := fileExtension(test.file)
			if ext != test.extension {
				t.Errorf("file extension should be %q, is: %q", test.extension, ext)
			}
		})
	}
}

var expected = `Boilerplate header is wrong for:
test/fail.go
test/fail.py`

func Test_run(t *testing.T) {
	var buf bytes.Buffer

	failed, err := run(&buf, "./test")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}
	if !failed {
		t.Errorf("should have failed")
	}
	output := strings.TrimSpace(buf.String())
	if output != expected {
		t.Errorf("unexpected messages printed:\n%s\nshould be:\n%s", output, expected)
	}
}
