package util

import (
	"reflect"
	"time"
)

var (
	defaultDateFormatInString  = "20060102"
	dateFormatInStringWithDash = "2006-01-02"
)

func FormatDateFromStringWithoutDash(dateString string) time.Time {
	return formatDateFromString(dateString, defaultDateFormatInString)
}

func FormatDateFromStringWithDash(dateString string) time.Time {
	return formatDateFromString(dateString, dateFormatInStringWithDash)
}

func formatDateFromString(dateString, format string) time.Time {
	date, err := time.Parse(format, dateString)
	if err != nil {
		Logger.Errorln(err)
	}
	return date
}

func FormatDateToStringWithoutDash(date time.Time) string {
	return formatDateToString(date, defaultDateFormatInString)
}

func FormatDateToStringWithDash(date time.Time) string {
	return formatDateToString(date, dateFormatInStringWithDash)
}

func formatDateToString(date time.Time, format string) string {
	return date.Format(format)
}

func IsDateTimeEmpty(dateTime time.Time) bool {
	return reflect.DeepEqual(dateTime, time.Time{})
}

// ParseDate parses a date string in either YYYYMMDD or YYYY-MM-DD format and returns a time.Time
// Returns an error if the date string is invalid
func ParseDate(dateStr string) (time.Time, error) {
	// Try parsing without dash first (YYYYMMDD)
	date, err := time.Parse(defaultDateFormatInString, dateStr)
	if err == nil {
		return date, nil
	}

	// Try parsing with dash (YYYY-MM-DD)
	date, err = time.Parse(dateFormatInStringWithDash, dateStr)
	if err == nil {
		return date, nil
	}

	// Both formats failed
	return time.Time{}, err
}
