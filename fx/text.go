package fx

import (
	"github.com/soderasen-au/go-common/util"
	"strings"
)

//IncludeOp Operator = "~="

func init() {
	RegisterNewBinOpExpCreator(IncludeOp, NewIncludeExp)
}

type IncludeExp struct {
	BinOpExp
}

func (e IncludeExp) Calc() *Value {
	lv := e.Lh.Calc()
	rv := e.Rh.Calc()
	if lv.HasError() {
		return Error(lv.Error.With("LH Error"))
	}
	if rv.HasError() {
		return Error(rv.Error.With("RH Error"))
	}
	if lv.IsNumeric {
		return Error(util.MsgError("IncludeExpOperands", "LV is numeric"))
	}
	if rv.IsNumeric {
		return Error(util.MsgError("IncludeExpOperands", "RV is numeric"))
	}
	if lv.Text == nil {
		return Error(util.MsgError("IncludeExpOperands", "LV Text is nil"))
	}
	if rv.Text == nil {
		return Error(util.MsgError("IncludeExpOperands", "RV Text is nil"))
	}
	return Bool(strings.Contains(*lv.Text, *rv.Text))
}

func NewIncludeExp(lh, rh *Value) Expression {
	return &IncludeExp{
		BinOpExp{
			Op: LessOp,
			Lh: lh,
			Rh: rh,
		},
	}
}
