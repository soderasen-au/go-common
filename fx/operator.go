package fx

import "github.com/soderasen-au/go-common/util"

type Operator string

const (
	//Unknown
	UnkownOp Operator = ""

	//Relation Ops
	LessOp      Operator = "<"
	GreaterOp   Operator = ">"
	EqualOp     Operator = "=="
	LessEqOp    Operator = "<="
	GreaterEqOp Operator = ">="
	NotEqOp     Operator = "!="

	//Text Ops
	IncludeOp Operator = "~="
)

func ParseOperator(str string) (Operator, *util.Result) {
	switch str {
	case string(LessOp):
		return LessOp, nil
	case string(GreaterOp):
		return GreaterOp, nil
	case string(EqualOp):
		return EqualOp, nil
	case string(LessEqOp):
		return LessEqOp, nil
	case string(GreaterEqOp):
		return GreaterEqOp, nil
	case string(NotEqOp):
		return NotEqOp, nil
	case string(IncludeOp):
		return IncludeOp, nil
	}
	return UnkownOp, util.MsgError(str, "Unknown operator")
}
