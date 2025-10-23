package util

import "testing"

func TestSumSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{
			name:     "empty slice",
			input:    []int{},
			expected: 0,
		},
		{
			name:     "single element",
			input:    []int{5},
			expected: 5,
		},
		{
			name:     "multiple positive numbers",
			input:    []int{1, 2, 3, 4, 5},
			expected: 15,
		},
		{
			name:     "mixed positive and negative",
			input:    []int{10, -5, 3, -2},
			expected: 6,
		},
		{
			name:     "all negative numbers",
			input:    []int{-1, -2, -3},
			expected: -6,
		},
		{
			name:     "with zeros",
			input:    []int{0, 5, 0, 10},
			expected: 15,
		},
		{
			name:     "all zeros",
			input:    []int{0, 0, 0},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SumSlice(tt.input)
			if result != tt.expected {
				t.Errorf("SumSlice(%v) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}
