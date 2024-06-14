package utils

import "time"

// FormatDateString formats datetime string to date string.
// Example: 2021-01-01T00:00:00Z -> 2021-01-01
func FormatDateString(datetimeStr string) string {
	return datetimeStr[:10]
}

// FormatTimeFromTime formats datetime time.Time to date string.
// Example: 2021-01-01T10:06:23Z -> 10:06
func FormatTimeFromTime(datetime time.Time) string {
	return datetime.Format("15:04")
}

// FormatDateTimeFromTime formats datetime time.Time to datetime string.
// Example: 2021-01-01T10:06:23.9999999Z -> 2021-01-01T10:06:23
func FormatDateTimeFromTime(datetime time.Time) string {
	return datetime.Format("2006-01-02T15:04:05")
}

// FormatDateTimeString formats datetime string to datetime string.
// Example: 2021-01-01T10:06:23.9999999Z -> 2021-01-01T10:06:23
func FormatDateTimeString(datetimeStr string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, datetimeStr)
	if err != nil {
		return time.Time{}
	}

	return t
}
