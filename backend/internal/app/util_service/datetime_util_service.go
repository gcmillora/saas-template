package utils

import "time"

func FormatDateTime(
	date time.Time,
) string {
	format := "2006/01/02 15:04"

	return date.Format(format)
}

