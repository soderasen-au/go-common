package fx

import (
	"testing"
)

func TestNewBinOpExp(t *testing.T) {
	tests := []struct {
		name    string
		op      Operator
		lh      *Value
		rh      *Value
		wantErr bool
	}{
		{
			name:    "LessOp",
			op:      LessOp,
			lh:      Number(5),
			rh:      Number(10),
			wantErr: false,
		},
		{
			name:    "GreaterOp",
			op:      GreaterOp,
			lh:      Number(10),
			rh:      Number(5),
			wantErr: false,
		},
		{
			name:    "EqualOp",
			op:      EqualOp,
			lh:      Number(5),
			rh:      Number(5),
			wantErr: false,
		},
		{
			name:    "unknown operator",
			op:      Operator("???"),
			lh:      Number(5),
			rh:      Number(10),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp, err := NewBinOpExp(tt.op, tt.lh, tt.rh)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBinOpExp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && exp == nil {
				t.Error("NewBinOpExp() returned nil expression")
			}
		})
	}
}

func TestRegisterNewBinOpExpCreator(t *testing.T) {
	// Create a custom operator
	customOp := Operator("CUSTOM")

	// Register a creator
	RegisterNewBinOpExpCreator(customOp, func(lh, rh *Value) Expression {
		return &LessExp{BinOpExp{Op: customOp, Lh: lh, Rh: rh}}
	})

	// Verify it was registered
	exp, err := NewBinOpExp(customOp, Number(1), Number(2))
	if err != nil {
		t.Errorf("NewBinOpExp() after registration failed: %v", err)
	}
	if exp == nil {
		t.Error("NewBinOpExp() returned nil after registration")
	}
}

func TestLessExp_Calc(t *testing.T) {
	tests := []struct {
		name string
		lh   *Value
		rh   *Value
		want bool
	}{
		{
			name: "5 < 10",
			lh:   Number(5),
			rh:   Number(10),
			want: true,
		},
		{
			name: "10 < 5",
			lh:   Number(10),
			rh:   Number(5),
			want: false,
		},
		{
			name: "5 < 5",
			lh:   Number(5),
			rh:   Number(5),
			want: false,
		},
		{
			name: "abc < bcd",
			lh:   Text("abc"),
			rh:   Text("bcd"),
			want: true,
		},
		{
			name: "bcd < abc",
			lh:   Text("bcd"),
			rh:   Text("abc"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewLessExp(tt.lh, tt.rh)
			result := exp.Calc()
			if tt.want && !result.True() {
				t.Errorf("LessExp.Calc() = false, want true")
			}
			if !tt.want && !result.False() {
				t.Errorf("LessExp.Calc() = true, want false")
			}
		})
	}
}

func TestGreatorExp_Calc(t *testing.T) {
	tests := []struct {
		name string
		lh   *Value
		rh   *Value
		want bool
	}{
		{
			name: "10 > 5",
			lh:   Number(10),
			rh:   Number(5),
			want: true,
		},
		{
			name: "5 > 10",
			lh:   Number(5),
			rh:   Number(10),
			want: false,
		},
		{
			name: "5 > 5",
			lh:   Number(5),
			rh:   Number(5),
			want: false,
		},
		{
			name: "bcd > abc",
			lh:   Text("bcd"),
			rh:   Text("abc"),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewGreatorExp(tt.lh, tt.rh)
			result := exp.Calc()
			if tt.want && !result.True() {
				t.Errorf("GreatorExp.Calc() = false, want true")
			}
			if !tt.want && !result.False() {
				t.Errorf("GreatorExp.Calc() = true, want false")
			}
		})
	}
}

