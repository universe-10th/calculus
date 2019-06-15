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


func (negated NegatedExpr) IsConstant() bool {
	return negated.arg.IsConstant()
}


func (negated NegatedExpr) String() string {
	switch v := negated.arg.(type) {
	case AddExpr:
		return fmt.Sprintf("-(%s)", v)
	default:
		return fmt.Sprintf("-%s", v)
	}
}


func Negated(arg Expression) Expression {
	if neg, ok := arg.(NegatedExpr); ok {
		return neg.arg
	} else {
		return NegatedExpr{arg}
	}
}


type InverseExpr struct {
	arg Expression
}


func (inverse InverseExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := inverse.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ops.Inv(result), nil
	}
}


func (inverse InverseExpr) Derivative(wrt Variable) (Expression, error) {
	if derivative, err := inverse.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		variables := Variables{}
		derivative.CollectVariables(variables)
		if _, ok := variables[wrt]; !ok {
			return Num(0), nil
		} else {
			return Negated(Mul(derivative, Pow(inverse.arg, Num(-2)))), nil
		}
	}
	return nil, nil
}


func (inverse InverseExpr) CollectVariables(variables Variables) {
	inverse.arg.CollectVariables(variables)
}


func (inverse InverseExpr) IsConstant() bool {
	return inverse.arg.IsConstant()
}


func (inverse InverseExpr) String() string {
	// TODO
	return ""
}


func Inverse(arg Expression) Expression {
	if inv, ok := arg.(InverseExpr); ok {
		return inv.arg
	} else {
		return InverseExpr{arg}
	}
}