package stderr

import (
	"fmt"
	"os"
)

// Printf formats according to a format specifier and writes to stderr.
// It returns the number of bytes written and any write error encountered.
func Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

// Println formats using the default formats for its operands and writes to stderr.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func Println(a ...any) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}

// Fatalf is equivalent to [Printf] followed by a call to os.Exit(1).
func Fatalf(format string, a ...any) {
	Printf(format, a...)
	os.Exit(1)
}
