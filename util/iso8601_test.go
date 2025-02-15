package util

import "testing"

func TestFromISO8601(t *testing.T) {
	tests := []struct {
		isoFormat string
		expected  string
	}{
		// Basic Date Formats
		{"YYYY-MM-DD", "2006-01-02"},
		{"YYYY-MMM-DD", "2006-Jan-02"},
		{"YYYY-MMMM-DD", "2006-January-02"},
		{"YY-MM-DD", "06-01-02"},

		// Date-Time Formats
		{"YYYY-MM-DD HH:mm:ss", "2006-01-02 15:04:05"},
		{"YYYY-MM-DD hh:mm:ss", "2006-01-02 03:04:05"},
		{"YYYY-MM-DD HH:mm", "2006-01-02 15:04"},
		{"YYYY-MM-DD HH", "2006-01-02 15"},

		// Date-Time with Timezone
		{"YYYY-MM-DDTHH:mm:ssZ", "2006-01-02T15:04:05Z07:00"},
		{"YYYY-MM-DDTHH:mm:ss±hh:mm", "2006-01-02T15:04:05-07:00"},

		// Time-Only Formats
		{"HH:mm:ss", "15:04:05"},
		{"hh:mm:ss", "03:04:05"},
		{"HH:mm", "15:04"},
		{"hh:mm", "03:04"},

		// Milliseconds & Microseconds
		{"YYYY-MM-DD HH:mm:ss.SSS", "2006-01-02 15:04:05.000"},
		{"YYYY-MM-DD HH:mm:ss.SSSSSS", "2006-01-02 15:04:05.000000"},

		// Edge Cases
		{"YYYY", "2006"},             // Only year
		{"YYYY-MM", "2006-01"},       // Year and month
		{"MM-DD", "01-02"},           // Month and day
		{"DD-MM-YYYY", "02-01-2006"}, // Different order
	}

	for _, test := range tests {
		got := FromISO8601(test.isoFormat)
		if got != test.expected {
			t.Errorf("FromISO8601(%q) = %q; want %q", test.isoFormat, got, test.expected)
		}
	}
}
