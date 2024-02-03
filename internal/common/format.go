package utils

import "time"

// FormatDate formats datetime string to date string.
// Example: 2021-01-01T00:00:00Z -> 2021-01-01
func FormatDate(datetimeStr string) string {
	return datetimeStr[:10]
}

// FormatTime formats datetime time.Time to date string.
// Example: 2021-01-01T10:06:23Z -> 10:06
func FormatTime(datetime time.Time) string {
	return datetime.Format("15:04")
}
