package crypto

import (
	"testing"
)

func TestSHA2656Hex(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		wantErr  bool
	}{
		{
			name:     "empty input",
			input:    []byte{},
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr:  false,
		},
		{
			name:     "simple text",
			input:    []byte("hello"),
			expected: "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824",
			wantErr:  false,
		},
		{
			name:     "another text",
			input:    []byte("test"),
			expected: "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08",
			wantErr:  false,
		},
		{
			name:     "longer text",
			input:    []byte("The quick brown fox jumps over the lazy dog"),
			expected: "d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e592",
			wantErr:  false,
		},
		{
			name:     "binary data",
			input:    []byte{0x00, 0x01, 0x02, 0x03, 0xff},
			expected: "ff5d8507b6a72bee2debce2c0054798deaccdc5d8a1b945b6280ce8aa9cba52e",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := SHA2656Hex(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SHA2656Hex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("SHA2656Hex() = %v, want %v", result, tt.expected)
			}
		})
	}
}
