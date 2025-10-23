package fx

import (
	"testing"

	"github.com/soderasen-au/go-common/util"
)

func TestValue_Less(t *testing.T) {
	type fields struct {
		Text      *string
		IsNumeric bool
		Number    *float64
		Error     *util.Result
	}
	type args struct {
		rh Value
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  *Value
	}{
		{
			name: "less",
			fields: fields{
				Text:      util.Ptr("abc"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      util.Ptr("bcd"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			}},
			want:  true,
			want1: True(),
		},
		{
			name: "great",
			fields: fields{
				Text:      util.Ptr("bcd"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      util.Ptr("bcd"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			}},
			want:  false,
			want1: False(),
		},
		{
			name: "equal",
			fields: fields{
				Text:      util.Ptr("abc"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      util.Ptr("abc"),
				IsNumeric: false,
				Number:    nil,
				Error:     nil,
			}},
			want:  false,
			want1: False(),
		},
		{
			name: "lessN",
			fields: fields{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.0),
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.1),
				Error:     nil,
			}},
			want:  true,
			want1: True(),
		},
		{
			name: "greatN",
			fields: fields{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.1),
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.0),
				Error:     nil,
			}},
			want:  false,
			want1: False(),
		},
		{
			name: "equalN",
			fields: fields{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.0),
				Error:     nil,
			},
			args: args{rh: Value{
				Text:      nil,
				IsNumeric: true,
				Number:    util.Ptr(1.0),
				Error:     nil,
			}},
			want:  false,
			want1: False(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lh := &Value{
				Text:      tt.fields.Text,
				IsNumeric: tt.fields.IsNumeric,
				Number:    tt.fields.Number,
				Error:     tt.fields.Error,
			}
			got := lh.Less(&tt.args.rh)
			if got.Equal(tt.want1).Not().True() {
				t.Errorf("Less() got = %v, want %v", util.JsonStr(got), util.JsonStr(tt.want1))
			}
		})
	}
}

func TestValue_Factories(t *testing.T) {
	t.Run("False", func(t *testing.T) {
		v := False()
		if !v.False() {
			t.Error("False() should create a false value")
		}
		if v.True() {
			t.Error("False() should not be true")
		}
	})

	t.Run("True", func(t *testing.T) {
		v := True()
		if !v.True() {
			t.Error("True() should create a true value")
		}
		if v.False() {
			t.Error("True() should not be false")
		}
	})

	t.Run("Nil", func(t *testing.T) {
		v := Nil()
		if !v.IsNil() {
			t.Error("Nil() should create a nil value")
		}
	})

	t.Run("Number", func(t *testing.T) {
		v := Number(42.5)
		if !v.IsNumeric {
			t.Error("Number() should create a numeric value")
		}
		if v.Number == nil || *v.Number != 42.5 {
			t.Errorf("Number() = %v, want 42.5", v.Number)
		}
	})

	t.Run("Text", func(t *testing.T) {
		v := Text("hello")
		if v.IsNumeric {
			t.Error("Text() should not be numeric")
		}
		if v.Text == nil || *v.Text != "hello" {
			t.Errorf("Text() = %v, want hello", v.Text)
		}
	})

	t.Run("Error", func(t *testing.T) {
		err := util.MsgError("test", "error message")
		v := Error(err)
		if !v.HasError() {
			t.Error("Error() should have error")
		}
		if v.Error == nil {
			t.Error("Error() should set Error field")
		}
	})

	t.Run("Bool true", func(t *testing.T) {
		v := Bool(true)
		if !v.True() {
			t.Error("Bool(true) should be true")
		}
	})

	t.Run("Bool false", func(t *testing.T) {
		v := Bool(false)
		if !v.False() {
			t.Error("Bool(false) should be false")
		}
	})

	t.Run("Dual numeric", func(t *testing.T) {
		v := Dual("123.45")
		if !v.IsNumeric {
			t.Error("Dual(\"123.45\") should be numeric")
		}
		if v.Number == nil || *v.Number != 123.45 {
			t.Errorf("Dual(\"123.45\") = %v, want 123.45", v.Number)
		}
	})

	t.Run("Dual text", func(t *testing.T) {
		v := Dual("hello")
		if v.IsNumeric {
			t.Error("Dual(\"hello\") should not be numeric")
		}
		if v.Text == nil || *v.Text != "hello" {
			t.Errorf("Dual(\"hello\") = %v, want hello", v.Text)
		}
	})
}

