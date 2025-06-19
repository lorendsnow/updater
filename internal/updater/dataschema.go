package updater

import (
	"log/slog"
	"strconv"
	"time"
)

/*
 *==================================================================================================
 * Date and Time Conversion Formats
 *==================================================================================================
 */

const DATE_TIME_FORMAT = "01/02/2006 1504"
const DATE_ONLY_FORMAT = "01/02/2006"

/*
 *==================================================================================================
 * Record Struct
 *==================================================================================================
 */

// Record represents a single crime record from the City of Porland's data.
type Record struct {
	Address         string
	CaseNumber      string
	CrimeAgainst    string
	Neighborhood    string
	OccurDateTime   time.Time
	OffenseCategory string
	OffenseType     string
	OpenDataLat     *float64 // we use pointers for primitives to allow for nil values
	OpenDataLon     *float64
	OpenDataX       *float64
	OpenDataY       *float64
	ReportDate      time.Time
	OffenseCount    *int
}

/*
 *==================================================================================================
 * Public Functions
 *==================================================================================================
 */

// NewRecord takes a row of strings from a CSV file and marshals the data into
// a Record.
func NewRecord(row []string, logger *slog.Logger) Record {
	if len(row) != 14 {
		logger.Error("bad data format - expected 14 columns", "row length", len(row))
		return Record{}
	}
	return Record{
		Address:         row[0],
		CaseNumber:      row[1],
		CrimeAgainst:    row[2],
		Neighborhood:    row[3],
		OccurDateTime:   parseDateTime(row[4], row[5], logger),
		OffenseCategory: row[6],
		OffenseType:     row[7],
		OpenDataLat:     parseFloat(row[8]),
		OpenDataLon:     parseFloat(row[9]),
		OpenDataX:       parseFloat(row[10]),
		OpenDataY:       parseFloat(row[11]),
		ReportDate:      parseDate(row[12], logger),
		OffenseCount:    parseInt(row[13]),
	}
}

/*
 *==================================================================================================
 * Private Functions
 *==================================================================================================
 */

// parseDate takes a date string in the format "MM/DD/YYYY" and returns a
// time.Time with UTC location. If the date string is empty or there's an error
// while parsing the string, it returns a default value of "01/01/1900".
func parseDate(date string, logger *slog.Logger) time.Time {
	formattedDate, err := time.Parse(DATE_ONLY_FORMAT, date)

	if err != nil {
		logger.Error(
			"Failed to parse date; using default value '01/01/1900'",
			"date",
			date,
			"error",
			err,
		)
		formattedDate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	return formattedDate
}

// parseDateTime takes a date string in the format "MM/DD/YYYY" and a time
// string in the format "HHMM" and returns a time.Time with UTC location. If
// the date or time string is empty or there's an error while parsing the
// strings, it returns a default value of "01/01/1900 00:00".
func parseDateTime(date string, timeOnly string, logger *slog.Logger) time.Time {
	timeStr := date + " " + timeOnly

	formattedDate, err := time.Parse(DATE_TIME_FORMAT, timeStr)

	if err != nil {
		logger.Error(
			"Failed to parse date and time; using default value '01/01/1900 00:00'",
			"time",
			timeOnly,
			"error",
			err,
		)
		formattedDate = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	return formattedDate
}

// parseFloat takes a string and returns a float64. If the string is empty or
// there's another error while parsing the string, it returns nil.
func parseFloat(s string) *float64 {
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}

// parseInt takes a string and returns an integer. If the string is empty or
// there's another error while parsing the string, it returns nil.
func parseInt(s string) *int {
	if s == "" {
		return nil
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &i
}
