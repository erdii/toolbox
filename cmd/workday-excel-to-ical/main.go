package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	ics "github.com/arran4/golang-ical"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"github.com/xuri/excelize/v2"
)

var (
	organiser string
	domain    string
	summary   string
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	cmd.Flags().StringVarP(&organiser, "organiser", "o", "", "organiser of created events")
	cmd.Flags().StringVarP(&domain, "domain", "d", "", "domain of created events")
	cmd.Flags().StringVarP(&summary, "summary", "s", "", "summary of created events")

	must(cmd.MarkFlagRequired("organiser"))
	must(cmd.MarkFlagRequired("domain"))
	must(cmd.MarkFlagRequired("summary"))
}

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}

var cmd = &cobra.Command{
	Use:   "workday-excel-to-ical <INPUT SHEET.xlsx>",
	Short: "get workday pto sheets into your calendar",
	Long:  `cheaply convert workday pto approval sheets into an ical events which can be imported into your company calendar`,
	RunE: func(_ *cobra.Command, args []string) error {
		return run(args[0])
	},
}

var ErrWrongDateSyntax = errors.New("wrong date syntax")

func run(inputFilePath string) error {
	var err error

	f, err := excelize.OpenFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	defer func() {
		cErr := f.Close()
		if cErr != nil && err != nil {
			err = fmt.Errorf("while closing sheet: %w, while: %w", cErr, err)
		}
	}()

	sheets := f.GetSheetList()
	if slen := len(sheets); slen < 1 {
		return fmt.Errorf("excel file must contain at least 1 sheet. got: %d", slen)
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return fmt.Errorf("getting rows: %w", err)
	}

	dates, err := datesFromRows(logr.Discard(), rows, func(rawDate string) (time.Time, error) {
		// switch order from mm-dd-yy to yy-mm-dd before parsing
		mdy := strings.Split(rawDate, "-")
		if len(mdy) != 3 {
			return time.Time{}, fmt.Errorf("parsing date: %w: %s", err, rawDate)
		}
		dmy := []string{"20" + mdy[2], mdy[0], mdy[1]}
		return time.Parse(time.DateOnly, strings.Join(dmy, "-"))
	})

	if err != nil {
		return fmt.Errorf("extracting dates from rows: %w", err)
	}

	cal, err := icalFromDates(dates)
	if err != nil {
		return fmt.Errorf("converting dates to ical format: %w", err)
	}

	fmt.Fprintf(os.Stderr, "Loaded %d days of PTO from Excel sheet:\n", len(dates))
	for _, date := range dates {
		fmt.Fprintf(os.Stderr, "- %s\n", date.Format(time.DateOnly))
	}

	fmt.Fprintln(os.Stdout, cal)
	return nil
}

func datesFromRows[T any](log logr.Logger, rows [][]string, convert func(string) (T, error)) ([]T, error) {
	dates := []T{}

	// Extract date from first cell in each row, skipping 2 header rows.
	for _, row := range rows[2:] {
		duration := fmt.Sprintf("%s %s", row[3], row[4])
		if duration != "1 Days" {
			return nil, fmt.Errorf("row [%s] must describe 1 day of pto", strings.Join(row, " "))
		}

		rawDate := row[0]
		date, err := convert(rawDate)
		if err != nil {
			return nil, fmt.Errorf("converting date: %w", err)
		}
		dates = append(dates, date)
	}

	return dates, nil
}

func icalFromDates(dates []time.Time) (string, error) {
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodPublish)

	for _, date := range dates {
		event := cal.AddEvent(fmt.Sprintf("%x@%s", md5.Sum([]byte(date.Format(time.DateOnly))), domain))

		event.SetCreatedTime(time.Now())
		event.SetDtStampTime(time.Now())
		event.SetModifiedAt(time.Now())
		event.SetAllDayStartAt(date)
		event.SetSummary(summary)
		event.SetOrganizer(organiser)
	}

	return cal.Serialize(), nil
}