func TestValue_Setters(t *testing.T) {
	t.Run("SetNum", func(t *testing.T) {
		v := &Value{}
		v.SetNum(99.9)
		if !v.IsNumeric {
			t.Error("SetNum should set IsNumeric to true")
		}
		if v.Number == nil || *v.Number != 99.9 {
			t.Errorf("SetNum(99.9) = %v, want 99.9", v.Number)
		}
		if v.Text != nil {
			t.Error("SetNum should clear Text")
		}
		if v.Error != nil {
			t.Error("SetNum should clear Error")
		}
	})

	t.Run("SetText", func(t *testing.T) {
		v := &Value{}
		v.SetText("test")
		if v.IsNumeric {
			t.Error("SetText should set IsNumeric to false")
		}
		if v.Text == nil || *v.Text != "test" {
			t.Errorf("SetText(\"test\") = %v, want test", v.Text)
		}
		if v.Number != nil {
			t.Error("SetText should clear Number")
		}
		if v.Error != nil {
			t.Error("SetText should clear Error")
		}
	})

	t.Run("SetError", func(t *testing.T) {
		v := &Value{}
		v.SetNum(123)
		err := util.MsgError("test", "error")
		v.SetError(err)
		if v.Text != nil {
			t.Error("SetError should clear Text")
		}
		if v.Number != nil {
			t.Error("SetError should clear Number")
		}
		if v.IsNumeric {
			t.Error("SetError should set IsNumeric to false")
		}
		if v.Error == nil {
			t.Error("SetError should set Error")
		}
	})

	t.Run("SetValue with number", func(t *testing.T) {
		v := &Value{}
		v.SetValue("456.78")
		if !v.IsNumeric {
			t.Error("SetValue(\"456.78\") should be numeric")
		}
		if v.Number == nil || *v.Number != 456.78 {
			t.Errorf("SetValue(\"456.78\") = %v, want 456.78", v.Number)
		}
	})

	t.Run("SetValue with text", func(t *testing.T) {
		v := &Value{}
		v.SetValue("not a number")
		if v.IsNumeric {
			t.Error("SetValue(\"not a number\") should not be numeric")
		}
		if v.Text == nil || *v.Text != "not a number" {
			t.Errorf("SetValue = %v, want 'not a number'", v.Text)
		}
	})

	t.Run("SetTrue", func(t *testing.T) {
		v := &Value{}
		v.SetTrue()
		if !v.True() {
			t.Error("SetTrue() should make value true")
		}
		if v.False() {
			t.Error("SetTrue() should not be false")
		}
	})

	t.Run("SetFalse", func(t *testing.T) {
		v := &Value{}
		v.SetFalse()
		if !v.False() {
			t.Error("SetFalse() should make value false")
		}
		if v.True() {
			t.Error("SetFalse() should not be true")
		}
	})
}

func TestValue_Predicates(t *testing.T) {
	t.Run("IsNil on nil value", func(t *testing.T) {
		v := Nil()
		if !v.IsNil() {
			t.Error("Nil value should return true for IsNil()")
		}
	})

	t.Run("IsNil on non-nil value", func(t *testing.T) {
		v := Number(5)
		if v.IsNil() {
			t.Error("Non-nil value should return false for IsNil()")
		}
	})

	t.Run("HasError with error", func(t *testing.T) {
		v := Error(util.MsgError("test", "error"))
		if !v.HasError() {
			t.Error("Value with error should return true for HasError()")
		}
	})

	t.Run("HasError without error", func(t *testing.T) {
		v := Number(5)
		if v.HasError() {
			t.Error("Value without error should return false for HasError()")
		}
	})

	t.Run("IsBool on true", func(t *testing.T) {
		v := True()
		if !v.IsBool() {
			t.Error("True value should return true for IsBool()")
		}
	})

	t.Run("IsBool on false", func(t *testing.T) {
		v := False()
		if !v.IsBool() {
			t.Error("False value should return true for IsBool()")
		}
	})

	t.Run("IsBool on non-bool", func(t *testing.T) {
		v := Number(5)
		if v.IsBool() {
			t.Error("Non-bool value should return false for IsBool()")
		}
	})
}

