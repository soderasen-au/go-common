package fx

//LessOp      Operator = "<"
//GreaterOp   Operator = ">"
//EqualOp     Operator = "=="
//LessEqOp    Operator = "<="
//GreaterEqOp Operator = ">="
//NotEqOp     Operator = "!="

func init() {
	RegisterNewBinOpExpCreator(LessOp, NewLessExp)
	RegisterNewBinOpExpCreator(GreaterOp, NewGreatorExp)
	RegisterNewBinOpExpCreator(EqualOp, NewEqualExp)
	RegisterNewBinOpExpCreator(LessEqOp, NewLessEqExp)
	RegisterNewBinOpExpCreator(GreaterEqOp, NewGreaterEqExp)
	RegisterNewBinOpExpCreator(NotEqOp, NewNotEqExp)
}

type LessExp struct {
	BinOpExp
}

func (e LessExp) Calc() *Value {
	return e.Lh.Calc().Less(e.Rh.Calc())
}

func NewLessExp(lh, rh *Value) Expression {
	return &LessExp{
		BinOpExp{
			Op: LessOp,
			Lh: lh,
			Rh: rh,
		},
	}
}

type GreatorExp struct {
	BinOpExp
}

func (e GreatorExp) Calc() *Value {
	return e.Rh.Calc().Less(e.Lh.Calc())
}

func NewGreatorExp(lh, rh *Value) Expression {
	return &GreatorExp{
		BinOpExp{
			Op: GreaterOp,
			Lh: lh,
			Rh: rh,
		},
	}
}

type EqualExp struct {
	BinOpExp
}

func (e EqualExp) Calc() *Value {
	return e.Lh.Calc().Equal(e.Rh.Calc())
}

func NewEqualExp(lh, rh *Value) Expression {
	return &EqualExp{
		BinOpExp{
			Op: EqualOp,
			Lh: lh,
			Rh: rh,
		},
	}
}

type LessEqExp struct {
	BinOpExp
}

func (e LessEqExp) Calc() *Value {
	lv := e.Lh.Calc()
	rv := e.Rh.Calc()

	eq := lv.Equal(rv)
	lt := lv.Less(rv)
	return lt.Or(eq)
}

func NewLessEqExp(lh, rh *Value) Expression {
	return &LessEqExp{
		BinOpExp{
			Op: LessEqOp,
			Lh: lh,
			Rh: rh,
		},
	}
}

type GreaterEqExp struct {
	BinOpExp
}

func (e GreaterEqExp) Calc() *Value {
	lv := e.Lh.Calc()
	rv := e.Rh.Calc()

	eq := lv.Equal(rv)
	gt := rv.Less(lv)
	return gt.Or(eq)
}

func NewGreaterEqExp(lh, rh *Value) Expression {
	return &GreaterEqExp{
		BinOpExp{
			Op: GreaterEqOp,
			Lh: lh,
			Rh: rh,
		},
	}
}

type NotEqExp struct {
	BinOpExp
}

func (e NotEqExp) Calc() *Value {
	lv := e.Lh.Calc()
	rv := e.Rh.Calc()

	eq := lv.Equal(rv)
	return eq.Not()
}

func NewNotEqExp(lh, rh *Value) Expression {
	return &NotEqExp{
		BinOpExp{
			Op: NotEqOp,
			Lh: lh,
			Rh: rh,
		},
	}
}
