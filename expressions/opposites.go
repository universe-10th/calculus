package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"fmt"
)

type NegatedExpr struct {
	arg Expression
}


func (negated NegatedExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := negated.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ops.Neg(result), nil
	}
}


func (negated NegatedExpr) Derivative(wrt Variable) (Expression, error) {
	if result, err := negated.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Negated(result), nil
	}
}


func (negated NegatedExpr) CollectVariables(variables Variables) {
	negated.arg.CollectVariables(variables)
}


func (negated NegatedExpr) String() string {
	switch v := negated.arg.(type) {
	case AddExpr:
		return fmt.Sprintf("-(%s)", v)
	default:
		return fmt.Sprintf("-%s", v)
	}
}


func Negated(arg Expression) NegatedExpr {
	return NegatedExpr{arg}
}