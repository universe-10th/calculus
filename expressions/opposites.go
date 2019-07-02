package expressions

import (
	"github.com/universe-10th-calculus/sets"
	"github.com/universe-10th-calculus/ops"
	"fmt"
	"github.com/universe-10th-calculus/errors"
)


// NegatedExpr is a negated expression, like -X or -(X + Y).
type NegatedExpr struct {
	arg Expression
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (negated NegatedExpr) Curry(args Arguments) (Expression, error) {
	if curried, err := negated.arg.Curry(args); err != nil {
		return nil, err
	} else {
		return Negated(curried).Simplify()
	}
}


// Evaluate computes the inner expression's evaluated value and negates it.
func (negated NegatedExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := negated.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ops.Neg(result), nil
	}
}


// Derivative just computes the derivative of the inner expression and changes its sign.
func (negated NegatedExpr) Derivative(wrt Variable) (Expression, error) {
	if result, err := negated.arg.Derivative(wrt); err != nil {
		return nil, err
	} else {
		return Negated(result).Simplify()
	}
}


// CollectVariables digs into the inner expression.
func (negated NegatedExpr) CollectVariables(variables Variables) {
	negated.arg.CollectVariables(variables)
}


// IsConstant tells whether the inner expression is constant with respect to the given value or not.
func (negated NegatedExpr) IsConstant(wrt Variable) bool {
	return negated.arg.IsConstant(wrt)
}


// Simplify evaluates the simplified value of the inner expression, and optimizes it.
// This means: constant simplified expressions are negated, and negated simplified expressions are unwrapped.
func (negated NegatedExpr) Simplify() (Expression, error) {
	if simplified, err := negated.arg.Simplify(); err != nil {
		return nil, err
	} else if num, ok := simplified.(Constant); ok {
		return Constant{ops.Neg(num.number)}, nil
	} else {
		return Negated(simplified), nil
	}
}


// String represents the negation appropriately, as -X or -(X).
func (negated NegatedExpr) String() string {
	switch v := negated.arg.(type) {
	case AddExpr:
		return fmt.Sprintf("-(%s)", v)
	default:
		return fmt.Sprintf("-%s", v)
	}
}


// Negated constructs a negated expression.
func Negated(arg Expression) Expression {
	if neg, ok := arg.(NegatedExpr); ok {
		return neg.arg
	} else if con, ok := arg.(Constant); ok && ops.IsNegative(con.number) {
		return Constant{ops.Neg(con.number)}
	} else {
		return NegatedExpr{arg}
	}
}


// InverseExpr is an inverted expression, like 1/X or 1/(1 + X).
type InverseExpr struct {
	arg Expression
}


// Curry tries currying the underlying expression first, and then attempts simplifying.
func (inverse InverseExpr) Curry(args Arguments) (Expression, error) {
	if curried, err := inverse.arg.Curry(args); err != nil {
		return nil, err
	} else {
		return Inverse(curried).Simplify()
	}
}


// Evaluate computes the value of the inner expression, and then divides 1 by it.
// It returns an error if a division-by-zero occurs.
func (inverse InverseExpr) Evaluate(args Arguments) (sets.Number, error) {
	if result, err := inverse.arg.Evaluate(args); err != nil {
		return nil, err
	} else {
		return ops.Inv(result), nil
	}
}


// Derivative computes the 1/X derivative rule, also applying chain rule appropriately.
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


// CollectVariables digs into the inner term.
func (inverse InverseExpr) CollectVariables(variables Variables) {
	inverse.arg.CollectVariables(variables)
}


// IsConstant returns whether the inner (to-invert) term is constant with respect to the given variable.
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


// Simplify tries simplifying the inner expression and, if constant, computing the inverse and returning it as constant.
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


// String represents the inverse of a value as X^-1 or (X)^-1.
func (inverse InverseExpr) String() string {
	if _, ok := inverse.arg.(SelfContained); !ok {
		return fmt.Sprintf("(%s)^-1", inverse.arg)
	} else {
		return fmt.Sprintf("%s^-1", inverse.arg)
	}
}


// Inverse constructs an inverted expression node.
func Inverse(arg Expression) Expression {
	if inv, ok := arg.(InverseExpr); ok {
		return inv.arg
	} else if neg, ok := arg.(NegatedExpr); ok {
		return NegatedExpr{InverseExpr{neg.arg}}
	} else {
		return InverseExpr{arg}
	}
}