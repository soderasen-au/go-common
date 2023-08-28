package fx

import "github.com/soderasen-au/go-common/util"

type Expression interface {
	Calc() *Value
}

type BinOpExp struct {
	Op Operator
	Lh Expression
	Rh Expression
}

type BinOpExpCreator func(lh *Value, rh *Value) Expression

var (
	binOpExpCreators = map[Operator]BinOpExpCreator{}
)

func RegisterNewBinOpExpCreator(op Operator, creator BinOpExpCreator) {
	binOpExpCreators[op] = creator
}

func NewBinOpExp(op Operator, lh, rh *Value) (Expression, *util.Result) {
	c, ok := binOpExpCreators[op]
	if !ok {
		return nil, util.MsgError("LookupCreator", "No creator for OP: "+string(op))
	}
	return c(lh, rh), nil
}
