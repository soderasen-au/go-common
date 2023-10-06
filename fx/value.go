package fx

import (
	"fmt"
	"github.com/soderasen-au/go-common/util"
	"strconv"
)

type (
	Value struct {
		Text      *string      `json:"text,omitempty" yaml:"text,omitempty"`
		IsNumeric bool         `json:"is_numeric" yaml:"is_numeric"`
		Number    *float64     `json:"number,omitempty" yaml:"number,omitempty"`
		Error     *util.Result `json:"error,omitempty" yaml:"error,omitempty"`
	}
)

func (v *Value) Calc() *Value {
	return v
}

func (v *Value) IsNil() bool {
	return v.Text == nil && v.Number == nil && v.Error == nil
}

func (v *Value) HasError() bool {
	return v.Error != nil
}

func (v *Value) SetError(err *util.Result) {
	v.Text = nil
	v.IsNumeric = false
	v.Number = nil
	v.Error = err
}

func (v *Value) SetFalse() *Value {
	v.Text = util.Ptr("false")
	v.IsNumeric = true
	v.Number = util.Ptr(0.0)
	v.Error = nil
	return v
}

func (v *Value) SetTrue() *Value {
	v.Text = util.Ptr("true")
	v.IsNumeric = true
	v.Number = util.Ptr(1.0)
	v.Error = nil
	return v
}

func (v *Value) False() bool {
	return v.Error == nil && util.MaybeNil(v.Text) == "false" && v.IsNumeric == true && util.MaybeNil(v.Number) == 0
}

func (v *Value) True() bool {
	return v.Error == nil && util.MaybeNil(v.Text) == "true" && v.IsNumeric == true && util.MaybeNil(v.Number) == 1.0
}

func (v *Value) IsBool() bool {
	return v.True() || v.False()
}

func (v *Value) SetNum(n float64) *Value {
	v.Text = nil
	v.IsNumeric = true
	v.Number = util.Ptr(n)
	v.Error = nil
	return v
}

func (v *Value) SetText(t string) *Value {
	v.Text = util.Ptr(t)
	v.IsNumeric = false
	v.Number = nil
	v.Error = nil
	return v
}

func (v *Value) SetValue(s string) *Value {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		v.SetText(s)
	} else {
		v.SetNum(f)
	}
	return v
}

func (lh *Value) Less(rh *Value) *Value {
	if lh.HasError() {
		return Error(lh.Error.With("LH Error"))
	}
	if rh.HasError() {
		return Error(rh.Error.With("RH Error"))
	}

	if lh.IsNumeric && rh.IsNumeric {
		if lh.Number == nil || rh.Number == nil {
			return Error(util.MsgError("NumberValueLess", "operand Number can't be nil"))
		}
		return Bool(*lh.Number < *rh.Number)
	} else if !lh.IsNumeric && !rh.IsNumeric {
		if lh.Text == nil || rh.Text == nil {
			return Error(util.MsgError("TextValueLess", "operand Text can't be nil"))
		}
		return Bool(*lh.Text < *rh.Text)
	}

	return Error(util.MsgError("ValueLess", "operand must be in same type"))
}

func (lh *Value) Equal(rh *Value) *Value {
	return Bool(lh.IsNumeric == rh.IsNumeric &&
		util.MaybeNil(lh.Text) == util.MaybeNil(rh.Text) &&
		util.MaybeNil(lh.Number) == util.MaybeNil(rh.Number))
}

func (lh *Value) Or(rh *Value) *Value {
	if !lh.IsBool() {
		return Error(util.MsgError("Or", "lh is not bool"))
	}
	if !rh.IsBool() {
		return Error(util.MsgError("Or", "rh is not bool"))
	}
	return Bool(lh.True() || rh.True())
}

func (lh *Value) And(rh *Value) *Value {
	if !lh.IsBool() {
		return Error(util.MsgError("Or", "lh is not bool"))
	}
	if !rh.IsBool() {
		return Error(util.MsgError("Or", "rh is not bool"))
	}
	return Bool(lh.True() && rh.True())
}

func (lh *Value) Not() *Value {
	if !lh.IsBool() {
		return Error(util.MsgError("Or", "lh is not bool"))
	}
	return Bool(lh.False())
}

func (v *Value) String() string {
	if v.IsNil() {
		return "<nil>"
	} else if v.IsNumeric {
		return fmt.Sprintf("%f", *v.Number)
	}
	return util.MaybeNil(v.Text)
}

func False() *Value {
	v := &Value{}
	return v.SetFalse()
}
func True() *Value {
	v := &Value{}
	return v.SetTrue()
}

func Nil() *Value {
	return &Value{
		Text:      nil,
		IsNumeric: false,
		Number:    nil,
		Error:     nil,
	}
}

func Number(n float64) *Value {
	v := Value{}
	v.SetNum(n)
	return &v
}

func Text(s string) *Value {
	v := Value{}
	v.SetText(s)
	return &v
}

func Error(err *util.Result) *Value {
	v := Value{}
	v.SetError(err)
	return &v
}

func Bool(b bool) *Value {
	v := Value{}
	if b {
		v.SetTrue()
	} else {
		v.SetFalse()
	}
	return &v
}

func Dual(str string) *Value {
	v := Value{}
	v.SetValue(str)
	return &v
}
