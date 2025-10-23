package util

import (
	"encoding/json"
	"testing"
)

func TestMaybeAssignStr(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		defStr   string
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			defStr:   "default",
			expected: "default",
		},
		{
			name:     "empty string pointer",
			input:    Ptr(""),
			defStr:   "default",
			expected: "default",
		},
		{
			name:     "non-empty string pointer",
			input:    Ptr("existing"),
			defStr:   "default",
			expected: "existing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.input
			MaybeAssignStr(&s, tt.defStr)
			if s == nil || *s != tt.expected {
				t.Errorf("MaybeAssignStr() = %v, want %v", StrMaybeNil(s), tt.expected)
			}
		})
	}
}

func TestStrMaybeNil(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		expected string
	}{
		{
			name:     "nil pointer",
			input:    nil,
			expected: "",
		},
		{
			name:     "empty string",
			input:    Ptr(""),
			expected: "",
		},
		{
			name:     "non-empty string",
			input:    Ptr("hello"),
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrMaybeNil(tt.input)
			if result != tt.expected {
				t.Errorf("StrMaybeNil() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestStrMaybeDefault(t *testing.T) {
	tests := []struct {
		name     string
		input    *string
		defStr   string
		expected string
	}{
		{
			name:     "nil pointer uses default",
			input:    nil,
			defStr:   "default",
			expected: "default",
		},
		{
			name:     "empty string returns empty",
			input:    Ptr(""),
			defStr:   "default",
			expected: "",
		},
		{
			name:     "non-empty string returns value",
			input:    Ptr("value"),
			defStr:   "default",
			expected: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StrMaybeDefault(tt.input, tt.defStr)
			if result != tt.expected {
				t.Errorf("StrMaybeDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIntMaybeDefault(t *testing.T) {
	tests := []struct {
		name     string
		input    *int
		def      int
		expected int
	}{
		{
			name:     "nil pointer uses default",
			input:    nil,
			def:      42,
			expected: 42,
		},
		{
			name:     "zero value returns zero",
			input:    Ptr(0),
			def:      42,
			expected: 0,
		},
		{
			name:     "non-zero value returns value",
			input:    Ptr(100),
			def:      42,
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IntMaybeDefault(tt.input, tt.def)
			if result != tt.expected {
				t.Errorf("IntMaybeDefault() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestBoolMaybeDeNil(t *testing.T) {
	tests := []struct {
		name     string
		input    *bool
		expected bool
	}{
		{
			name:     "nil pointer returns false",
			input:    nil,
			expected: false,
		},
		{
			name:     "false value returns false",
			input:    Ptr(false),
			expected: false,
		},
		{
			name:     "true value returns true",
			input:    Ptr(true),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BoolMaybeDeNil(tt.input)
			if result != tt.expected {
				t.Errorf("BoolMaybeDeNil() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMaybeNil(t *testing.T) {
	t.Run("nil string pointer", func(t *testing.T) {
		var s *string
		result := MaybeNil(s)
		if result != "" {
			t.Errorf("MaybeNil() = %v, want empty string", result)
		}
	})

	t.Run("non-nil string pointer", func(t *testing.T) {
		s := Ptr("test")
		result := MaybeNil(s)
		if result != "test" {
			t.Errorf("MaybeNil() = %v, want test", result)
		}
	})

	t.Run("nil int pointer", func(t *testing.T) {
		var i *int
		result := MaybeNil(i)
		if result != 0 {
			t.Errorf("MaybeNil() = %v, want 0", result)
		}
	})

	t.Run("non-nil int pointer", func(t *testing.T) {
		i := Ptr(42)
		result := MaybeNil(i)
		if result != 42 {
			t.Errorf("MaybeNil() = %v, want 42", result)
		}
	})
}

func TestMaybeDefault(t *testing.T) {
	t.Run("nil string pointer uses default", func(t *testing.T) {
		var s *string
		result := MaybeDefault(s, "default")
		if result != "default" {
			t.Errorf("MaybeDefault() = %v, want default", result)
		}
	})

	t.Run("non-nil string pointer returns value", func(t *testing.T) {
		s := Ptr("value")
		result := MaybeDefault(s, "default")
		if result != "value" {
			t.Errorf("MaybeDefault() = %v, want value", result)
		}
	})

	t.Run("nil int pointer uses default", func(t *testing.T) {
		var i *int
		result := MaybeDefault(i, 42)
		if result != 42 {
			t.Errorf("MaybeDefault() = %v, want 42", result)
		}
	})

	t.Run("non-nil int pointer returns value", func(t *testing.T) {
		i := Ptr(100)
		result := MaybeDefault(i, 42)
		if result != 100 {
			t.Errorf("MaybeDefault() = %v, want 100", result)
		}
	})
}

func TestMax(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int
		expected int
	}{
		{"x greater", 10, 5, 10},
		{"y greater", 5, 10, 10},
		{"equal", 7, 7, 7},
		{"negative numbers", -5, -10, -5},
		{"mixed signs", -5, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Max(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("Max(%d, %d) = %d, want %d", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int
		expected int
	}{
		{"x smaller", 5, 10, 5},
		{"y smaller", 10, 5, 5},
		{"equal", 7, 7, 7},
		{"negative numbers", -5, -10, -10},
		{"mixed signs", -5, 5, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Min(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("Min(%d, %d) = %d, want %d", tt.x, tt.y, result, tt.expected)
			}
		})
	}
}

func TestGetAPIUrl(t *testing.T) {
	tests := []struct {
		name     string
		baseUrl  string
		endpoint string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple endpoint",
			baseUrl:  "https://example.com",
			endpoint: "/users",
			expected: "https://example.com/api/v1/users",
			wantErr:  false,
		},
		{
			name:     "endpoint with /api/v1 prefix",
			baseUrl:  "https://example.com",
			endpoint: "/api/v1/users",
			expected: "https://example.com/api/v1/users",
			wantErr:  false,
		},
		{
			name:     "base url with path",
			baseUrl:  "https://example.com/base",
			endpoint: "/users",
			expected: "https://example.com/api/v1/users",
			wantErr:  false,
		},
		{
			name:     "invalid base url",
			baseUrl:  "://invalid",
			endpoint: "/users",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetAPIUrl(tt.baseUrl, tt.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAPIUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("GetAPIUrl() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetUrl(t *testing.T) {
	tests := []struct {
		name     string
		baseUrl  string
		endpoint string
		expected string
		wantErr  bool
	}{
		{
			name:     "simple endpoint",
			baseUrl:  "https://example.com",
			endpoint: "/users",
			expected: "https://example.com/users",
			wantErr:  false,
		},
		{
			name:     "base url with path",
			baseUrl:  "https://example.com/api",
			endpoint: "/users",
			expected: "https://example.com/api/users",
			wantErr:  false,
		},
		{
			name:     "nested endpoint",
			baseUrl:  "https://example.com",
			endpoint: "/api/v2/users",
			expected: "https://example.com/api/v2/users",
			wantErr:  false,
		},
		{
			name:     "invalid base url",
			baseUrl:  "://invalid",
			endpoint: "/users",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetUrl(tt.baseUrl, tt.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("GetUrl() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewString(t *testing.T) {
	str := "test"
	result := NewString(str)
	if result == nil {
		t.Error("NewString() returned nil")
		return
	}
	if *result != str {
		t.Errorf("NewString() = %v, want %v", *result, str)
	}
}

func TestPtr(t *testing.T) {
	t.Run("string pointer", func(t *testing.T) {
		str := "test"
		result := Ptr(str)
		if result == nil || *result != str {
			t.Errorf("Ptr() = %v, want %v", result, &str)
		}
	})

	t.Run("int pointer", func(t *testing.T) {
		num := 42
		result := Ptr(num)
		if result == nil || *result != num {
			t.Errorf("Ptr() = %v, want %v", result, &num)
		}
	})

	t.Run("bool pointer", func(t *testing.T) {
		b := true
		result := Ptr(b)
		if result == nil || *result != b {
			t.Errorf("Ptr() = %v, want %v", result, &b)
		}
	})
}

func TestInt32to64(t *testing.T) {
	tests := []struct {
		name     string
		input    int32
		expected int64
	}{
		{"positive number", 42, 42},
		{"negative number", -42, -42},
		{"zero", 0, 0},
		{"max int32", 2147483647, 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input
			result := Int32to64(&input)
			if result == nil {
				t.Error("Int32to64() returned nil")
				return
			}
			if *result != tt.expected {
				t.Errorf("Int32to64() = %v, want %v", *result, tt.expected)
			}
		})
	}
}

func TestStupidDeepCopy(t *testing.T) {
	type TestStruct struct {
		Name  string
		Value int
		Tags  []string
	}

	tests := []struct {
		name    string
		input   *TestStruct
		wantErr bool
	}{
		{
			name: "simple struct",
			input: &TestStruct{
				Name:  "test",
				Value: 42,
				Tags:  []string{"tag1", "tag2"},
			},
			wantErr: false,
		},
		{
			name: "empty struct",
			input: &TestStruct{
				Name:  "",
				Value: 0,
				Tags:  nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := StupidDeepCopy(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StupidDeepCopy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && result != nil {
				// Compare using JSON to avoid pointer comparison issues
				inputJSON, _ := json.Marshal(tt.input)
				resultJSON, _ := json.Marshal(result)
				if string(inputJSON) != string(resultJSON) {
					t.Errorf("StupidDeepCopy() = %v, want %v", string(resultJSON), string(inputJSON))
				}
				// Ensure it's a different pointer
				if result == tt.input {
					t.Error("StupidDeepCopy() returned same pointer, not a copy")
				}
			}
		})
	}
}
