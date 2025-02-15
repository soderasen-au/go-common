package util

import "strings"

// Convert ISO 8601 format string to Go's reference layout
func FromISO8601(isoFormat string) string {
	replacer := strings.NewReplacer(
		"YYYY", "2006",
		"YY", "06",
		"MMMM", "January", // Full month name
		"MMM", "Jan", // Abbreviated month name
		"MM", "01",
		"DD", "02",
		"HH", "15",
		"hh", "03",
		"mm", "04",
		"ss", "05",
		".SSSSSS", ".000000", // Microseconds (ensure dot placement)
		".SSS", ".000", // Milliseconds (ensure dot placement)
		"Z", "Z07:00", // Timezone offset
		"±hh:mm", "-07:00",
	)
	return replacer.Replace(isoFormat)
}
