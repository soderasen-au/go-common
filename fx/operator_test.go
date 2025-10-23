package fx

import (
	"testing"
)

func TestParseOperator(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Operator
		wantErr bool
	}{
		{
			name:    "less than",
			input:   "<",
			want:    LessOp,
			wantErr: false,
		},
		{
			name:    "greater than",
			input:   ">",
			want:    GreaterOp,
			wantErr: false,
		},
		{
			name:    "equal",
			input:   "==",
			want:    EqualOp,
			wantErr: false,
		},
		{
			name:    "less or equal",
			input:   "<=",
			want:    LessEqOp,
			wantErr: false,
		},
		{
			name:    "greater or equal",
			input:   ">=",
			want:    GreaterEqOp,
			wantErr: false,
		},
		{
			name:    "not equal",
			input:   "!=",
			want:    NotEqOp,
			wantErr: false,
		},
		{
			name:    "include",
			input:   "~=",
			want:    IncludeOp,
			wantErr: false,
		},
		{
			name:    "unknown operator",
			input:   "???",
			want:    UnkownOp,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    UnkownOp,
			wantErr: true,
		},
		{
			name:    "invalid operator",
			input:   "&&",
			want:    UnkownOp,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOperator(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseOperator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseOperator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOperatorConstants(t *testing.T) {
	tests := []struct {
		name string
		op   Operator
		want string
	}{
		{"UnknownOp", UnkownOp, ""},
		{"LessOp", LessOp, "<"},
		{"GreaterOp", GreaterOp, ">"},
		{"EqualOp", EqualOp, "=="},
		{"LessEqOp", LessEqOp, "<="},
		{"GreaterEqOp", GreaterEqOp, ">="},
		{"NotEqOp", NotEqOp, "!="},
		{"IncludeOp", IncludeOp, "~="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.op) != tt.want {
				t.Errorf("Operator %s = %v, want %v", tt.name, tt.op, tt.want)
			}
		})
	}
}
