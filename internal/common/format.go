package utils

import "time"

// FormatDate formats datetime string to date string.
// Example: 2021-01-01T00:00:00Z -> 2021-01-01
func FormatDate(time time.Time) string {
	return time.Format("2006-01-02")
}

// DateStrToDate converts date string to time.Time.
func DateStrToDate(date string) time.Time {
	t, _ := time.Parse("2006-01-02", date)
	return t
}

// FormatTime formats datetime time.Time to date string.
// Example: 2021-01-01T10:06:23Z -> 10:06
func FormatTime(datetime time.Time) string {
	return datetime.Format("15:04")
}

// FormatDatetime formats datetime time.Time to datetime string.
// Example: 2021-01-01T10:06:23Z -> 2021-01-01T10:06:23Z
func FormatDatetime(datetime time.Time) string {
	return datetime.Format("2006-01-02T15:04:05Z")
}
