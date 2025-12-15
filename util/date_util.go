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
