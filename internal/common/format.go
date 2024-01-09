package utils

// FormatDate formats datetime string to date string.
// Example: 2021-01-01T00:00:00Z -> 2021-01-01
func FormatDate(datetimeStr string) string {
	return datetimeStr[:10]
}
