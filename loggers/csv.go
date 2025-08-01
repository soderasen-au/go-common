package loggers

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// CSVWriter wraps a CSV writer for zerolog
type CSVWriter struct {
	writer    *csv.Writer
	delimiter rune
	headers   []string
}

// NewCSVWriter creates a new CSV writer
func NewCSVWriter(w io.Writer, delimiter rune, headers []string) *CSVWriter {
	csvWriter := csv.NewWriter(w)
	csvWriter.Comma = delimiter

	cw := &CSVWriter{
		writer:    csvWriter,
		delimiter: delimiter,
		headers:   headers,
	}

	// Write headers
	cw.writer.Write(headers)
	cw.writer.Flush()

	return cw
}

// Write implements io.Writer interface for zerolog
func (cw *CSVWriter) Write(p []byte) (n int, err error) {
	// Parse the JSON log entry from zerolog
	logStr := string(p)
	logStr = strings.TrimSpace(logStr)

	// Simple JSON parsing for demonstration
	// In production, you might want to use encoding/json
	record := make([]string, len(cw.headers))

	for i, header := range cw.headers {
		value := extractJSONValue(logStr, header)
		record[i] = value
	}

	err = cw.writer.Write(record)
	if err != nil {
		return 0, err
	}

	cw.writer.Flush()
	return len(p), nil
}

// Simple JSON value extractor (basic implementation)
func extractJSONValue(jsonStr, key string) string {
	// Look for "key":"value" pattern
	pattern := fmt.Sprintf(`"%s":"`, key)
	start := strings.Index(jsonStr, pattern)
	if start == -1 {
		// Try without quotes for numbers/booleans
		pattern = fmt.Sprintf(`"%s":`, key)
		start = strings.Index(jsonStr, pattern)
		if start == -1 {
			return ""
		}
		start += len(pattern)

		// Find the end (comma or closing brace)
		end := start
		for end < len(jsonStr) && jsonStr[end] != ',' && jsonStr[end] != '}' {
			end++
		}
		return strings.TrimSpace(jsonStr[start:end])
	}

	start += len(pattern)

	// Find closing quote
	end := start
	for end < len(jsonStr) && jsonStr[end] != '"' {
		if jsonStr[end] == '\\' {
			end++ // Skip escaped character
		}
		end++
	}

	if end < len(jsonStr) {
		return jsonStr[start:end]
	}
	return ""
}
