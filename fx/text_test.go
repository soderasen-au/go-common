package fx

import (
	"testing"

	"github.com/soderasen-au/go-common/util"
)

func TestIncludeExp_Calc(t *testing.T) {
	tests := []struct {
		name    string
		lh      *Value
		rh      *Value
		want    bool
		wantErr bool
	}{
		{
			name:    "contains substring",
			lh:      Text("hello world"),
			rh:      Text("world"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "does not contain substring",
			lh:      Text("hello world"),
			rh:      Text("xyz"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "contains at beginning",
			lh:      Text("hello world"),
			rh:      Text("hello"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "contains at end",
			lh:      Text("hello world"),
			rh:      Text("world"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "empty substring",
			lh:      Text("hello"),
			rh:      Text(""),
			want:    true,
			wantErr: false,
		},
		{
			name:    "both empty",
			lh:      Text(""),
			rh:      Text(""),
			want:    true,
			wantErr: false,
		},
		{
			name:    "case sensitive",
			lh:      Text("Hello World"),
			rh:      Text("world"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "exact match",
			lh:      Text("test"),
			rh:      Text("test"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "left has error",
			lh:      Error(util.MsgError("test", "error")),
			rh:      Text("test"),
			want:    false,
			wantErr: true,
		},
		{
			name:    "right has error",
			lh:      Text("test"),
			rh:      Error(util.MsgError("test", "error")),
			want:    false,
			wantErr: true,
		},
		{
			name:    "left is numeric",
			lh:      Number(123),
			rh:      Text("test"),
			want:    false,
			wantErr: true,
		},
		{
			name:    "right is numeric",
			lh:      Text("test"),
			rh:      Number(123),
			want:    false,
			wantErr: true,
		},
		{
			name:    "left text is nil",
			lh:      &Value{IsNumeric: false, Text: nil},
			rh:      Text("test"),
			want:    false,
			wantErr: true,
		},
		{
			name:    "right text is nil",
			lh:      Text("test"),
			rh:      &Value{IsNumeric: false, Text: nil},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewIncludeExp(tt.lh, tt.rh)
			result := exp.Calc()

			if (result.HasError()) != tt.wantErr {
				t.Errorf("IncludeExp.Calc() error = %v, wantErr %v", result.HasError(), tt.wantErr)
				return
			}

			if !tt.wantErr {
				if tt.want && !result.True() {
					t.Errorf("IncludeExp.Calc() = false, want true")
				}
				if !tt.want && !result.False() {
					t.Errorf("IncludeExp.Calc() = true, want false")
				}
			}
		})
	}
}

func TestNewIncludeExp(t *testing.T) {
	lh := Text("hello")
	rh := Text("world")
	exp := NewIncludeExp(lh, rh)

	if exp == nil {
		t.Fatal("NewIncludeExp() returned nil")
	}

	// Verify it's an IncludeExp
	includeExp, ok := exp.(*IncludeExp)
	if !ok {
		t.Error("NewIncludeExp() did not return *IncludeExp")
	}

	if includeExp.Lh != lh {
		t.Error("NewIncludeExp() Lh not set correctly")
	}

	if includeExp.Rh != rh {
		t.Error("NewIncludeExp() Rh not set correctly")
	}
}