func TestEqualExp_Calc(t *testing.T) {
	tests := []struct {
		name string
		lh   *Value
		rh   *Value
		want bool
	}{
		{
			name: "5 == 5",
			lh:   Number(5),
			rh:   Number(5),
			want: true,
		},
		{
			name: "5 == 10",
			lh:   Number(5),
			rh:   Number(10),
			want: false,
		},
		{
			name: "abc == abc",
			lh:   Text("abc"),
			rh:   Text("abc"),
			want: true,
		},
		{
			name: "abc == bcd",
			lh:   Text("abc"),
			rh:   Text("bcd"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewEqualExp(tt.lh, tt.rh)
			result := exp.Calc()
			if tt.want && !result.True() {
				t.Errorf("EqualExp.Calc() = false, want true")
			}
			if !tt.want && !result.False() {
				t.Errorf("EqualExp.Calc() = true, want false")
			}
		})
	}
}

func TestLessEqExp_Calc(t *testing.T) {
	tests := []struct {
		name string
		lh   *Value
		rh   *Value
		want bool
	}{
		{
			name: "5 <= 10",
			lh:   Number(5),
			rh:   Number(10),
			want: true,
		},
		{
			name: "5 <= 5",
			lh:   Number(5),
			rh:   Number(5),
			want: true,
		},
		{
			name: "10 <= 5",
			lh:   Number(10),
			rh:   Number(5),
			want: false,
		},
		{
			name: "abc <= abc",
			lh:   Text("abc"),
			rh:   Text("abc"),
			want: true,
		},
		{
			name: "abc <= bcd",
			lh:   Text("abc"),
			rh:   Text("bcd"),
			want: true,
		},
		{
			name: "bcd <= abc",
			lh:   Text("bcd"),
			rh:   Text("abc"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewLessEqExp(tt.lh, tt.rh)
			result := exp.Calc()
			if tt.want && !result.True() {
				t.Errorf("LessEqExp.Calc() = false, want true")
			}
			if !tt.want && !result.False() {
				t.Errorf("LessEqExp.Calc() = true, want false")
			}
		})
	}
}

func TestGreaterEqExp_Calc(t *testing.T) {
	tests := []struct {
		name string
		lh   *Value
		rh   *Value
		want bool
	}{
		{
			name: "10 >= 5",
			lh:   Number(10),
			rh:   Number(5),
			want: true,
		},
		{
			name: "5 >= 5",
			lh:   Number(5),
			rh:   Number(5),
			want: true,
		},
		{
			name: "5 >= 10",
			lh:   Number(5),
			rh:   Number(10),
			want: false,
		},
		{
			name: "bcd >= abc",
			lh:   Text("bcd"),
			rh:   Text("abc"),
			want: true,
		},
		{
			name: "abc >= abc",
			lh:   Text("abc"),
			rh:   Text("abc"),
			want: true,
		},
		{
			name: "abc >= bcd",
			lh:   Text("abc"),
			rh:   Text("bcd"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewGreaterEqExp(tt.lh, tt.rh)
			result := exp.Calc()
			if tt.want && !result.True() {
				t.Errorf("GreaterEqExp.Calc() = false, want true")
			}
			if !tt.want && !result.False() {
				t.Errorf("GreaterEqExp.Calc() = true, want false")
			}
		})
	}
}

func TestNotEqExp_Calc(t *testing.T) {
	tests := []struct {
		name string
		lh   *Value
		rh   *Value
		want bool
	}{
		{
			name: "5 != 10",
			lh:   Number(5),
			rh:   Number(10),
			want: true,
		},
		{
			name: "5 != 5",
			lh:   Number(5),
			rh:   Number(5),
			want: false,
		},
		{
			name: "abc != bcd",
			lh:   Text("abc"),
			rh:   Text("bcd"),
			want: true,
		},
		{
			name: "abc != abc",
			lh:   Text("abc"),
			rh:   Text("abc"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exp := NewNotEqExp(tt.lh, tt.rh)
			result := exp.Calc()
			if tt.want && !result.True() {
				t.Errorf("NotEqExp.Calc() = false, want true")
			}
			if !tt.want && !result.False() {
				t.Errorf("NotEqExp.Calc() = true, want false")
			}
		})
	}
}
