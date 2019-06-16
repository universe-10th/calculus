package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"fmt"
	"github.com/universe-10th-calculus/errors"
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
		return Negated(result).Simplify()
	}
}


func (negated NegatedExpr) CollectVariables(variables Variables) {
	negated.arg.CollectVariables(variables)
}


func (negated NegatedExpr) IsConstant(wrt Variable) bool {
	return negated.arg.IsConstant(wrt)
}


func (negated NegatedExpr) Simplify() (Expression, error) {
	if simplified, err := negated.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		return Constant{ops.Neg(num.number)}, nil
	} else {
		return Negated(simplified), nil
	}
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
			return Negated(Mul(derivative, Pow(inverse.arg, Num(-2)))).Simplify()
		}
	}
}


func (inverse InverseExpr) CollectVariables(variables Variables) {
	inverse.arg.CollectVariables(variables)
}


func (inverse InverseExpr) IsConstant(wrt Variable) bool {
	return inverse.arg.IsConstant(wrt)
}


func (inverse InverseExpr) wrappedInverse(value sets.Number) (result sets.Number, err error) {
	defer func(){
		if r := recover(); r != nil {
			result = nil
			err = errors.DivisionByZero
		}
	}()
	result = ops.Inv(value)
	return
}


func (inverse InverseExpr) Simplify() (Expression, error) {
	if simplified, err := inverse.arg.Simplify(); err != nil {
		return nil, err
	} else {
		if num, ok := simplified.(Constant); ok {
			if result, err := inverse.wrappedInverse(num.number); err != nil {
				return nil, err
			} else {
				return Constant{result}, nil
			}
		} else {
			return Inverse(simplified), nil
		}
	}
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