func TestValue_String(t *testing.T) {
	tests := []struct {
		name string
		v    *Value
		want string
	}{
		{
			name: "nil value",
			v:    Nil(),
			want: "<nil>",
		},
		{
			name: "numeric value",
			v:    Number(42.5),
			want: "42.500000",
		},
		{
			name: "text value",
			v:    Text("hello"),
			want: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.String()
			if got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValue_Or(t *testing.T) {
	tests := []struct {
		name    string
		lh      *Value
		rh      *Value
		wantVal bool
		wantErr bool
	}{
		{
			name:    "true OR true",
			lh:      True(),
			rh:      True(),
			wantVal: true,
			wantErr: false,
		},
		{
			name:    "true OR false",
			lh:      True(),
			rh:      False(),
			wantVal: true,
			wantErr: false,
		},
		{
			name:    "false OR true",
			lh:      False(),
			rh:      True(),
			wantVal: true,
			wantErr: false,
		},
		{
			name:    "false OR false",
			lh:      False(),
			rh:      False(),
			wantVal: false,
			wantErr: false,
		},
		{
			name:    "non-bool OR bool",
			lh:      Number(5),
			rh:      True(),
			wantVal: false,
			wantErr: true,
		},
		{
			name:    "bool OR non-bool",
			lh:      True(),
			rh:      Number(5),
			wantVal: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lh.Or(tt.rh)
			if (got.HasError()) != tt.wantErr {
				t.Errorf("Or() error = %v, wantErr %v", got.HasError(), tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.wantVal && !got.True() {
					t.Errorf("Or() = false, want true")
				}
				if !tt.wantVal && !got.False() {
					t.Errorf("Or() = true, want false")
				}
			}
		})
	}
}

func TestValue_And(t *testing.T) {
	tests := []struct {
		name    string
		lh      *Value
		rh      *Value
		wantVal bool
		wantErr bool
	}{
		{
			name:    "true AND true",
			lh:      True(),
			rh:      True(),
			wantVal: true,
			wantErr: false,
		},
		{
			name:    "true AND false",
			lh:      True(),
			rh:      False(),
			wantVal: false,
			wantErr: false,
		},
		{
			name:    "false AND true",
			lh:      False(),
			rh:      True(),
			wantVal: false,
			wantErr: false,
		},
		{
			name:    "false AND false",
			lh:      False(),
			rh:      False(),
			wantVal: false,
			wantErr: false,
		},
		{
			name:    "non-bool AND bool",
			lh:      Text("test"),
			rh:      True(),
			wantVal: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lh.And(tt.rh)
			if (got.HasError()) != tt.wantErr {
				t.Errorf("And() error = %v, wantErr %v", got.HasError(), tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.wantVal && !got.True() {
					t.Errorf("And() = false, want true")
				}
				if !tt.wantVal && !got.False() {
					t.Errorf("And() = true, want false")
				}
			}
		})
	}
}

func TestValue_Not(t *testing.T) {
	tests := []struct {
		name    string
		v       *Value
		wantVal bool
		wantErr bool
	}{
		{
			name:    "NOT true",
			v:       True(),
			wantVal: false,
			wantErr: false,
		},
		{
			name:    "NOT false",
			v:       False(),
			wantVal: true,
			wantErr: false,
		},
		{
			name:    "NOT non-bool",
			v:       Number(5),
			wantVal: false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v.Not()
			if (got.HasError()) != tt.wantErr {
				t.Errorf("Not() error = %v, wantErr %v", got.HasError(), tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tt.wantVal && !got.True() {
					t.Errorf("Not() = false, want true")
				}
				if !tt.wantVal && !got.False() {
					t.Errorf("Not() = true, want false")
				}
			}
		})
	}
}

func TestValue_Calc(t *testing.T) {
	v := Number(42)
	result := v.Calc()
	if result != v {
		t.Error("Calc() should return the value itself")
	}
}

func TestValue_Less_ErrorCases(t *testing.T) {
	t.Run("left has error", func(t *testing.T) {
		lh := Error(util.MsgError("test", "lh error"))
		rh := Number(5)
		result := lh.Less(rh)
		if !result.HasError() {
			t.Error("Less() should return error when left has error")
		}
	})

	t.Run("right has error", func(t *testing.T) {
		lh := Number(5)
		rh := Error(util.MsgError("test", "rh error"))
		result := lh.Less(rh)
		if !result.HasError() {
			t.Error("Less() should return error when right has error")
		}
	})

	t.Run("type mismatch", func(t *testing.T) {
		lh := Number(5)
		rh := Text("hello")
		result := lh.Less(rh)
		if !result.HasError() {
			t.Error("Less() should return error for type mismatch")
		}
	})

	t.Run("nil number on left", func(t *testing.T) {
		lh := &Value{IsNumeric: true, Number: nil}
		rh := Number(5)
		result := lh.Less(rh)
		if !result.HasError() {
			t.Error("Less() should return error for nil number")
		}
	})

	t.Run("nil number on right", func(t *testing.T) {
		lh := Number(5)
		rh := &Value{IsNumeric: true, Number: nil}
		result := lh.Less(rh)
		if !result.HasError() {
			t.Error("Less() should return error for nil number")
		}
	})

	t.Run("nil text on left", func(t *testing.T) {
		lh := &Value{IsNumeric: false, Text: nil}
		rh := Text("hello")
		result := lh.Less(rh)
		if !result.HasError() {
			t.Error("Less() should return error for nil text")
		}
	})

	t.Run("nil text on right", func(t *testing.T) {
		lh := Text("hello")
		rh := &Value{IsNumeric: false, Text: nil}
		result := lh.Less(rh)
		if !result.HasError() {
			t.Error("Less() should return error for nil text")
		}
	})
}
